package apierrors

type IStatus interface {
	GetCode() int32

	GetReason() string

	GetMessage() string

	GetMetadata() map[string]string

	GetPretty() string
}
