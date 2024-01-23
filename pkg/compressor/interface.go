package compressor

import (
	"io"

	"github.com/gin-gonic/gin"
)

// CodecWriter 是一个接口，定义了压缩编码器的写入操作。
// CodecWriter is an interface that defines the write operations of a compression encoder.
type CodecWriter interface {
	gin.ResponseWriter

	// Write 将字节切片写入响应。返回写入的字节数和可能的错误。
	// Write writes a byte slice to the response. Returns the number of bytes written and the possible error.
	Write(msg []byte) (int, error)

	// WriteString 将字符串写入响应。返回写入的字节数和可能的错误。
	// WriteString writes a string to the response. Returns the number of bytes written and the possible error.
	WriteString(msg string) (int, error)

	// WriteHeader 设置响应的状态码
	// WriteHeader sets the status code of the response
	WriteHeader(code int)

	// ResetCompressWriter 重置压缩编码器的写入器。	参数 w 是新的写入器，返回可能的错误。
	// ResetCompressWriter resets the writer of the compression encoder. The parameter w is the new writer, and the possible error is returned.
	ResetCompressWriter(w io.Writer) error

	// ResetResponseWriter 重置响应写入器。参数 rw 是新的响应写入器，返回可能的错误。
	// ResetResponseWriter resets the response writer. The parameter rw is the new response writer, and the possible error is returned.
	ResetResponseWriter(rw gin.ResponseWriter) error

	// ContentEncoding 返回响应的 Content-Encoding
	// ContentEncoding returns the Content-Encoding of the response
	ContentEncoding() string

	// Stop 停止压缩编码器的操作
	// Stop stops the operation of the compression encoder
	Stop()
}
