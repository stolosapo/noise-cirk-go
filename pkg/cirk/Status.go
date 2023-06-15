package cirk

// Status represent the status of the circuit breaker
type Status string

const (
	// Closed state is when the breaker works normal and allow all requests
	Closed = "closed"

	// Opened state is when the breaker has problem and not allowing any request
	Opened = "opened"

	// HalfOpened state is when the breaker has problem but allowing only one small percentage of requests
	HalfOpened = "half opened"
)
