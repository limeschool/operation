package minio

import (
	"strings"

	"github.com/minio/minio-go/v7"

	core "operation/tools/upload"
)

func NewListObjectResult(entries []minio.ObjectInfo, hasNext bool) *core.ListObjectResult {
	return &core.ListObjectResult{
		Files:      getFiles(entries),
		IsFinished: !hasNext,
	}
}

func getFiles(objects []minio.ObjectInfo) []core.File {
	var files []core.File

	for _, item := range objects {
		if strings.HasSuffix(item.Key, "/") {
			continue
		}

		files = append(files, &File{obj: item})
	}

	return files
}
