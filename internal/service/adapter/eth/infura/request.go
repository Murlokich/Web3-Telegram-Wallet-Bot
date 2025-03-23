package infura

type Request struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

func NewRequest(method string, params interface{}) Request {
	return Request{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
}

type Response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (r *Response) HasError() bool {
	return r.Error != nil
}
