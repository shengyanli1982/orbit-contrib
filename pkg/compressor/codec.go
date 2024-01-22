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
	gzipWriter, _ := gzip.NewWriterLevel(rw, config.level)
	writer := GZipWriter{
		ResponseWriter: rw,
		writer:         gzipWriter,
	}
	return &writer
}

// Write 实现了 io.Writer 接口
// Write implements the io.Writer interface
func (gw *GZipWriter) Write(msg []byte) (int, error) {
	return gw.writer.Write(msg)
}

// WriteString 将 string 写入到 ResponseWriter
// WriteString writes string to ResponseWriter
func (gw *GZipWriter) WriteString(msg string) (int, error) {
	return gw.Write(covt.StringToBytes(msg))
}

// Reset 重置 ResponseWriter 的 io.Writer, 更换前将原有数据全部写入底层 ResponseWriter
// Reset resets the io.Writer of ResponseWriter, write all data to underlying ResponseWriter before replacing
func (gw *GZipWriter) Reset(w io.Writer) error {
	err := gw.writer.Flush()
	defer gw.writer.Reset(w)
	if err != nil {
		return err
	}
	return nil
}

// Stop 停止 GZipWriter
// Stop stops GZipWriter
func (gw *GZipWriter) Stop() {
	gw.writer.Close()
}
