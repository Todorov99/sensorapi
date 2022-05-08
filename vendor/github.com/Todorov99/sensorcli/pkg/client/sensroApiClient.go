package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/Todorov99/sensorcli/pkg/logger"
	"github.com/Todorov99/sensorcli/pkg/sensor"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
)

type HttpResponseError struct {
	Err        error
	StatusCode int
}

type APIClient struct {
	logger    *logrus.Entry
	restyClt  *resty.Client
	token     string
	isExpered bool
}

func NewAPIClient(ctx context.Context, baseURL, rootCAPemFilePath string) *APIClient {
	var token string
	restyClt := resty.New().SetBaseURL(baseURL)

	if rootCAPemFilePath != "" {
		restyClt = restyClt.SetRootCertificate(rootCAPemFilePath)
	} else {
		restyClt = restyClt.SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	}

	return &APIClient{
		logger:    logger.NewLogrus("apiclient"),
		restyClt:  restyClt,
		isExpered: true,
		token:     token,
	}
}

func (a APIClient) GetToken(ctx context.Context, username, password string) (string, HttpResponseError) {
	usr := UserDto{
		Username: username,
		Password: password,
	}

	resp, err := a.restyClt.R().
		SetBasicAuth(username, password).
		SetContext(ctx).
		SetBody(usr).Get("/api/users/login")

	if err != nil {
		return "", HttpResponseError{
			Err: err,
		}
	}

	if resp.StatusCode() != http.StatusOK {
		return "", HttpResponseError{
			Err:        fmt.Errorf(string(resp.Body())),
			StatusCode: resp.StatusCode(),
		}
	}

	return resp.Header().Get("Token"), HttpResponseError{
		Err:        nil,
		StatusCode: resp.StatusCode(),
	}
}

func (a APIClient) SendMetrics(ctx context.Context, username, password string, measurements sensor.Measurment) HttpResponseError {
	a.logger.Debug("Sending metrics...")
	if a.isExpered || a.token == "" {
		a.logger.Debug("Getting new user token...")
		token, respError := a.GetToken(ctx, username, password)
		if respError.Err != nil {
			return respError
		}

		a.token = token
		a.isExpered = false
		a.logger.Debug("Token successfully updated")
	}

	resp, err := a.restyClt.R().
		SetContext(ctx).
		SetAuthToken(a.token).
		SetBody(measurements).
		Post("/api/measurement")

	if err != nil {
		return HttpResponseError{
			Err: err,
		}
	}

	if resp.StatusCode() != http.StatusOK {
		if resp.StatusCode() == http.StatusForbidden && strings.Contains(err.Error(), "is expired") {
			a.isExpered = true
			return a.SendMetrics(ctx, username, password, measurements)
		}

		return HttpResponseError{
			Err: fmt.Errorf(string(resp.Body())),
		}
	}

	return HttpResponseError{
		Err:        nil,
		StatusCode: resp.StatusCode(),
	}
}
