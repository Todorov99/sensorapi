package mailsender

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/Todorov99/server/pkg/dto"
	"github.com/Todorov99/server/pkg/server/config"
	"github.com/go-resty/resty/v2"
)

var mailsenderLogger = logger.NewLogrus("mailsender", os.Stdout)

type Client struct {
	restyClient *resty.Client
}

func New() *Client {
	mailsenderLogger.Debugf("Initializing mail sender client...")

	resty := resty.New().
		SetBaseURL(fmt.Sprintf("http://%s:%s", config.GetMailSenderCfg().GetServiceName(), config.GetMailSenderCfg().GetPort()))

	return &Client{
		restyClient: resty,
	}
}

// SendWithAttachments sends a mail with attachments to a concrete user.
func (c *Client) SendWithAttachments(ctx context.Context, sender dto.MailSenderDto, attachments []string) error {
	mailsenderLogger.Debugf("Sending report file from: %q to: %q", attachments, sender.To)

	senderInfo, err := json.Marshal(sender)
	if err != nil {
		return err
	}

	resp, err := c.restyClient.R().
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

func (c *Client) Send(ctx context.Context, sender dto.MailSenderDto) error {
	mailsenderLogger.Debugf("Sending mail to: %q", sender.To)

	resp, err := c.restyClient.R().
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
