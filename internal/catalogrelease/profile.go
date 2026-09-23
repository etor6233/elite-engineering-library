package catalogrelease

// AUTHORED bounded profile loading; activation never accepts an unpinned file.
import (
	"bytes"
	"crypto/sha256"
	"elite.local/enterprise/internal/approval"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

func LoadProfile(path, expected string) (Profile, error) {
	var p Profile
	if !filepath.IsAbs(path) || !ValidSHA(expected) {
		return p, ErrInvalid
	}
	info, e := os.Lstat(path)
	if e != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() < 2 || info.Size() > 16384 {
		return p, ErrInvalid
	}
	f, e := os.Open(path)
	if e != nil {
		return p, ErrInvalid
	}
	defer f.Close()
	actual, e := f.Stat()
	if e != nil || !os.SameFile(info, actual) {
		return p, ErrInvalid
	}
	raw, e := io.ReadAll(io.LimitReader(f, 16385))
	if e != nil || len(raw) > 16384 {
		return p, ErrInvalid
	}
	digest := sha256.Sum256(raw)
	if hex.EncodeToString(digest[:]) != expected {
		return p, ErrInvalid
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return p, ErrInvalid
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&p) != nil || dec.Decode(new(any)) != io.EOF || !p.Valid() {
		return Profile{}, ErrInvalid
	}
	return p, nil
}

func LoadFeedProfile(path, expected string) (FeedProfile, error) {
	var p FeedProfile
	if !filepath.IsAbs(path) || !ValidSHA(expected) {
		return p, ErrInvalid
	}
	info, e := os.Lstat(path)
	if e != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() < 2 || info.Size() > 16384 {
		return p, ErrInvalid
	}
	f, e := os.Open(path)
	if e != nil {
		return p, ErrInvalid
	}
	defer f.Close()
	actual, e := f.Stat()
	if e != nil || !os.SameFile(info, actual) {
		return p, ErrInvalid
	}
	raw, e := io.ReadAll(io.LimitReader(f, 16385))
	if e != nil || len(raw) > 16384 {
		return p, ErrInvalid
	}
	digest := sha256.Sum256(raw)
	if hex.EncodeToString(digest[:]) != expected {
		return p, ErrInvalid
	}
	if _, _, e = approval.CanonicalPayload(raw); e != nil {
		return p, ErrInvalid
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	if dec.Decode(&p) != nil || dec.Decode(new(any)) != io.EOF || !p.Valid() {
		return FeedProfile{}, ErrInvalid
	}
	return p, nil
}
