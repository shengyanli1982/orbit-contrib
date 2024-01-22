package compressor

import (
	"io"

	"github.com/gin-gonic/gin"
)

type CodecWriter interface {
	gin.ResponseWriter
	Write(msg []byte) (int, error)
	Reset(w io.Writer) error
	Stop()
}
