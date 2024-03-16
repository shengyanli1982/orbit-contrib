package compressor

import (
	"compress/flate"
	"compress/gzip"
	"io"

	"github.com/gin-gonic/gin"
	covt "github.com/shengyanli1982/orbit-contrib/internal/convertor"
)

// 定义 GZip 和 Deflate 的内容编码常量
// Define content encoding constants for GZip and Deflate
const (
	// GZip 内容编码
	// GZip content encoding
	GZipContentEncoding = "gzip"

	// Deflate 内容编码
	// Deflate content encoding
	DeflateContentEncoding = "deflate"
)

// GZipWriter 是一个 GZip 压缩的 ResponseWriter
// GZipWriter is a ResponseWriter for GZip compression
type GZipWriter struct {
	// 继承 gin 的 ResponseWriter
	// Inherits gin's ResponseWriter
	gin.ResponseWriter

	// GZip 压缩写入器
	// GZip compression writer
	writer *gzip.Writer
}

// NewGZipWriter 创建一个新的 GZipWriter 实例
// NewGZipWriter creates a new GZipWriter instance
func NewGZipWriter(config *Config, rw gin.ResponseWriter) *GZipWriter {
	// 定义一个 GZip 写入器
	// Define a GZip writer
	var gzipWriter *gzip.Writer

	// 如果 ResponseWriter 不为空，则创建一个新的 GZip 写入器，否则创建一个写入到 io.Discard 的 GZip 写入器
	// If ResponseWriter is not null, create a new GZip writer, otherwise create a GZip writer that writes to io.Discard
	if rw != nil {
		// 创建一个新的 GZip 写入器
		// Create a new GZip writer
		gzipWriter, _ = gzip.NewWriterLevel(rw, config.level)
	} else {
		// 创建一个写入到 io.Discard 的 GZip 写入器
		// Create a GZip writer that writes to io.Discard
		gzipWriter, _ = gzip.NewWriterLevel(io.Discard, config.level)
	}

	// 返回一个新的 GZipWriter 实例
	// Return a new GZipWriter instance
	return &GZipWriter{
		// 设置 ResponseWriter
		// Set ResponseWriter
		ResponseWriter: rw,

		// 设置 GZip 写入器
		// Set GZip writer
		writer: gzipWriter,
	}
}

// GZipWriter 的 Write 方法，删除 "Content-Length" 头部，然后写入消息
// Write method of GZipWriter, deletes the "Content-Length" header, then writes the message
func (gw *GZipWriter) Write(msg []byte) (int, error) {
	// 删除 "Content-Length" 头部
	// Deletes the "Content-Length" header
	gw.Header().Del("Content-Length")

	// 写入消息
	// Writes the message
	return gw.writer.Write(msg)
}

// GZipWriter 的 WriteString 方法，将字符串转换为字节并写入
// WriteString method of GZipWriter, converts the string to bytes and writes it
func (gw *GZipWriter) WriteString(msg string) (int, error) {
	return gw.Write(covt.StringToBytes(msg))
}

// GZipWriter 的 ResetCompressWriter 方法，重置压缩写入器
// ResetCompressWriter method of GZipWriter, resets the compression writer
func (gw *GZipWriter) ResetCompressWriter(w io.Writer) error {
	// 如果写入器不为空，则重置写入器
	// If the writer is not null, reset the writer
	if w != nil {
		gw.writer.Reset(w)
	}

	// 返回 nil 表示没有错误
	// Returns nil indicating no error
	return nil
}

// GZipWriter 的 ResetResponseWriter 方法，重置响应写入器
// ResetResponseWriter method of GZipWriter, resets the response writer
func (gw *GZipWriter) ResetResponseWriter(rw gin.ResponseWriter) error {
	// 如果响应写入器不为空，则重置响应写入器
	// If the response writer is not null, reset the response writer
	if rw != nil {
		gw.ResponseWriter = rw
	}

	// 返回 nil 表示没有错误
	// Returns nil indicating no error
	return nil
}

// GZipWriter 的 WriteHeader 方法，删除 "Content-Length" 头部，然后写入状态码
// WriteHeader method of GZipWriter, deletes the "Content-Length" header, then writes the status code
func (gw *GZipWriter) WriteHeader(code int) {
	// 删除 "Content-Length" 头部
	// Deletes the "Content-Length" header
	gw.Header().Del("Content-Length")

	// 写入状态码
	// Writes the status code
	gw.ResponseWriter.WriteHeader(code)
}

// GZipWriter 的 Stop 方法，关闭写入器
// Stop method of GZipWriter, close the writer
func (gw *GZipWriter) Stop() {
	gw.writer.Close()
}

