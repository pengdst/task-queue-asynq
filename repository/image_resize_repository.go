package repository

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"task-queue-asynq/model/kraken"
)

type ImageResizeRepository interface {
	Resize(request kraken.Request) (*kraken.Response, error)
}

type ImageResizeRepositoryImpl struct {
	client *resty.Client
}

func (i *ImageResizeRepositoryImpl) Resize(request kraken.Request) (*kraken.Response, error) {
	var response *kraken.Response
	var errResponse *kraken.ErrorResponse
	resp, err := i.client.R().
		SetBody(request).
		SetResult(&response).
		SetError(&errResponse).
		Post("v1/url")

	if err != nil {
		return nil, err
	}

	if !resp.IsSuccess() {
		return nil, errors.New(errResponse.Message)
	}

	return response, nil
}

func NewImageResizeRepository() ImageResizeRepository {
	client := resty.New()
	client.SetBaseURL("https://api.kraken.io")
	return &ImageResizeRepositoryImpl{client: client}
}
