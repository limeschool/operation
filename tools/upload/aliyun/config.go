package aliyun

import (
	core "operation/tools/upload"
)

type config struct {
	Endpoint        string
	Bucket          string
	AccessKeyID     string
	AccessKeySecret string
}

func getConfig() *config {
	return &config{
		Endpoint:        core.Config.Endpoint,
		Bucket:          core.Config.Bucket,
		AccessKeyID:     core.Config.AccessKeyID,
		AccessKeySecret: core.Config.AccessKeySecret,
	}
}
