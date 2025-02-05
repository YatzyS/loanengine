package rest

type RestHandler interface{}

type restHandler struct{}

func NewRestHandler() RestHandler {
	return &restHandler{}
}