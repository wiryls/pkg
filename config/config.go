package config

import (
	"bytes"
	"encoding/json"
)

// Marshaler for provide a set of method to Marshal and Unmarshal.
type Marshaler interface {
	CanMarshal(tag string, val interface{}) bool
	Marshal(val interface{}) ([]byte, error)
	CanUnmarshal(tag string, bin []byte) bool
	Unmarshal(bin []byte, val interface{}) error
}

// Marshalers is used to read or write configs from files.
type Marshalers []Marshaler

// export default methods.
var (
	Save           = defaults.Save
	Load           = defaults.Load
	LoadSome       = defaults.LoadSome
	LoadOrSaveSome = defaults.LoadOrSaveSome
)

////////////////////////////// the default //////////////////////////////////

var defaults = Marshalers{
	&Wrapper{
		CanEnc: func(tag string, _ interface{}) bool { return tag == "json" },
		CanDec: func(tag string, ______ []byte) bool { return tag == "json" },
		Decode: json.Unmarshal,
		Encode: func(val interface{}) ([]byte, error) {
			buf := new(bytes.Buffer)
			enc := json.NewEncoder(buf)
			enc.SetIndent("", "    ")
			err := enc.Encode(val)
			bin := buf.Bytes()
			if err != nil {
				bin = nil
			}
			return bin, err
		},
	},
}

// Wrapper provide a simple way to create Marshaler.
type Wrapper struct {
	CanEnc func(tag string, val interface{}) bool
	CanDec func(tag string, bin []byte) bool
	Encode func(val interface{}) ([]byte, error)
	Decode func(bin []byte, val interface{}) error
}

// CanMarshal test if it can marshal val.
func (w *Wrapper) CanMarshal(tag string, val interface{}) bool {
	return w.CanEnc(tag, val)
}

// Marshal object to bytes.
func (w *Wrapper) Marshal(val interface{}) ([]byte, error) {
	return w.Encode(val)
}

// CanUnmarshal test if it can unmarshal.
func (w *Wrapper) CanUnmarshal(tag string, bin []byte) bool {
	return w.CanDec(tag, bin)
}

// Unmarshal bytes to object.
func (w *Wrapper) Unmarshal(bin []byte, val interface{}) error {
	return w.Decode(bin, val)
}
