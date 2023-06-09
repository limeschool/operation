package qiniu

import (
	core "operation/tools/upload"
)

type config struct {
	Bucket    string
	AccessKey string
	SecretKey string
	Zone      string
	Domain    string
	Private   bool
}

func getConfig() *config {
	return &config{
		Bucket:    core.Config.Bucket,
		AccessKey: core.Config.AccessKey,
		SecretKey: core.Config.SecretKey,
		Zone:      core.Config.Zone,
		Domain:    core.Config.Domain,
		Private:   core.Config.Private,
	}
}
