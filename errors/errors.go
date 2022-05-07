// Package errors
//  "github.com/pkg/errors" wrapper
package errors

import (
	stderrors "errors"

	"github.com/alkaid/goerrors/internal/errors"
)

var (
	// New https://pkg.go.dev/errors#New 的别名
	//  注意不会像 https://pkg.go.dev/github.com/pkg/errors#New 一样添加堆栈
	//  若要实例化的同时添加堆栈请使用 NewWithStack
	//  常用于 sentinel errors 模式
	New = stderrors.New

	// Is https://pkg.go.dev/errors#Is 的别名
	Is = stderrors.Is

	// As https://pkg.go.dev/errors#As 的别名
	As = stderrors.As

	// Unwrap https://pkg.go.dev/errors#Unwrap 的别名
	Unwrap = stderrors.Unwrap

	// Wrap https://pkg.go.dev/github.com/pkg/errors#Wrap 的别名
	Wrap = errors.Wrap

	// Cause https://pkg.go.dev/github.com/pkg/errors#Cause 的别名
	Cause = errors.Cause

	// WithMessage https://pkg.go.dev/github.com/pkg/errors#WithMessage 的别名
	WithMessage = errors.WithMessage

	// WithStack https://pkg.go.dev/github.com/pkg/errors#WithStack 的别名
	WithStack = errors.WithStack

	// NewWithStack https://pkg.go.dev/github.com/pkg/errors#New 的别名
	//  实例化的同时会添加堆栈信息
	//  若实例化时不添加堆栈请使用 New
	NewWithStack = errors.New
)
