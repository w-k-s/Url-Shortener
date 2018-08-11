package error

type Error struct {
	code    Code              `json:"code"`
	domain  string            `json:"domain"`
	message string            `json:"message"`
	fields  map[string]string `json:"fields"`
}

func (e Error) Code() Code {
	return e.code
}

func (e Error) Domain() string {
	return e.domain
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Fields() map[string]string {
	return e.fields
}

func NewError(code Code, domain string, message string, fields map[string]string) *Error {
	return &Error{
		code:    code,
		domain:  domain,
		message: message,
		fields:  fields,
	}
}
