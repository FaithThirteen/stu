package codec

import "io"

type Header struct {
	ServiceMethod string // 格式为 "Service.Method"
	Seq           uint64 // 相当于请求的id
	Error         string
}

// Codec 抽象的编码器接口，可实现不同的编解码
type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

type NewCodecFunc func(io.ReadWriteCloser) Codec

type Type string

const (
	GobType  Type = "application/gob"
	JsonType Type = "application/json"
)

var NewCodecFuncMap map[Type]NewCodecFunc

// 类似工厂模式，不过是返回创建方法，而不是实例
func init() {
	NewCodecFuncMap = make(map[Type]NewCodecFunc)
	NewCodecFuncMap[GobType] = NewGobCodec
}
