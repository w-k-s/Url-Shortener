package error

type Code int

type Err interface {
	Domain() string
	Code() Code
	Error() string
	Fields() map[string]string
}
