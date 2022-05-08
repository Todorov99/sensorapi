package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type MailSenderClient struct {
	logger      *logrus.Entry
	restyClient *resty.Client
}

func NewMailSenderTLSClient(baseUrl, rootCAPemFilePath string) *MailSenderClient {
	resty := resty.New().
		SetBaseURL(baseUrl)

	if rootCAPemFilePath != "" {
		resty = resty.SetRootCertificate(rootCAPemFilePath)
	} else {
		fmt.Println("Skipping TLS verification...")
		resty = resty.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}

	return &MailSenderClient{
		restyClient: resty,
		logger:      logger.NewLogrus("mailsenderclient"),
	}
}

// SendWithAttachments sends a mail with attachments to a concrete user.
func (m *MailSenderClient) SendWithAttachments(ctx context.Context, sender MailSender, attachments []string) error {
	m.logger.Debugf("Sending report file from: %q to: %q", attachments, sender.To)

	senderInfo, err := json.Marshal(sender)
	if err != nil {
		return err
	}

	resp, err := m.restyClient.R().
		SetContext(ctx).
		SetFormDataFromValues(setFormDataFiles(attachments)).
		SetMultipartFields(&resty.MultipartField{
			Param:       "mailInfo",
			FileName:    "",
			ContentType: "application/json",
			Reader:      bytes.NewReader(senderInfo),
		}).Post("/api/mail/attachment/send")

	if err != nil {
		return err
	}

	if respCode := resp.StatusCode(); respCode != http.StatusOK {
		return fmt.Errorf("failed sending metric report: %q", string(resp.Body()))
	}

	return nil
}

func (m *MailSenderClient) Send(ctx context.Context, sender MailSender) error {
	m.logger.Debugf("Sending mail to: %q", sender.To)

	resp, err := m.restyClient.R().
		SetContext(ctx).
		SetBody(sender).
		Post("/api/mail/send")

	if err != nil {
		return err
	}

	if respCode := resp.StatusCode(); respCode != http.StatusOK {
		return fmt.Errorf("failed sending metric report: %q", string(resp.Body()))
	}

	return nil
}

func setFormDataFiles(files []string) url.Values {
	var urlValues url.Values = make(map[string][]string)
	for _, f := range files {
		urlValues.Add("@"+"files", f)
	}

	return urlValues
}
