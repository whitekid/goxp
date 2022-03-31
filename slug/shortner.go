package slug

import "math/big"

type Shortner struct {
	slugger *Slug
}

func NewShortner(encoding string) *Shortner {
	return &Shortner{
		slugger: New(encoding),
	}
}

func (sg *Shortner) Encode(v int64) string {
	bi := big.NewInt(v)
	return sg.slugger.Encode(bi.Bytes())
}

func (sg *Shortner) Decode(s string) (int64, error) {
	dec, err := sg.slugger.Decode(s)
	if err != nil {
		return 0, err
	}

	bi := big.NewInt(0)
	bi.SetBytes(dec)

	return bi.Int64(), nil
}
