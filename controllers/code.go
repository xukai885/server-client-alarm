package controllers

type ResCode int

const (
	CodeSuccess ResCode = 1000 + iota
	CodeError
	CodeServerBusy
)

var CodeMsgMap = map[ResCode]string{
	CodeSuccess:    "success",
	CodeError:      "未知错误",
	CodeServerBusy: "系统繁忙",
}

func (c ResCode) Msg() string {
	msg, ok := CodeMsgMap[c]
	if !ok {
		msg = CodeMsgMap[CodeError]
	}
	return msg
}
