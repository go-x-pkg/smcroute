package smcroute

import "fmt"

type Message struct {
	Code MessageCode `json:"code,omitempty" yaml:"code"`
	Text string      `json:"text,omitempty" yaml:"text"`
}

func (msg *Message) GetCode() MessageCode { return msg.Code }
func (msg *Message) GetText() string      { return msg.Text }

func (msg *Message) Is(code MessageCode) bool { return msg.Code == code }
func (msg *Message) Error() string            { return msg.String() }
func (msg *Message) String() string           { return msg.Text }

func (msg *Message) SetCode(v MessageCode) *Message { msg.Code = v; return msg }
func (msg *Message) SetText(v string) *Message      { msg.Text = v; return msg }

func NewMessageInitialize(it *Message) {}

func Errorf(code MessageCode, params ...interface{}) *Message {
	text := MessageCodeText(code)
	if len(params) > 0 {
		text = fmt.Sprintf(MessageCodeText(code), params...)
	}

	it := &Message{
		Code: code,
		Text: text,
	}
	NewMessageInitialize(it)
	return it
}

func NewMessage() *Message {
	it := new(Message)
	NewMessageInitialize(it)
	return it
}
