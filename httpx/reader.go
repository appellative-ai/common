package httpx

import "io"

var (
	EmptyReader = new(emptyReader)
)

type emptyReader struct{}

func (r *emptyReader) Read(p []byte) (int, error) {
	return 0, io.EOF
}

func (r *emptyReader) Close() error {
	return nil
}

func readAll(body io.Reader) ([]byte, error) {
	if body == nil {
		return nil, nil
	}
	return io.ReadAll(body)
}
