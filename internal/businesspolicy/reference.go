package businesspolicy

// AUTHORED reference-configuration binding. The original values live in JSON,
// not business logic. A changed embedded artifact fails startup validation.
import _ "embed"

//go:embed reference-profile.json
var referenceBytes []byte

const ReferenceSHA256 = "da3feda4d66a19806f20234773720c2f98f6c7d4251640d6fb12cb58b24be5c5"

func Reference() *Profile {
	p, err := Load(referenceBytes, ReferenceSHA256)
	if err != nil {
		panic("invalid embedded business policy profile")
	}
	return p
}

func ReferenceJSON() []byte { return append([]byte(nil), referenceBytes...) }
