package adapter

type Session interface {
	Message() string
	UserID() string
	UserName() string
	Platform() string
	ChannelID() string
	GetSession() interface{}
	Log() string
}

type Bot interface {
	Send(session Session, message *Message) error
}

type Segment struct {
	Type string
	Data string
}

type Message struct {
	Segments []*Segment
}

func (message *Message) Text(data string) *Message {
	if message == nil {
		message = &Message{}
	}
	if message.Segments == nil {
		message.Segments = make([]*Segment, 0)
	}
	message.Segments = append(message.Segments, &Segment{
		Type: "text",
		Data: data,
	})
	return message
}

func (message *Message) Image(data string) *Message {
	if message == nil {
		message = &Message{}
	}
	if message.Segments == nil {
		message.Segments = make([]*Segment, 0)
	}
	message.Segments = append(message.Segments, &Segment{
		Type: "image",
		Data: data,
	})
	return message
}
