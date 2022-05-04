// Package errors
//  "github.com/pkg/errors" wrapper
package errors

import (
	stderrors "errors"

	"github.com/alkaid/goerrors/internal/errors"
)

var (
	New         = stderrors.New
	Is          = stderrors.Is
	As          = stderrors.As
	Unwrap      = stderrors.Unwrap
	Wrap        = errors.Wrap
	Cause       = errors.Cause
	WithMessage = errors.WithMessage
	WithStack   = errors.WithStack
)
