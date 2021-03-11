package request

type Codec interface {
	Marshaller
	Unmarshaller
}

type Marshaller interface {
	Marshal(v interface{}) ([]byte, error)
}

type Unmarshaller interface {
	Unmarshal(data []byte, v interface{}) error
}
type CodecFuncs struct {
	MarshalFunc
	UnmarshalFunc
}

type MarshalFunc func(v interface{}) ([]byte, error)

func (f MarshalFunc) Marshal(v interface{}) ([]byte, error) { return f(v) }

type UnmarshalFunc func(data []byte, v interface{}) error

func (f UnmarshalFunc) Unmarshal(data []byte, v interface{}) error { return f(data, v) }

func MarshalFuncOf(m Marshaller) MarshalFunc {
	switch mm := m.(type) {
	case CodecFuncs:
		return mm.MarshalFunc
	case MarshalFunc:
		return mm
	}
	return m.Marshal
}

func UnmarshalFuncOf(u Unmarshaller) UnmarshalFunc {
	switch uu := u.(type) {
	case CodecFuncs:
		return uu.UnmarshalFunc
	case UnmarshalFunc:
		return uu
	}
	return u.Unmarshal
}
