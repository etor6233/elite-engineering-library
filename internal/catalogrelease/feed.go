package catalogrelease

// AUTHORED reference feed protocol over pinned Go net/http. Provider-specific
// marketplace SDKs remain separate. The receiver must preserve per-key receipts
// and reject an older generation after a newer one; it is not a blind webhook.
import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/tls"
	"elite.local/enterprise/internal/approval"
	"elite.local/enterprise/internal/outbounddelivery"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type FeedProfile struct {
	Schema     string `json:"schema"`
	ReceiverID string `json:"receiver_id"`
	BaseURL    string `json:"base_url"`
	Mode       string `json:"mode"`
}

func (p FeedProfile) Valid() bool {
	if p.Schema != "elite-catalog-feed/v1" || !ValidID(p.ReceiverID) || len(p.BaseURL) > 1024 {
		return false
	}
	u, e := url.Parse(p.BaseURL)
	if e != nil || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" || u.Path != "" || u.Opaque != "" {
		return false
	}
	if p.Mode == "fixture" {
		return u.Scheme == "http" && (u.Hostname() == "127.0.0.1" || u.Hostname() == "::1") && u.Port() != ""
	}
	return p.Mode == "https" && u.Scheme == "https"
}

type FeedAck struct {
	Status        string    `json:"status"`
	Generation    int64     `json:"generation,string"`
	PayloadSHA256 string    `json:"payload_sha256"`
	RemoteID      string    `json:"remote_id"`
	AcceptedAt    time.Time `json:"accepted_at"`
}
type FeedHTTP struct {
	profile    FeedProfile
	profileSHA string
	token      string
	client     *http.Client
	transport  *http.Transport
}

func NewFeedHTTP(p FeedProfile, token string) (*FeedHTTP, error) {
	if !p.Valid() || !ValidText(token, 4096) {
		return nil, ErrInvalid
	}
	_, sha, e := Canonical(p)
	if e != nil {
		return nil, e
	}
	transport := &http.Transport{TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12}, ResponseHeaderTimeout: 5 * time.Second, MaxResponseHeaderBytes: 16384, DisableKeepAlives: true}
	client := &http.Client{Transport: transport, Timeout: 8 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	return &FeedHTTP{p, sha, token, client, transport}, nil
}
func (c *FeedHTTP) Profile() FeedProfile  { return c.profile }
func (c *FeedHTTP) ProfileSHA256() string { return c.profileSHA }
func (c *FeedHTTP) Close()                { c.transport.CloseIdleConnections() }
func (c *FeedHTTP) exchange(ctx context.Context, method, key string, generation int64, payload []byte, payloadSHA string) (outbounddelivery.Receipt, error) {
	if !ValidText(key, 128) || generation < 1 || !ValidSHA(payloadSHA) {
		return outbounddelivery.Receipt{}, ErrInvalid
	}
	req, e := http.NewRequestWithContext(ctx, method, c.profile.BaseURL+"/publications/"+url.PathEscape(key), bytes.NewReader(payload))
	if e != nil {
		return outbounddelivery.Receipt{}, ErrInvalid
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", key)
	req.Header.Set("X-Catalog-Generation", fmt.Sprint(generation))
	req.Header.Set("X-Content-SHA256", payloadSHA)
	response, e := c.client.Do(req)
	if e != nil {
		return outbounddelivery.Receipt{}, errors.New("feed response unconfirmed")
	}
	defer response.Body.Close()
	raw, e := io.ReadAll(io.LimitReader(response.Body, 16385))
	if e != nil || len(raw) > 16384 {
		return outbounddelivery.Receipt{}, errors.New("feed response exceeds boundary")
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return outbounddelivery.Receipt{}, errors.New("feed receipt malformed")
	}
	var ack FeedAck
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&ack) != nil || dec.Decode(new(any)) != io.EOF || ack.Generation != generation || ack.PayloadSHA256 != payloadSHA {
		return outbounddelivery.Receipt{}, errors.New("feed receipt binding mismatch")
	}
	digest := sha256.Sum256(raw)
	evidence := hex.EncodeToString(digest[:])
	if ack.Status == "rejected" && (response.StatusCode == 422 || method == "GET" && response.StatusCode == 200) {
		return outbounddelivery.Receipt{}, outbounddelivery.NewTerminalFailure("FEED_REJECTED", evidence, nil)
	}
	if ack.Status != "accepted" || (response.StatusCode != 200 && response.StatusCode != 201) || strings.TrimSpace(ack.RemoteID) == "" {
		return outbounddelivery.Receipt{}, errors.New("feed receipt not accepted")
	}
	receipt := outbounddelivery.Receipt{ProviderMessageID: ack.RemoteID, EvidenceSHA256: evidence, AcceptedAt: ack.AcceptedAt}
	if e = receipt.Validate(); e != nil {
		return outbounddelivery.Receipt{}, e
	}
	return receipt, nil
}
func (c *FeedHTTP) Send(ctx context.Context, key string, generation int64, payload []byte, payloadSHA string) (outbounddelivery.Receipt, error) {
	sum := sha256.Sum256(payload)
	if hex.EncodeToString(sum[:]) != payloadSHA || len(payload) > 524288 {
		return outbounddelivery.Receipt{}, ErrInvalid
	}
	return c.exchange(ctx, "POST", key, generation, payload, payloadSHA)
}
func (c *FeedHTTP) Reconcile(ctx context.Context, key string, generation int64, payloadSHA string) (outbounddelivery.Receipt, error) {
	return c.exchange(ctx, "GET", key, generation, nil, payloadSHA)
}

type FeedStatus struct {
	DeliveryKey   string `json:"delivery_key"`
	Generation    int64  `json:"generation,string"`
	State         string `json:"state"`
	PayloadSHA256 string `json:"payload_sha256"`
	SourceSHA256  string `json:"source_sha256"`
}
