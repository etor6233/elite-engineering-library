package catalogrelease

import (
	"bytes"
	"crypto/sha256"
	"elite.local/enterprise/internal/approval"
	"encoding/hex"
	"image"
	"image/png"
	"testing"
)

func FuzzCatalogPublicationBoundaries(f *testing.F) {
	var seed bytes.Buffer
	if e := png.Encode(&seed, image.NewNRGBA(image.Rect(0, 0, 2, 2))); e != nil {
		f.Fatal(e)
	}
	for _, raw := range [][]byte{seed.Bytes(), append(append([]byte(nil), seed.Bytes()...), []byte("uninterpreted trailer")...),
		[]byte(`{"command_id":"publish-one","draft_id":"draft-one","expected_generation":"2","reason":"approved snapshot"}`),
		[]byte(`{"command_id":"one","command_id":"two"}`), []byte(`{"amount_minor_units":"9223372036854775807"}`), []byte{}, []byte("<svg/>")} {
		f.Add(raw)
	}
	f.Fuzz(func(t *testing.T, raw []byte) {
		if len(raw) > 1<<20 {
			return
		}
		if canonical, hash, e := approval.CanonicalPayload(raw); e == nil {
			again, againHash, e := approval.CanonicalPayload(canonical)
			if e != nil || !bytes.Equal(canonical, again) || hash != againHash {
				t.Fatal("canonical command identity unstable")
			}
		}
		media, clean, e := NormalizePNG(raw)
		if e != nil {
			return
		}
		originalHash := sha256.Sum256(raw)
		cleanHash := sha256.Sum256(clean)
		if media.OriginalSHA256 != hex.EncodeToString(originalHash[:]) || media.SHA256 != hex.EncodeToString(cleanHash[:]) ||
			media.Width < 1 || media.Height < 1 || media.Width > 2048 || media.Height > 2048 || media.Width*media.Height > 1048576 || len(clean) > 4<<20 {
			t.Fatal("normalized media escaped declared boundary")
		}
		next, twice, e := NormalizePNG(clean)
		if e != nil || next.SHA256 != media.SHA256 || !bytes.Equal(clean, twice) {
			t.Fatal("normalized media is not stable")
		}
	})
}
