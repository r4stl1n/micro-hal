package messages

import "github.com/vmihailenco/msgpack/v5"

type ResultType int

const (
	SuccessResult ResultType = 1
	FailureResult ResultType = 2

	UnknownResult ResultType = 99
)

type Result struct {
	Type ResultType
	Text string
}

func (result *Result) Init(rType ResultType, text string) *Result {

	*result = Result{
		Type: rType,
		Text: text,
	}

	return result
}

func (result *Result) Pack() []byte {
	bytes, _ := msgpack.Marshal(&result)
	return bytes
}

func (result *Result) Unpack(data []byte) error {
	return msgpack.Unmarshal(data, &result)
}
