package storage

type User struct {
	User    string
	Balance float32
}

type Transaction struct {
	MsgID    string
	Sender   string
	Receiver string
	Amount   float32
	Term     int
}
