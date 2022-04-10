package dto

type MailSenderDto struct {
	Subject string   `json:"subject,omitempty"`
	Cc      []string `json:"cc,omitempty"`
	To      []string `json:"to,omitempty"`
	Body    string   `json:"body,omitempty"`
}
