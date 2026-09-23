package catalogrelease

import (
	"bytes"
	"image"
	"image/png"
	"testing"
)

func TestCatalogMediaStructuralBoundary(t *testing.T) {
	var b bytes.Buffer
	if e := png.Encode(&b, image.NewNRGBA(image.Rect(0, 0, 4, 3))); e != nil {
		t.Fatal(e)
	}
	original := append(append([]byte(nil), b.Bytes()...), []byte("<script>trailer</script>")...)
	m, clean, e := NormalizePNG(original)
	if e != nil || m.Width != 4 || m.Height != 3 || bytes.Contains(clean, []byte("script")) || m.SHA256 == m.OriginalSHA256 {
		t.Fatal("PNG normalization", m, e)
	}
	again, out, e := NormalizePNG(clean)
	if e != nil || !bytes.Equal(out, clean) || again.SHA256 != m.SHA256 {
		t.Fatal("non-deterministic PNG")
	}
	for _, bad := range [][]byte{[]byte("<svg onload='x'/>"), b.Bytes()[:20], make([]byte, (1<<20)+1)} {
		if _, _, e := NormalizePNG(bad); e == nil {
			t.Fatal("invalid input accepted")
		}
	}
	b.Reset()
	if e = png.Encode(&b, image.NewNRGBA(image.Rect(0, 0, 2049, 1))); e != nil {
		t.Fatal(e)
	}
	if _, _, e = NormalizePNG(b.Bytes()); e == nil {
		t.Fatal("dimension budget ignored")
	}
}
