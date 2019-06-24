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

var POST_URL = "https://api.imgur.com/3/image"
var GET_URL = "https://api.imgur.com/3/account/me/images"

func ImageAPICall(url string, base64imagehash string) *http.Response {
	method := "GET"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	if url == POST_URL {
		method = "POST"
		err := writer.WriteField("image", base64imagehash)

		if err != nil {
			fmt.Println(err)
		}
	}

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