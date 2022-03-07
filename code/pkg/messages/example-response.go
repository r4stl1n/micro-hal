package messages

import "github.com/vmihailenco/msgpack/v5"

type ExampleResponse struct {
	Result
}

func (exampleResponse *ExampleResponse) Success() *ExampleResponse {

	*exampleResponse = ExampleResponse{
		Result: Result{
			Type: SuccessResult,
			Text: "",
		},
	}

	return exampleResponse
}

func (exampleResponse *ExampleResponse) Failure(reason string) *ExampleResponse {

	*exampleResponse = ExampleResponse{
		Result: Result{
			Type: FailureResult,
			Text: reason,
		},
	}

	return exampleResponse
}

func (exampleResponse *ExampleResponse) Pack() []byte {
	bytes, _ := msgpack.Marshal(&exampleResponse)
	return bytes
}

func (exampleResponse *ExampleResponse) Unpack(data []byte) error {
	return msgpack.Unmarshal(data, &exampleResponse)
}
