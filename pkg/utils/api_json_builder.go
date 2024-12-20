package utils

// JSONResult -
type JSONResult struct {
	Code   Code        `json:"code"`
	Msg    Msg         `json:"msg"`
	Detail string      `json:"detail,omitempty"`
	Data   interface{} `json:"data"`
}

// SuccessJSONData -
func SuccessJSONData(data interface{}) *JSONResult {
	return &JSONResult{
		Code: Success,
		Msg:  SuccessMsg,
		Data: data,
	}
}

// FailJSONData -
func FailJSONData(code Code, msg Msg, err error) *JSONResult {
	var detail string
	if err != nil {
		detail = err.Error()
	}
	return &JSONResult{
		Code:   code,
		Msg:    msg,
		Detail: detail,
		Data:   nil,
	}
}
