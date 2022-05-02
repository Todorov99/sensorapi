package client

type UserDto struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type MailSender struct {
	Subject string   `json:"subject,omitempty"`
	Cc      []string `json:"cc,omitempty"`
	To      []string `json:"to,omitempty"`
	Body    string   `json:"body,omitempty"`
}
