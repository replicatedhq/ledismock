package ledismock

// Result represents a mock result object to return from the mock db.
type Result struct {
	value interface{}
}

// NewResult constructs a result object.
func NewResult(value interface{}) *Result {
	r := Result{
		value: value,
	}

	return &r
}
