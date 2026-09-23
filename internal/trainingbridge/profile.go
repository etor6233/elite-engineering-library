// Package trainingbridge composes versioned help, participation evidence and
// the existing human approval owner. It never scores or grants permissions.
package trainingbridge

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"

	"elite.local/enterprise/internal/approval"
)

var ErrContract = errors.New("training: invalid or inactive contract")
var ErrScope = errors.New("training: not authorized or not found")
var ErrConflict = errors.New("training: recorded identity or content differs")
var codeRE = regexp.MustCompile(`^[a-z][a-z0-9-]{0,63}$`)
var hexRE = regexp.MustCompile(`^[0-9a-f]{64}$`)

const Method = "HUMAN_REVIEW_NO_GRANTS_V1"

type Article struct {
	ID         string   `json:"id"`
	Version    string   `json:"version"`
	Title      string   `json:"title"`
	Paragraphs []string `json:"paragraphs"`
}
type ContentBundle struct {
	Schema       string    `json:"schema"`
	SourceSHA256 string    `json:"source_sha256"`
	Articles     []Article `json:"articles"`
}
type Prompt struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}
type Course struct {
	ID      string   `json:"id"`
	Title   string   `json:"title"`
	Role    string   `json:"role"`
	Lessons []string `json:"lessons"`
	Prompts []Prompt `json:"prompts"`
}
type ProfileDocument struct {
	Schema         string   `json:"schema"`
	ID             string   `json:"id"`
	Revision       int      `json:"revision"`
	TenantID       string   `json:"tenant_id"`
	OrganizationID string   `json:"organization_id"`
	Method         string   `json:"method"`
	ContentSHA256  string   `json:"content_sha256"`
	SourceSHA256   string   `json:"source_sha256"`
	Courses        []Course `json:"courses"`
}
type Activation struct {
	Enabled        bool
	ID             string
	Revision       int
	SHA256         string
	TenantID       string
	OrganizationID string
}
type Profile struct {
	document ProfileDocument
	hash     string
	articles map[string]Article
}
type CourseView struct {
	Course          Course    `json:"course"`
	Articles        []Article `json:"articles"`
	ProfileID       string    `json:"profile_id"`
	ProfileRevision int       `json:"profile_revision"`
	ProfileSHA256   string    `json:"profile_sha256"`
	ContentSHA256   string    `json:"content_sha256"`
	Method          string    `json:"method"`
}

