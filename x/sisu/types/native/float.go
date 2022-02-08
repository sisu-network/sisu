package native

import (
	"math/big"
)

type Float struct {
	f *big.Float
}

func NewBigFloat(f *big.Float) *Float {
	return &Float{
		f: f,
	}
}

func NewFloat(f float64) *Float {
	return &Float{
		f: big.NewFloat(f),
	}
}

// Marshal implements the gogo proto custom type interface.
func (f Float) Marshal() ([]byte, error) {
	if f.f == nil {
		f.f = new(big.Float)
	}
	return f.f.MarshalText()
}

func (f *Float) MarshalTo(data []byte) (n int, err error) {
	if f.f == nil {
		f.f = new(big.Float)
	}
	if len(f.f.String()) == 0 {
		copy(data, []byte{0x30})
		return 1, nil
	}

	bz, err := f.Marshal()
	if err != nil {
		return 0, err
	}

	copy(data, bz)
	return len(bz), nil
}

// Unmarshal implements the gogo proto custom type interface.
func (f *Float) Unmarshal(data []byte) error {
	if len(data) == 0 {
		f = nil
		return nil
	}

	if f.f == nil {
		f.f = new(big.Float)
	}

	if err := f.f.UnmarshalText(data); err != nil {
		return err
	}

	return nil
}

// Size implements the gogo proto custom type interface.
func (f *Float) Size() int {
	bz, _ := f.Marshal()
	return len(bz)
}

func (f *Float) Cmp(other *Float) int {
	return f.f.Cmp(other.f)
}
func (f *Float) String() string {
	return f.f.String()
}

func (f *Float) Bytes() []byte {
	return []byte(f.f.String())
}
