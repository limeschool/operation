package system

type LoginLogRequest struct {
	Page   int    `json:"page" form:"page" binding:"required" sql:"-"`
	Count  int    `json:"count" form:"count"  binding:"required,max=50"  sql:"-"`
	Phone  string `json:"phone" form:"phone"`
	Status *bool  `json:"status" form:"status"`
	Start  int64  `json:"start" form:"start" sql:"> ?" field:"created_at"`
	End    int64  `json:"end" form:"end" sql:"< ?" field:"created_at"`
}
