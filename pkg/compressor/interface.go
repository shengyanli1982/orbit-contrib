package compressor

import (
	"io"

	"github.com/gin-gonic/gin"
)

type CodecWriter interface {
	gin.ResponseWriter
	Write(msg []byte) (int, error)
	WriteString(msg string) (int, error)
	WriteHeader(code int)
	ResetCompressWriter(w io.Writer) error
	ResetResponseWriter(rw gin.ResponseWriter) error
	Stop()
}