func digest(b []byte) string { s := sha256.Sum256(b); return hex.EncodeToString(s[:]) }
func strict(raw []byte, out any) error {
	if len(raw) == 0 || len(raw) > 65536 || !utf8.Valid(raw) {
		return ErrContract
	}
	if _, _, e := approval.CanonicalPayload(raw); e != nil {
		return ErrContract
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(out) != nil {
		return ErrContract
	}
	var trailing any
	if dec.Decode(&trailing) != io.EOF {
		return ErrContract
	}
	return nil
}
func text(s string, max int) bool {
	if len(s) < 1 || len(s) > max || !utf8.ValidString(s) {
		return false
	}
	for _, r := range s {
		if (r < 32 && r != '\n' && r != '\t') || r == 127 {
			return false
		}
	}
	return true
}
func readBounded(path string) ([]byte, error) {
	if !filepath.IsAbs(path) {
		return nil, ErrContract
	}
	before, e := os.Lstat(path)
	if e != nil || !before.Mode().IsRegular() || before.Size() > 65536 {
		return nil, ErrContract
	}
	f, e := os.Open(path)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	st, e := f.Stat()
	if e != nil || !st.Mode().IsRegular() || st.Size() > 65536 || !os.SameFile(before, st) {
		return nil, ErrContract
	}
	b, e := io.ReadAll(io.LimitReader(f, 65537))
	if len(b) > 65536 {
		return nil, ErrContract
	}
	return b, e
}
func LoadProfile(path, contentPath string, a Activation) (*Profile, error) {
	raw, e := readBounded(path)
	if e != nil {
		return nil, e
	}
	content, e := readBounded(contentPath)
	if e != nil {
		return nil, e
	}
	return ParseProfile(raw, content, a)
}
func ParseProfile(raw, content []byte, a Activation) (*Profile, error) {
	var d ProfileDocument
	var b ContentBundle
	if !a.Enabled || !hexRE.MatchString(a.SHA256) || digest(raw) != a.SHA256 || strict(raw, &d) != nil || strict(content, &b) != nil {
		return nil, ErrContract
	}
	if d.Schema != "elite-training-profile/v1" || d.Method != Method || !codeRE.MatchString(d.ID) || d.ID != a.ID || d.Revision < 1 || d.Revision != a.Revision || d.TenantID != a.TenantID || d.OrganizationID != a.OrganizationID || !text(d.TenantID, 128) || !text(d.OrganizationID, 128) || !hexRE.MatchString(d.SourceSHA256) || digest(content) != d.ContentSHA256 || b.Schema != "elite-training-content/v1" || b.SourceSHA256 != d.SourceSHA256 || len(b.Articles) < 1 || len(b.Articles) > 64 || len(d.Courses) < 1 || len(d.Courses) > 8 {
		return nil, ErrContract
	}
	p := &Profile{document: d, hash: a.SHA256, articles: map[string]Article{}}
	for _, v := range b.Articles {
		if !codeRE.MatchString(v.ID) || !regexp.MustCompile(`^\d{1,6}\.\d{1,6}\.\d{1,6}$`).MatchString(v.Version) || !text(v.Title, 160) || len(v.Paragraphs) < 1 || len(v.Paragraphs) > 20 {
			return nil, ErrContract
		}
		if _, ok := p.articles[v.ID]; ok {
			return nil, ErrContract
		}
		for _, s := range v.Paragraphs {
			if !text(s, 4000) {
				return nil, ErrContract
			}
		}
		p.articles[v.ID] = v
	}
	seen := map[string]bool{}
	for _, c := range d.Courses {
		if !codeRE.MatchString(c.ID) || seen[c.ID] || !text(c.Title, 160) || len(c.Lessons) < 1 || len(c.Lessons) > 4 || len(c.Prompts) < 1 || len(c.Prompts) > 8 {
			return nil, ErrContract
		}
		seen[c.ID] = true
		switch c.Role {
		case "owner", "admin", "employee", "customer":
		default:
			return nil, ErrContract
		}
		lessons := map[string]bool{}
		for _, id := range c.Lessons {
			if _, ok := p.articles[id]; !ok || lessons[id] {
				return nil, ErrContract
			}
			lessons[id] = true
		}
		prompts := map[string]bool{}
		for _, q := range c.Prompts {
			if !codeRE.MatchString(q.ID) || !text(q.Text, 1000) || prompts[q.ID] {
				return nil, ErrContract
			}
			prompts[q.ID] = true
		}
		view, _ := p.Course(c.ID)
		encoded, _ := json.Marshal(view)
		if len(encoded) > 32768 {
			return nil, ErrContract
		}
	}
	return p, nil
}
func (p *Profile) Scope() (string, string) { return p.document.TenantID, p.document.OrganizationID }
func (p *Profile) Hash() string            { return p.hash }
func (p *Profile) Courses() []CourseView {
	v := []CourseView{}
	for _, c := range p.document.Courses {
		x, _ := p.Course(c.ID)
		v = append(v, x)
	}
	return v
}
func (p *Profile) Course(id string) (CourseView, error) {
	for _, c := range p.document.Courses {
		if c.ID == id {
			c.Lessons = append([]string(nil), c.Lessons...)
			c.Prompts = append([]Prompt(nil), c.Prompts...)
			v := CourseView{Course: c, ProfileID: p.document.ID, ProfileRevision: p.document.Revision, ProfileSHA256: p.hash, ContentSHA256: p.document.ContentSHA256, Method: Method}
			for _, id := range c.Lessons {
				article := p.articles[id]
				article.Paragraphs = append([]string(nil), article.Paragraphs...)
				v.Articles = append(v.Articles, article)
			}
			return v, nil
		}
	}
	return CourseView{}, ErrContract
}
