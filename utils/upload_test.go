package utils

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestUpload(t *testing.T) {
	e := echo.New()

	path := "../../pic.jpg"
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	var byts []byte
	buf := bytes.NewBuffer(byts)
	mulW := multipart.NewWriter(buf)
	//encoding
	u, err := mulW.CreateFormFile("file", path)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	_, err = io.Copy(u, file)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	err = mulW.Close()
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/participants", buf)
	req.Header.Set("Content-Type", mulW.FormDataContentType())
	os.Setenv("CLOUDINARY_URL", "cloudinary://412863797731755:Ek8qQIKfbRl-NFiet2MXDeg9HTU@deklcb9ul")
	rec := httptest.NewRecorder()
	eCtx := e.NewContext(req, rec)
	fl, err := eCtx.FormFile("file")
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	result, err := GetUploadedPic(fl)
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	fmt.Printf("%v\n", result)
}

func TestCloudinaryUpload(t *testing.T) {
	//os.set
	os.Setenv("CLOUDINARY_URL", "cloudinary://412863797731755:Ek8qQIKfbRl-NFiet2MXDeg9HTU@deklcb9ul")
	result, err := UploadPicToCloudinary("C:\\Users\\TemitopeAlabi\\Downloads\\projects\\Arravo\\hackathon-backend\\uploads\\file.txt")
	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
	fmt.Printf("%v\n", result)
}
