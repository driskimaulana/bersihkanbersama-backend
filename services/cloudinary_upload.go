package services

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"time"
)

var cld, _ = cloudinary.NewFromParams(os.Getenv("CLOUDINARY_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))

func UploadImage(c *gin.Context) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	file, err := c.FormFile("coverImage")
	defer cancel()
	if err != nil {
		return "", err
	}

	fileHandle, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func(fileHandle multipart.File) {
		err := fileHandle.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(fileHandle)

	result, err := cld.Upload.Upload(ctx, fileHandle, uploader.UploadParams{
		Folder: "bersihkanbersama",
	})

	if err != nil {
		return "", err
	}

	fmt.Println(result.URL)

	return result.URL, nil
}
