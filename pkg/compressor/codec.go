package compressor

import (
	"compress/gzip"
	"io"
	"sync"

	"github.com/gin-gonic/gin"
	covt "github.com/shengyanli1982/orbit-contrib/internal/convertor"
)

type GZipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
	once   sync.Once
}

func NewGZipWriter(config *Config, rw gin.ResponseWriter) *GZipWriter {
	z, _ := gzip.NewWriterLevel(io.Discard, config.level)
	writer := GZipWriter{
		ResponseWriter: rw,
		writer:         z,
		once:           sync.Once{},
	}
	return &writer
}

func (gw *GZipWriter) WriteHeader(code int) {
	gw.ResponseWriter.WriteHeader(code)
}

func (gw *GZipWriter) Write(msg []byte) (int, error) {
	return gw.writer.Write(msg)
}

func (gw *GZipWriter) WriteString(msg string) (int, error) {
	return gw.Write(covt.StringToBytes(msg))
}

func (gw *GZipWriter) Reset(w io.Writer) error {
	err := gw.writer.Flush()
	defer gw.writer.Reset(w)
	if err != nil {
		return err
	}
	return nil
}

func (gw *GZipWriter) Stop() {
	gw.once.Do(func() {
		gw.writer.Close()
	})
}
