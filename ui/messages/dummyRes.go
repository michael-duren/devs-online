package messages

type DummyResponse struct {
	StatusCode int
	Message    *string
	Err        error
}
