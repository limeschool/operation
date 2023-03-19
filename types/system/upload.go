package system

type UploadRequest struct {
	Dir string `json:"dir" form:"dir"  binding:"required"`
}
