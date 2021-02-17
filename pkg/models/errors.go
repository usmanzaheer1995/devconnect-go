package models

type ClientError interface {
	Error() string
	ResponseStatus() int
}

type HttpError struct {
	Cause  error  `json:"-"`
	Detail string `json:"message"`
	Status int    `json:"-"`
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

func(e *HttpError) ResponseStatus() int {
	return e.Status
}

func NewHttpError(err error, status int, detail string) error {
	return &HttpError{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}
