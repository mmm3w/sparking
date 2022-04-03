package wsk

type CacheMsg struct {
	MsgId    string `json:"id"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}

func getUnreadMessages(p1 string, p2 string) []*CacheMsg {
	data := make([]*CacheMsg, 0)
	// data = append(data, &CacheMsg{
	// 	MsgId:    "asdfasdf",
	// 	Sender:   p1,
	// 	Receiver: p2,
	// 	Type:     "text",
	// 	Content:  "test",
	// })
	return data
}
