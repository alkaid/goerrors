package apierrors

import (
	"errors"
	"fmt"
	"io"

	status2 "github.com/alkaid/goerrors/apierrors/http/status"

	pkgerrors "github.com/alkaid/goerrors/errors"
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

type errKey string

var errs = map[errKey]*Error{}

// Register 注册错误信息
func Register(e *Error) {
	errs[errKey(e.Reason)] = e
}

// Error is a status error.
type Error struct {
	Status
	cause error
}

func (e *Error) Error() string {
	if e.cause != nil {
		return e.fmtError() + e.cause.Error()
	}
	return e.fmtError()
}
func (e *Error) fmtError() string {
	return fmt.Sprintf("code=%d,reason=%s,message=%s,metadata=%v,cause:%s", e.Code, e.Reason, e.Message, e.Metadata, e.cause.Error())
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
//
//	若cause是包装过stack的会保留stack,打印时会打印cause的堆栈
func (e *Error) WithCause(cause error) *Error {
	err := Clone(e)
	err.cause = cause
	return err
}

// WithStack with the stub empty error for stack
//
//	注意 FromError 会自动调用 WithCause,若已经 WithCause 且 cause 已经包装过stack的请勿重复调用 WithStack
func (e *Error) WithStack() error {
	return pkgerrors.WithStack(e)
}

// WithMessage set message to current Error
//
//	注意不会添加stack
func (e *Error) WithMessage(msg string) *Error {
	err := Clone(e)
	err.Message = msg
	return err
}

// WithAppend 在原本消息后添加 fmt.Sprintf
//
//	注意不会添加stack
func (e *Error) WithAppend(format string, a ...any) *Error {
	err := Clone(e)
	err.Message = fmt.Sprintf(err.Message+". "+format, a...)
	return err
}

// WithTail 在原本消息前添加 fmt.Sprintf
//
//	注意不会添加stack
func (e *Error) WithTail(format string, a ...any) *Error {
	err := Clone(e)
	err.Message = fmt.Sprintf(format+". err="+e.Message, a...)
	return err
}

// WithPretty set pretty to current Error
//
//	注意不会添加stack
func (e *Error) WithPretty(pretty string) *Error {
	err := Clone(e)
	err.Pretty = pretty
	return err
}

// WithMetadata with an MD formed by the mapping of key, value.
//
//	注意不会添加stack
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
			io.WriteString(s, w.fmtError())
			fmt.Fprintf(s, "%+v\n", w.Cause())
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

// New returns an error object for the code, message.
//
//	不带 stack
func New(code int, reason string, message string, pretty string) *Error {
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
//
//	1.若原 error 带stack,则会保留stack
//	2.若原 error 是 Wrap 过的 Error,则可能丢失 stack
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
				e, ok := errs[errKey(d.Reason)]
				if ok {
					return e.WithMessage(gs.Message()).WithMetadata(d.Metadata)
				}
				return New(
					status2.FromGRPCCode(gs.Code()),
					d.Reason,
					gs.Message(),
					"",
				).WithMetadata(d.Metadata).WithCause(err)
			default:
				// do nothing
			}
		}
		return ret
	}
	return New(UnknownCode, UnknownReason, err.Error(), "").WithCause(err)
}

// FromStatus 将 IStatus 转为 Error 并 Error.WithStack
//
//	@param status
//	@return error
func FromStatus(status IStatus) error {
	metadata := make(map[string]string, len(status.GetMetadata()))
	for k, v := range status.GetMetadata() {
		metadata[k] = v
	}
	e := &Error{Status: Status{
		Code:     status.GetCode(),
		Reason:   status.GetReason(),
		Message:  status.GetMessage(),
		Metadata: metadata,
		Pretty:   status.GetPretty(),
	}}
	return e.WithStack()
}

// FromStatusWithoutStack 将 IStatus 转为 Error, 不带 stack
//
//	@param status
//	@return *Error
func FromStatusWithoutStack(status IStatus) *Error {
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
