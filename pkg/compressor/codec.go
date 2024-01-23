package compressor

import (
	"compress/flate"
	"compress/gzip"
	"io"

	"github.com/gin-gonic/gin"
	covt "github.com/shengyanli1982/orbit-contrib/internal/convertor"
)

const (
	GZipContentEncoding    = "gzip"
	DeflateContentEncoding = "deflate"
)

// =================== GZip ===================

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

// ContentEncoding 返回 Content-Encoding
// ContentEncoding returns Content-Encoding
func (gw *GZipWriter) ContentEncoding() string {
	return GZipContentEncoding
}

// =================== Deflate ===================

// DeflateWriter 是一个 Deflate 压缩的 ResponseWriter
// DeflateWriter is a ResponseWriter of Deflate compression
type DeflateWriter struct {
	gin.ResponseWriter
	writer *flate.Writer
}

// NewDeflateWriter 创建一个新的 DeflateWriter 实例
// NewDeflateWriter creates a new DeflateWriter instance
func NewDeflateWriter(config *Config, rw gin.ResponseWriter) *DeflateWriter {
	var flateWriter *flate.Writer
	if rw != nil {
		flateWriter, _ = flate.NewWriter(rw, config.level)
	} else {
		flateWriter, _ = flate.NewWriter(io.Discard, config.level)
	}
	return &DeflateWriter{
		ResponseWriter: rw,
		writer:         flateWriter,
	}
}

// Write 实现了 io.Writer 接口
// Write implements the io.Writer interface
func (gw *DeflateWriter) Write(msg []byte) (int, error) {
	gw.Header().Del("Content-Length")
	return gw.writer.Write(msg)
}

// WriteString 将 string 写入到 ResponseWriter
// WriteString writes string to ResponseWriter
func (dw *DeflateWriter) WriteString(msg string) (int, error) {
	return dw.Write(covt.StringToBytes(msg))
}

// ResetCompressWriter 重置 GZip 的 io.Writer
// ResetCompressWriter resets the io.Writer of GZip
func (dw *DeflateWriter) ResetCompressWriter(w io.Writer) error {
	if w != nil {
		dw.writer.Reset(w)
	}
	return nil
}

// ResetResponseWriter 重置 ResponseWriter
// ResetResponseWriter resets the ResponseWriter
func (dw *DeflateWriter) ResetResponseWriter(rw gin.ResponseWriter) error {
	if rw != nil {
		dw.ResponseWriter = rw
	}
	return nil
}

// WriteHeader 设置 ResponseWriter 的状态码
// WriteHeader sets the status code of ResponseWriter
func (dw *DeflateWriter) WriteHeader(code int) {
	dw.Header().Del("Content-Length")
	dw.ResponseWriter.WriteHeader(code)
}

// Stop 停止 GZipWriter
// Stop stops GZipWriter
func (dw *DeflateWriter) Stop() {
	dw.writer.Close()
}

// ContentEncoding 返回 Content-Encoding
// ContentEncoding returns Content-Encoding
func (dw *DeflateWriter) ContentEncoding() string {
	return DeflateContentEncoding
}
