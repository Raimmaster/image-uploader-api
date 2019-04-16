package ImageService

import (
	"bytes"
	"io/ioutil"
	"time"

	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
	"net/http"
)

type ImageResponse struct {
	Uploaded []string `json:"uploaded"`
}

type URLRequest struct {
	URLs []string `json:"urls"`
}

type UploadedStruct struct {
	Pending  []string `json:"pending"`
	Complete []string `json:"complete"`
	Failed   []string `json:"failed"`
}

type UploadStatusResponse struct {
	ID       string         `json:"id"`
	Created  string         `json:"created"`
	Finished string         `json:"finished"`
	Status   string         `json:"status"`
	Uploaded UploadedStruct `json:"uploaded"`
}

type UploadResponse struct {
	JobID string `json:"jobId"`
}

var statusResponseMap map[string]UploadStatusResponse

func StartServer() {
	router := mux.NewRouter()
	statusResponseMap = make(map[string]UploadStatusResponse)
	router.HandleFunc("/v1/images/upload", postUploadImages).Methods("POST")
	router.HandleFunc("/v1/images/upload/{jobId}", getImagesUploadStatus).Methods("GET", "HEAD")
	router.HandleFunc("/v1/images", getImages).Methods("GET", "HEAD")

	config.access_token = os.Getenv("ACCESS_TOKEN")

	fmt.Println("Starting server")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func postUploadImages(responseWriter http.ResponseWriter, request *http.Request) {
	var urls URLRequest
	err := json.NewDecoder(request.Body).Decode(&urls)

	if err != nil {
		fmt.Println(err)
		return
	}

	jobID := uuid.Must(uuid.NewV4())

	response := UploadResponse{
		JobID: jobID.String(),
	}

	go asyncUploadImages(response.JobID, urls)
	fmt.Println("/v1/images/upload-", jobID.String())
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(&response)
}

func imageToBase64(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(response.Body)
	newString := buffer.String()

	base64EncodedString := base64.StdEncoding.EncodeToString([]byte(newString))

	return base64EncodedString
}

// removing duplicates function adapted from https://kylewbanks.com/blog/creating-unique-slices-in-go
func removing_duplicates(input []string) []string {
	new_slice := make([]string, 0, len(input))
	original := make(map[string]bool)

	for _, val := range input {
		if _, ok := original[val]; !ok {
			original[val] = true
			new_slice = append(new_slice, val)
		}
	}

	return new_slice
}

func asyncUploadImages(jobID string, urls URLRequest) {
	urls.URLs = removing_duplicates(urls.URLs)
	uploadStatusResponse := UploadStatusResponse{
		ID:       jobID,
		Created:  string(time.Now().Format(time.RFC3339)),
		Finished: "",
		Status:   "pending",
		Uploaded: UploadedStruct{
			Pending:  urls.URLs,
			Complete: []string{},
			Failed:   []string{},
		},
	}

	statusResponseMap[jobID] = uploadStatusResponse

	for _, url := range urls.URLs {
		uploadStatusResponse.Status = "in-progress"
		base64EncodedString := imageToBase64(url)
		response := ImageUpload(base64EncodedString)
		if response.StatusCode == 200 {
			defer response.Body.Close()
			body, _ := ioutil.ReadAll(response.Body)
			imageURL := gjson.Get(string(body), "data.link")
			uploadStatusResponse.Uploaded.Complete = append(uploadStatusResponse.Uploaded.Complete, imageURL.String())
		} else {
			fmt.Println("Response: ", response.StatusCode)
			uploadStatusResponse.Uploaded.Failed = append(uploadStatusResponse.Uploaded.Failed, url)
		}
		if len(uploadStatusResponse.Uploaded.Pending) > 0 {
			uploadStatusResponse.Uploaded.Pending = uploadStatusResponse.Uploaded.Pending[:len(uploadStatusResponse.Uploaded.Pending)-1]
		}
		statusResponseMap[jobID] = uploadStatusResponse
	}
	uploadStatusResponse.Finished = string(time.Now().Format(time.RFC3339))
	uploadStatusResponse.Status = "complete"
	statusResponseMap[jobID] = uploadStatusResponse
}

func getImagesUploadStatus(responseWriter http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	jobID := params["jobId"]

	response := statusResponseMap[jobID]

	fmt.Printf("/v1/images/%s\n", response.ID)
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(&response)
}

func getImages(responseWriter http.ResponseWriter, request *http.Request) {
	accountImagesResponse := GetAccountImages()
	defer accountImagesResponse.Body.Close()
	body, _ := ioutil.ReadAll(accountImagesResponse.Body)
	imageURLs := gjson.Get(string(body), "data.#.link")
	response := &ImageResponse{
		Uploaded: []string{},
	}
	for _, imageLink := range imageURLs.Array() {
		response.Uploaded = append(response.Uploaded, imageLink.String())
	}

	fmt.Println("/v1/images")
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(response)
}
