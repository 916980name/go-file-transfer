package third

import (
	"github.com/cloudinary/cloudinary-go/v2"
)

var (
	CLOUDINARY_FOLDER          string
	CLOUDINARY_MONGODB         string
	CLOUDINARY_MONGOCOLLECTION string
)

func Cloudinary_credentials() (*cloudinary.Cloudinary, error) {
	// Add your Cloudinary credentials, set configuration parameter
	// Secure=true to return "https" URLs, and create a context
	//===================
	cld, err := cloudinary.New()
	cld.Config.URL.Secure = true
	return cld, err
}
