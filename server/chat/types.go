package chat

type IncomingMessage struct {
	TokenString	string	`json:"tokenString"`
	Content		string	`json:"content"`
}

type TriggerMediaPayload struct {
	Bucket	string		`json:"bucket"`
	Key		string		`json:"key"`
}
