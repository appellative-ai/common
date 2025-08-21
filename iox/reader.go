package iox

import (
	"io"
	"net/http"
)

type EncodingReader interface {
	io.ReadCloser
}

func NewEncodingReader(r io.Reader, h http.Header) (EncodingReader, error) {
	encoding := contentEncoding(h)
	switch encoding {
	case GzipEncoding:
		return NewGzipReader(r)
	case BrotliEncoding, DeflateEncoding, CompressEncoding:
		return nil, newStatusContentEncodingError(encoding)
	default:
		return NewIdentityReader(r), nil
	}
}

type identityReader struct {
	reader io.Reader
}

// NewIdentityReader - The default (identity) encoding; the use of no transformation whatsoever
func NewIdentityReader(r io.Reader) EncodingReader {
	ir := new(identityReader)
	ir.reader = r
	return ir
}

func (i *identityReader) Read(p []byte) (n int, err error) {
	return i.reader.Read(p)
}

func (i *identityReader) Close() error {
	return nil
}
