package requests

import (
	"mime/multipart"

	"github.com/goravel/framework/contracts/http"
)

type StoreAttachmentRequest struct {
	File *multipart.FileHeader `form:"file" json:"file"`
}

func (r *StoreAttachmentRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *StoreAttachmentRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{
		"file": "required|file|max:5120|mimes:png,jpg,pdf,txt",
	}
}
