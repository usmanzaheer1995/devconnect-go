package errors

type formErrors map[string][]string

type ClientError interface {
	Error() string
	ResponseStatus() int
	ResponseData() formErrors
}

type HttpError struct {
	Cause  error               `json:"-"`
	Detail string              `json:"message"`
	Status int                 `json:"-"`
	Data   formErrors `json:"data"`
}

func (e *HttpError) ResponseData() formErrors {
	return e.Data
}

func (e *HttpError) Error() string {
	if e.Detail == "" && e.Cause != nil {
		return e.Cause.Error()
	}
	if e.Cause == nil {
		return e.Detail
	}
	if e.Detail == "" {
		return "Something went wrong"
	}
	return e.Detail + ": " + e.Cause.Error()
}

func (e *HttpError) ResponseStatus() int {
	return e.Status
}

func NewHttpError(err error, status int, detail string) error {
	return &HttpError{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}

func NewHttpError2(status int, detail string, data map[string][]string) error {
	return &HttpError{
		Detail: detail,
		Status: status,
		Data: data,
	}
}
