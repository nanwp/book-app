package helper

type ErrBadRequest struct {
	Message string
}

func NewErrBadRequest(message string) *ErrBadRequest {
	return &ErrBadRequest{Message: message}
}

func (e ErrBadRequest) Error() string {
	return e.Message
}

type ErrNotFound struct {
	Message string
}

func NewErrNotFound(message string) *ErrNotFound {
	return &ErrNotFound{Message: message}
}

func (e ErrNotFound) Error() string {
	return e.Message
}

type ErrUnauthorized struct {
	Message string
}

func (e ErrUnauthorized) Error() string {
	return e.Message
}

func NewErrUnauthorized(message string) *ErrUnauthorized {
	return &ErrUnauthorized{Message: message}
}

type ErrForbidden struct {
	Message string
}

func NewErrForbidden(message string) *ErrForbidden {
	return &ErrForbidden{Message: message}
}

func (e ErrForbidden) Error() string {
	return e.Message
}

type ErrInternalServer struct {
	Message string
}

func (e ErrInternalServer) Error() string {
	return e.Message
}

func NewErrInternalServer(message string) error {
	return &ErrInternalServer{Message: message}
}

func IsErrNotFound(err error) bool {
	_, ok := err.(*ErrNotFound)
	return ok
}
