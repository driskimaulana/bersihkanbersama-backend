package services

import (
	"context"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"time"
)

var cld, _ = cloudinary.NewFromParams("dscbb3cu2", "236375965318535", "rMTpUebQW7YudNxweG2na60Tqfs")

func UploadImage(c *gin.Context) (string, error) {
	fmt.Println("1")
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
	fmt.Println("3" + file.Filename)

	result, err := cld.Upload.Upload(ctx, fileHandle, uploader.UploadParams{
		Folder: "bersihkanbersama",
	})

	if err != nil {
		return "", err
	}
	fmt.Println("2")

	fmt.Println(result)

	return result.URL, nil
}
