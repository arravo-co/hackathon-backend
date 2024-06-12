package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/h2non/filetype"
)

func init() {

}

type UploadOpts struct {
	Folder string
}

func GetUploadedPic(file *multipart.FileHeader, opts ...UploadOpts) (string, error) {
	op := UploadOpts{
		Folder: "../uploads",
	}
	for _, opt := range opts {
		if opt.Folder != "" {
			op.Folder = opt.Folder
		}
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	fileName := file.Filename
	fmt.Println(fileName)
	var byt []byte
	_, err = src.Read(byt)
	if err != nil {
		return "", err
	}
	fmt.Printf("\nprofile picture file  %d     %s\n", file.Size, string(byt))
	is_img := filetype.IsImage(byt)
	if !is_img {
		//return "", fmt.Errorf("error: not an image type ")
	}
	path := filepath.Join(op.Folder, fileName)
	osFile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer osFile.Close()
	wB, err := io.Copy(osFile, src)
	if err != nil {
		return "", err
	}
	fmt.Printf("\n%d bytes written..\n", wB)
	return path, nil
}

func UploadPicToCloudinary(imgPath string) (*uploader.UploadResult, error) {
	fmt.Println("Preparing to upload... reading file...")
	byt, err := os.ReadFile(imgPath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	buf := bytes.NewBuffer(byt)
	cld, err := cloudinary.New()
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	fmt.Println("Preparing to upload... pushing to cloudinary...")
	uploadResult, err := cld.Upload.Upload(context.Background(), buf, uploader.UploadParams{})
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	err = CleanupUpload(imgPath)
	if err != nil {
		fmt.Printf("Failed to clean up: %s\n", err.Error())
	}
	return uploadResult, nil
}

func CleanupUpload(imgPath string) error {
	fmt.Println("file to delete:                            ", imgPath)
	return os.Remove(imgPath)
}
