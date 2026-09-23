package catalogrelease

// DEPENDENCY_PIN: existing Go1.26.8 standard library image/png performs PNG
// decoding/reencoding. This AUTHORED wrapper bounds resources and records hashes;
// it does not claim malware scanning or native-image-library repair.
import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image/png"
)

func NormalizePNG(raw []byte) (Media, []byte, error) {
	if len(raw) == 0 || len(raw) > 1<<20 {
		return Media{}, nil, ErrInvalid
	}
	cfg, err := png.DecodeConfig(bytes.NewReader(raw))
	if err != nil || cfg.Width < 1 || cfg.Height < 1 || cfg.Width > 2048 || cfg.Height > 2048 || cfg.Width*cfg.Height > 1048576 {
		return Media{}, nil, ErrInvalid
	}
	im, err := png.Decode(bytes.NewReader(raw))
	if err != nil {
		return Media{}, nil, ErrInvalid
	}
	var out bytes.Buffer
	if err = png.Encode(&out, im); err != nil || out.Len() > 4<<20 {
		return Media{}, nil, ErrInvalid
	}
	original := sha256.Sum256(raw)
	clean := sha256.Sum256(out.Bytes())
	return Media{OriginalSHA256: hex.EncodeToString(original[:]), SHA256: hex.EncodeToString(clean[:]), Width: cfg.Width, Height: cfg.Height}, out.Bytes(), nil
}