// GZipWriter 的 ContentEncoding 方法，返回内容编码
// ContentEncoding method of GZipWriter, return the content encoding
func (gw *GZipWriter) ContentEncoding() string {
	return GZipContentEncoding
}

// DeflateWriter 是一个 Deflate 压缩的 ResponseWriter
// DeflateWriter is a ResponseWriter for Deflate compression
type DeflateWriter struct {
	// 继承 gin 的 ResponseWriter
	// Inherits gin's ResponseWriter
	gin.ResponseWriter

	// Deflate 压缩写入器
	// Deflate compression writer
	writer *flate.Writer
}

// NewDeflateWriter 创建一个新的 DeflateWriter 实例
// NewDeflateWriter creates a new DeflateWriter instance
func NewDeflateWriter(config *Config, rw gin.ResponseWriter) *DeflateWriter {
	// 定义一个 Deflate 写入器
	// Define a Deflate writer
	var flateWriter *flate.Writer

	// 如果 ResponseWriter 不为空，则创建一个新的 Deflate 写入器，否则创建一个写入到 io.Discard 的 Deflate 写入器
	// If ResponseWriter is not null, create a new Deflate writer, otherwise create a Deflate writer that writes to io.Discard
	if rw != nil {
		// 创建一个新的 Deflate 写入器
		// Create a new Deflate writer
		flateWriter, _ = flate.NewWriter(rw, config.level)
	} else {
		// 创建一个写入到 io.Discard 的 Deflate 写入器
		// Create a Deflate writer that writes to io.Discard
		flateWriter, _ = flate.NewWriter(io.Discard, config.level)
	}

	// 返回一个新的 DeflateWriter 实例
	// Return a new DeflateWriter instance
	return &DeflateWriter{
		// 设置 ResponseWriter
		// Set ResponseWriter
		ResponseWriter: rw,

		// 设置 Deflate 写入器
		// Set Deflate writer
		writer: flateWriter,
	}
}

// DeflateWriter 的 Write 方法，删除 "Content-Length" 头部，然后写入消息
// Write method of DeflateWriter, deletes the "Content-Length" header, then writes the message
func (dw *DeflateWriter) Write(msg []byte) (int, error) {
	// 删除 "Content-Length" 头部
	// Deletes the "Content-Length" header
	dw.Header().Del("Content-Length")

	// 写入消息
	// Writes the message
	return dw.writer.Write(msg)
}

// DeflateWriter 的 WriteString 方法，将字符串转换为字节并写入
// WriteString method of DeflateWriter, converts the string to bytes and writes it
func (dw *DeflateWriter) WriteString(msg string) (int, error) {
	// 将字符串转换为字节并写入
	// Converts the string to bytes and writes it
	return dw.Write(covt.StringToBytes(msg))
}

// DeflateWriter 的 ResetCompressWriter 方法，重置压缩写入器
// ResetCompressWriter method of DeflateWriter, resets the compression writer
func (dw *DeflateWriter) ResetCompressWriter(w io.Writer) error {
	// 如果写入器不为空，则重置写入器
	// If the writer is not null, reset the writer
	if w != nil {
		dw.writer.Reset(w)
	}

	// 返回 nil 表示没有错误
	// Returns nil indicating no error
	return nil
}

// DeflateWriter 的 ResetResponseWriter 方法，重置响应写入器
// ResetResponseWriter method of DeflateWriter, resets the response writer
func (dw *DeflateWriter) ResetResponseWriter(rw gin.ResponseWriter) error {
	// 如果响应写入器不为空，则重置响应写入器
	// If the response writer is not null, reset the response writer
	if rw != nil {
		dw.ResponseWriter = rw
	}

	// 返回 nil 表示没有错误
	// Returns nil indicating no error
	return nil
}

// DeflateWriter 的 WriteHeader 方法，删除 "Content-Length" 头部，然后写入状态码
// WriteHeader method of DeflateWriter, deletes the "Content-Length" header, then writes the status code
func (dw *DeflateWriter) WriteHeader(code int) {
	// 删除 "Content-Length" 头部
	// Deletes the "Content-Length" header
	dw.Header().Del("Content-Length")

	// 写入状态码
	// Writes the status code
	dw.ResponseWriter.WriteHeader(code)
}

// DeflateWriter 的 Stop 方法，关闭写入器
// Stop method of DeflateWriter, closes the writer
func (dw *DeflateWriter) Stop() {
	// 关闭写入器
	// Closes the writer
	dw.writer.Close()
}

// DeflateWriter 的 ContentEncoding 方法，返回内容编码
// ContentEncoding method of DeflateWriter, returns the content encoding
func (dw *DeflateWriter) ContentEncoding() string {
	// 返回内容编码
	// Returns the content encoding
	return DeflateContentEncoding
}
