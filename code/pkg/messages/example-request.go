package messages

import "github.com/vmihailenco/msgpack/v5"

type ExampleRequest struct {
}

func (exampleRequest *ExampleRequest) Init() *ExampleRequest {
	*exampleRequest = ExampleRequest{}
	return exampleRequest
}

func (exampleRequest *ExampleRequest) Pack() []byte {
	bytes, _ := msgpack.Marshal(&exampleRequest)
	return bytes
}

func (exampleRequest *ExampleRequest) Unpack(data []byte) error {
	return msgpack.Unmarshal(data, &exampleRequest)
}
