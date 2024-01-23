package compressor

import (
	"compress/gzip"
	"io"

	"github.com/gin-gonic/gin"
	covt "github.com/shengyanli1982/orbit-contrib/internal/convertor"
)

// GZipWriter 是一个 GZip 压缩的 ResponseWriter
// GZipWriter is a ResponseWriter of GZip compression
type GZipWriter struct {
	gin.ResponseWriter
	writer *gzip.Writer
}

// NewGZipWriter 创建一个新的 GZipWriter 实例
// NewGZipWriter creates a new GZipWriter instance
func NewGZipWriter(config *Config, rw gin.ResponseWriter) *GZipWriter {
	var gzipWriter *gzip.Writer
	if rw != nil {
		gzipWriter, _ = gzip.NewWriterLevel(rw, config.level)
	} else {
		gzipWriter, _ = gzip.NewWriterLevel(io.Discard, config.level)
	}
	return &GZipWriter{
		ResponseWriter: rw,
		writer:         gzipWriter,
	}
}

// Write 实现了 io.Writer 接口
// Write implements the io.Writer interface
func (gw *GZipWriter) Write(msg []byte) (int, error) {
	gw.Header().Del("Content-Length")
	return gw.writer.Write(msg)
}

// WriteString 将 string 写入到 ResponseWriter
// WriteString writes string to ResponseWriter
func (gw *GZipWriter) WriteString(msg string) (int, error) {
	return gw.Write(covt.StringToBytes(msg))
}

// ResetCompressWriter 重置 GZip 的 io.Writer
// ResetCompressWriter resets the io.Writer of GZip
func (gw *GZipWriter) ResetCompressWriter(w io.Writer) error {
	if w != nil {
		gw.writer.Reset(w)
	}
	return nil
}

// ResetResponseWriter 重置 ResponseWriter
// ResetResponseWriter resets the ResponseWriter
func (gw *GZipWriter) ResetResponseWriter(rw gin.ResponseWriter) error {
	if rw != nil {
		gw.ResponseWriter = rw
	}
	return nil
}

// WriteHeader 设置 ResponseWriter 的状态码
// WriteHeader sets the status code of ResponseWriter
func (gw *GZipWriter) WriteHeader(code int) {
	gw.Header().Del("Content-Length")
	gw.ResponseWriter.WriteHeader(code)
}

// Stop 停止 GZipWriter
// Stop stops GZipWriter
func (gw *GZipWriter) Stop() {
	gw.writer.Close()
}
