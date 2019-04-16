package ImageService

import (
	"bytes"
	"fmt"

	"mime/multipart"
	"net/http"
)

type Configuration struct {
	access_token string
}

var config Configuration

func ImageUpload(base64imagehash string) *http.Response {
	url := "https://api.imgur.com/3/image"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("image", base64imagehash)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{
		CheckRedirect: func(request *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	request, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}

	access_token := config.access_token
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", access_token))

	request.Header.Set("Content-Type", writer.FormDataContentType())
	response, err := client.Do(request)

	if err != nil {
		fmt.Println(err)
	}

	return response
}

func GetAccountImages() *http.Response {
	url := "https://api.imgur.com/3/account/me/images"
	method := "GET"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	err := writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}

	access_token := config.access_token
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", access_token))

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)

	return res
}
