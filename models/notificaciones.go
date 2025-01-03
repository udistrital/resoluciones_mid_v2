package models

type Email struct {
	Destination Destination `json:"Destination"`
	Message     Message     `json:"Message"`
	SourceEmail string      `json:"SourceEmail"`
	SourceName  string      `json:"SourceName"`
}

type Destination struct {
	ToAddresses []string `json:"ToAddresses"`
}

type Content struct {
	Data string `json:"Data"`
}

type Body struct {
	Html Content `json:"Html"`
	Text Content `json:"Text"`
}

type Message struct {
	Body        Body         `json:"Body"`
	Subject     Content      `json:"Subject"`
	Attachments []Attachment `json:"Attachments"`
}

type Attachment struct {
	ContentType string `json:"ContentType"`
	FileName    string `json:"FileName"`
	Base64File  string `json:"Base64File"`
}

type EmailResponse struct {
	Result struct {
		MessageId string `json:"MessageId"`
	} `json:"Result"`
}
