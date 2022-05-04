package apierrors

import (
	"errors"
	"fmt"
	"io"

	status2 "github.com/alkaid/goerrors/apierrors/http/status"

	pkgerrors "github.com/alkaid/goerrors/internal/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

var Is = errors.Is
var As = errors.As
var Unwrap = errors.Unwrap

//go:generate protoc -I. --go_out=paths=source_relative:. errors.proto

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = 500
	// UnknownReason is unknown reason for error info.
	UnknownReason = ""
	// SupportPackageIsVersion1 this constant should not be referenced by any other code.
	SupportPackageIsVersion1 = true
)

// Error is a status error.
type Error struct {
	Status
	cause error
	*pkgerrors.Stack
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v cause = %v", e.Code, e.Reason, e.Message, e.Metadata, e.cause)
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.cause }
func (e *Error) Cause() error  { return e.cause }

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code && se.Reason == e.Reason
	}
	return false
}

// WithCause with the underlying cause of the error.
func (e *Error) WithCause(cause error) *Error {
	err := Clone(e)
	err.cause = cause
	err.Stack = pkgerrors.Callers()
	return err
}

// WithStack with the stub empty error for stack
func (e *Error) WithStack() *Error {
	err := Clone(e)
	err.Stack = pkgerrors.Callers()
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
func (e *Error) WithMetadata(md map[string]string) *Error {
	err := Clone(e)
	err.Metadata = md
	return err
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(status2.ToGRPCCode(int(e.Code)), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.Reason,
			Metadata: e.Metadata,
		})
	return s
}

func (w *Error) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			if w.cause == nil {
				fmt.Fprintf(s, "%+v", w.Error())
			} else {
				fmt.Fprintf(s, "%+v", w.Cause())
			}
			w.Stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
	}
}

// New returns an error object for the code, message.
func New(code int, reason, message string, pretty string) *Error {
	return &Error{
		Status: Status{
			Code:    int32(code),
			Message: message,
			Reason:  reason,
			Pretty:  pretty,
		},
	}
}

// Code returns the http code for an error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return 200 //nolint:gomnd
	}
	return int(FromError(err).Code)
}

// Reason returns the reason for a particular error.
// It supports wrapped errors.
func Reason(err error) string {
	if err == nil {
		return UnknownReason
	}
	return FromError(err).Reason
}

// Clone deep clone error to a new error.
func Clone(err *Error) *Error {
	metadata := make(map[string]string, len(err.Metadata))
	for k, v := range err.Metadata {
		metadata[k] = v
	}
	return &Error{
		cause: err.cause,
		Status: Status{
			Code:     err.Code,
			Reason:   err.Reason,
			Message:  err.Message,
			Metadata: metadata,
			Pretty:   err.Pretty,
		},
	}
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	gs, ok := status.FromError(err)
	if ok {
		ret := New(
			status2.FromGRPCCode(gs.Code()),
			UnknownReason,
			gs.Message(),
			"",
		)
		for _, detail := range gs.Details() {
			switch d := detail.(type) {
			case *errdetails.ErrorInfo:
				ret.Reason = d.Reason
				return ret.WithMetadata(d.Metadata)
			}
		}
		return ret
	}
	return New(UnknownCode, UnknownReason, err.Error(), "")
}

func FromStatus(status IStatus) *Error {
	metadata := make(map[string]string, len(status.GetMetadata()))
	for k, v := range status.GetMetadata() {
		metadata[k] = v
	}
	return &Error{Status: Status{
		Code:     status.GetCode(),
		Reason:   status.GetReason(),
		Message:  status.GetMessage(),
		Metadata: metadata,
		Pretty:   status.GetPretty(),
	}}
}
