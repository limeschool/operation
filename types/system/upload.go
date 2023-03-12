package system

type UploadRequest struct {
	Dir string `form:"name"  binding:"required"`
}
