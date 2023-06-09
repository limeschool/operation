package aliyun

import (
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	core "operation/tools/upload"
)

func NewListObjectResult(result oss.ListObjectsResult) *core.ListObjectResult {
	return &core.ListObjectResult{
		Files:      getFiles(result.Objects),
		IsFinished: !result.IsTruncated,
	}
}

func getFiles(objects []oss.ObjectProperties) []core.File {
	var files []core.File

	for _, properties := range objects {
		if strings.HasSuffix(properties.Key, "/") {
			continue
		}

		files = append(files, &File{properties: properties})
	}

	return files
}
