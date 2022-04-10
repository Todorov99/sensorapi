package config

import (
	"os"
	"strings"
)

type mailSender struct {
	serviceName string
	port        string
}

func NewMailSender(applicationProperties *ApplicationProperties) *mailSender {
	port := applicationProperties.MailSender.Port

	if strings.HasPrefix(port, "${") {
		str := strings.TrimPrefix(port, "${")
		str = strings.TrimSuffix(str, "}")
		port = os.Getenv(str)
	}

	return &mailSender{
		serviceName: applicationProperties.MailSender.ServiceName,
		port:        port,
	}
}

func (m *mailSender) GetServiceName() string {
	return m.serviceName
}

func (m *mailSender) GetPort() string {
	return m.port
}
