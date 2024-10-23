package errors

type InternalServerError struct {
	Message    string
	StatusCode int
}

func (h *InternalServerError) Error() string {
	return h.Message
}

type BadRequestError struct {
	Message    string
	StatusCode int
}

func (h *BadRequestError) Error() string {
	return h.Message
}
