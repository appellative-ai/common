package core

import (
	"fmt"
	"net/http"
)

const (
	StatusInvalidContent             = int(90)  // Content is not available, is nil, or is of the wrong type, usually found via unmarshalling
	StatusIOError                    = int(91)  // I/O operation failed
	StatusJsonDecodeError            = int(92)  // Json decoding failed
	StatusJsonEncodeError            = int(93)  // Json encoding failed
	StatusContentEncodingError       = int(94)  // Content encoding error
	StatusContentEncodingInvalidType = int(95)  // Content encoding error
	StatusNotProvided                = int(96)  // No status is available
	StatusRateLimited                = int(97)  // Rate limited
	StatusNotStarted                 = int(98)  // Not started
	StatusHaveContent                = int(99)  // Content is available
	StatusGzipEncodingError          = int(100) // Gzip encoding error
	StatusGzipDecodingError          = int(101) // Gzip decoding error
	StatusExecError                  = int(105) // Execution error, as in a database call
	StatusInvalidArgument            = 3        //codes.InvalidArgument    // The client specified an invalid argument. Note that this differs from FAILED_PRECONDITION. INVALID_ARGUMENT indicates arguments that are problematic regardless of the state of the system (e.g., a malformed file name).
	StatusDeadlineExceeded           = 4        //codes.DeadlineExceeded   // The deadline expired before the operation could complete. For operations that change the state of the system, this error may be returned even if the operation has completed successfully. For example, a successful response from a server could have been delayed long
)

var (
	StatusOK       = NewStatus(http.StatusOK, nil)
	StatusNotFound = NewStatus(http.StatusNotFound, nil)
)

type Status struct {
	Err      error
	Code     int
	Location string
}

func NewStatus(code int, err error) *Status {
	s := new(Status)
	s.Code = code
	s.Err = err
	return s
}

func (s *Status) OK() bool       { return s.Code == http.StatusOK }
func (s *Status) NotFound() bool { return s.Code == http.StatusNotFound }
func (s *Status) String() string {
	if s.Err != nil {
		return fmt.Sprintf("%v - %v", HttpStatus(s.Code), s.Err)
	}
	return fmt.Sprintf("%v", HttpStatus(s.Code))

}

func (s *Status) SetLocation(location string) *Status {
	s.Location = location
	return s
}

/*
func (s *Status) Error() string {
	if s.Err != nil {
		return s.Err.Error()
	}
	return ""
}


*/
