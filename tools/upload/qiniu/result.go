package qiniu

import (
	"strings"

	"github.com/qiniu/go-sdk/v7/storage"
	core "operation/tools/upload"
)

func NewListObjectResult(entries []storage.ListItem, hasNext bool) *core.ListObjectResult {
	return &core.ListObjectResult{
		Files:      getFiles(entries),
		IsFinished: !hasNext,
	}
}

func getFiles(entries []storage.ListItem) []core.File {
	var files []core.File

	for _, item := range entries {
		if strings.HasSuffix(item.Key, "/") {
			continue
		}

		files = append(files, &File{item: item})
	}

	return files
}
