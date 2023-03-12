package tencent

import (
	core "operation/tools/upload"
)

type config struct {
	Url       string
	SecretId  string
	SecretKey string
}

func getConfig() *config {
	return &config{
		Url:       core.Config.Endpoint,
		SecretId:  core.Config.SecretID,
		SecretKey: core.Config.SecretKey,
	}
}
