package tencent

import (
	"net/http"
	"net/url"

	"github.com/tencentyun/cos-go-sdk-v5"
	core "operation/tools/upload"
)

type Driver struct{}

func NewDriver() core.Driver {
	return &Driver{}
}

func (d *Driver) Storage() (core.Storage, error) {
	conf := getConfig()

	if conf.Url == "" || conf.SecretId == "" || conf.SecretKey == "" {
		return nil, core.ErrorConfigEmpty
	}

	u, err := url.Parse(conf.Url)
	if err != nil {
		return nil, err
	}

	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.SecretId,
			SecretKey: conf.SecretKey,
		},
	})

	store := Store{
		client: client,
	}

	return core.NewStorage(&store), nil
}

func (d Driver) Name() string {
	return "tencent"
}
