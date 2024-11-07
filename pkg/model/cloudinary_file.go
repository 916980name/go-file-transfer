package model

import (
	"time"

	cl "github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryFile struct {
	AssetID      string    `json:"asset_id"`
	PublicID     string    `json:"public_id"`
	AssetFolder  string    `json:"asset_folder"`
	DisplayName  string    `json:"display_name"`
	Format       string    `json:"format"`
	Width        int       `json:"width,omitempty"`
	Height       int       `json:"height,omitempty"`
	ResourceType string    `json:"resource_type"`
	CreatedAt    time.Time `json:"created_at"`
	Bytes        int       `json:"bytes"`
	Type         string    `json:"type"`
	URL          string    `json:"url"`
	SecureURL    string    `json:"secure_url"`
	AccessMode   string    `json:"access_mode"`
	Title        string    `json:"title"`
	Desc         string    `json:"desc"`
}

func ClCopy(result *cl.UploadResult) *CloudinaryFile {
	if result == nil {
		return nil
	}
	return &CloudinaryFile{
		AssetID:      result.AssetID,
		PublicID:     result.PublicID,
		AssetFolder:  result.AssetFolder,
		DisplayName:  result.DisplayName,
		Format:       result.Format,
		Width:        result.Width,
		Height:       result.Height,
		ResourceType: result.ResourceType,
		CreatedAt:    result.CreatedAt,
		Bytes:        result.Bytes,
		Type:         result.Type,
		URL:          result.URL,
		SecureURL:    result.SecureURL,
		AccessMode:   result.AccessMode,
	}
}
