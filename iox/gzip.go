package iox

import (
	"compress/gzip"
	"io"
)

type gzipWriter struct {
	writer *gzip.Writer
}

func NewGzipWriter(w io.Writer) EncodingWriter {
	zw := new(gzipWriter)
	zw.writer = gzip.NewWriter(w)
	return zw
}

func (g *gzipWriter) Write(p []byte) (n int, err error) {
	return g.writer.Write(p)
}

func (g *gzipWriter) ContentEncoding() string {
	return GzipEncoding
}

func (g *gzipWriter) Close() error {
	err0 := g.writer.Flush()
	err1 := g.writer.Close()
	if err0 == nil && err1 == nil {
		return nil
	}
	if err1 != nil {
		return err1
	}
	return err0
}

type gzipReader struct {
	reader *gzip.Reader
}

func NewGzipReader(r io.Reader) (EncodingReader, error) {
	zr := new(gzipReader)
	var err error
	zr.reader, err = gzip.NewReader(r)
	if err != nil {
		return nil, err
	}
	return zr, nil
}

func (g *gzipReader) Read(p []byte) (n int, err error) {
	return g.reader.Read(p)
}

func (g *gzipReader) Close() error {
	return g.reader.Close()
}
