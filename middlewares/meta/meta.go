package meta

import (
	"encoding/json"
	"github.com/limeschool/gin"
	"operation/consts"
	"operation/errors"
)

type Metadata struct {
	UserID    int64  `json:"UserID"`
	RoleID    int64  `json:"RoleID"`
	RoleKey   string `json:"RoleKey"`
	Username  string `json:"Username"`
	DataScope string `json:"DataScope"`
	TeamID    int64  `json:"TeamID"`
}

func (m *Metadata) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func Get(ctx *gin.Context) (*Metadata, error) {
	val, is := ctx.Get(consts.JwtMapClaims)
	if !is {
		return nil, errors.TokenDataError
	}

	meta, is := val.(*Metadata)
	if !is {
		return nil, errors.TokenDataError
	}
	return meta, nil
}

func Parse(m any) (*Metadata, error) {
	meta := Metadata{}

	// 序列化
	b, err := json.Marshal(m)
	if err != nil {
		return nil, errors.TokenDataError
	}

	// 反序列化
	if json.Unmarshal(b, &meta) != nil {
		return nil, errors.TokenDataError
	}
	return &meta, nil
}
