package model

type SendEmailRequest struct {
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	IsHTML      bool
	Attachments []string
}
