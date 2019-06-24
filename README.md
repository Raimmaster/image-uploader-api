# Image Uploader API

This is a small API that uploads images to Imgur. The function of `ImageAPICall` is adapted by the Go functions provided by Imgur in their [API docs](https://apidocs.imgur.com/).

This service is dockerized, so in order to run it, you'll need an access token. In order to generate one, you can look into [Authorization and OAuth](https://apidocs.imgur.com/#authorization-and-oauth) in Imgur's API docs. Once you have the token, you can set it up as an environment variable and utilize it as an argument while running the container.

You can run the container in the following manner:
```Go
export ACCESS_TOKEN=token
docker build --build-arg ACCESS_TOKEN=$ACCESS_TOKEN -t image-service:0.0.1 .
docker run -p 8000:8000 image-service:0.0.1 
```

You should see a `Starting Server` output to know it has successfully started, and you should be able to make requests to `localhost:8000` if you're running locally.

## Endpoints

### POST: /v1/images/upload
**Attributes:**
- **urls:** An array of URLs to images that will be uploaded. 

**Example request body:**

```JSON
{
	"urls": [
		"https://images.shazam.com/coverart/t44239036-b40621303_s400.jpg",
		"https://ih0.redbubble.net/image.77106686.9243/mp,550x550,gloss,ffffff,t.3.sssjpg",
		"https://leadiq.com/img/logo.png"
		]
}
```

**Response**
On success, returns immediately with an appropriate status code with the id of the job.

**Response body**
Attributes:
- **jobId:** The id of the upload job that was just submitted.

**Example response body:**
```JSON
{ 
	"jobId": "55355b7c-9b86-4a1a-b32e-6cdd6db07183" 
}
```

### GET /v1/images/upload/:jobId

The request has no body and no query parameters. `:jobId` is an ID returned from the POST upload images API.

**Response**
On success, returns immediately with an appropriate status code with the id of the job.

**Response body**
Attributes:
- **id:** The id of the upload job.
- **created:** When job was created. In ISO8601 format (YYYY-MM-DDTHH:mm:ss.sssZ) for GMT.
- **finished:** When job was completed. In same format as created. Is null​ if status is not complete.
- **status:** The status of the entire upload job. Is one of:
	- **pending:** indicates job has not started processing.
	- **in-progress:** job has started processing.
	- **complete:** job is complete.
- **uploaded:** An object of arrays containing the set of URLs submitted, in several arrays indicating the status of that image URL upload (pending, complete, failed).

**Example Response Body:**
```JSON
{
    "id": "6c0f7894-4bc7-47b8-b680-cc3350595300",
    "created": "2019-04-16T18:19:50Z",
    "finished": "2019-04-16T18:19:53Z",
    "status": "complete",
    "uploaded": {
        "pending": [],
        "complete": [
            "https://i.imgur.com/tXWwhBP.jpg",
            "https://i.imgur.com/lXweEKv.png"
        ],
        "failed": [
            "https://ih0.redbubble.net/image.77106686.9243/mp,550x550,gloss,ffffff,t.3.sssjpg"
        ]
    }
}
```

### GET /v1/images
Gets the links of all images uploaded to Imgur. These links will be accessible by anyone.

**Response**
On success, return an array of the Imgur links to the successfully uploaded images. The links should be public.

**Response body**
Attributes:
- **uploaded:** An array of the Imgur links to the uploaded images.

**Example Response Body**
```JSON
{
    "uploaded": [
        "https://i.imgur.com/n98GxmC.jpg",
        "https://i.imgur.com/xhS3LIV.jpg",
        "https://i.imgur.com/YjvqG5i.png",
        "https://i.imgur.com/uAUqui2.jpg",
        "https://i.imgur.com/DcaPtrj.jpg",
        "https://i.imgur.com/MMwnWAY.jpg",
        "https://i.imgur.com/KWv1Olx.jpg",
        "https://i.imgur.com/9mfhUmD.png",
        "https://i.imgur.com/YVaC9C8.png",
        "https://i.imgur.com/Zqje7EA.jpg"
    ]
}
```

