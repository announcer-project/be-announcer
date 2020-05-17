package main

import (
	"be_nms/database"
	"be_nms/routes"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo/v4/middleware"
)

// func TestBlob() {
// 	// From the Azure portal, get your storage account name and key and set environment variables.
// 	accountName, accountKey := "sqlvafao4cvoaektc2", "SW/QzVdcTUK5xusO3X6Wjmg45XwBj7POHnfRy8lMpUtkILuS/aoFd6fQTSOq63+c1fTnBNPjnNsSQMXXtgz17A=="
// 	if len(accountName) == 0 || len(accountKey) == 0 {
// 		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
// 	}

// 	// Create a default request pipeline using your storage account name and account key.
// 	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
// 	if err != nil {
// 		log.Fatal("Invalid credentials with error: " + err.Error())
// 	}
// 	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

// 	// Create a random string for the quick start container
// 	containerName := "images"

// 	// From the Azure portal, get your storage account blob service URL endpoint.
// 	URL, _ := url.Parse(
// 		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

// 	// Create a ContainerURL object that wraps the container URL and a request
// 	// pipeline to make requests.
// 	containerURL := azblob.NewContainerURL(*URL, p)

// 	// Create the container
// 	// fmt.Printf("Creating a container named %s\n", containerName)
// 	ctx := context.Background() // This example uses a never-expiring context
// 	// _, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
// 	// handleErrors(err)
// 	// Create a file to test the upload and download.
// 	fmt.Printf("Creating a dummy file to test the upload and download\n")
// 	data := []byte("hello world this is a blob\n")
// 	fileName := "test-1"
// 	err = ioutil.WriteFile(fileName, data, 0700)

// 	// Here's how to upload a blob.
// 	blobURL := containerURL.NewBlockBlobURL(fileName)
// 	file, err := os.Open(fileName)

// 	// You can use the low-level Upload (PutBlob) API to upload files. Low-level APIs are simple wrappers for the Azure Storage REST APIs.
// 	// Note that Upload can upload up to 256MB data in one shot. Details: https://docs.microsoft.com/rest/api/storageservices/put-blob
// 	// To upload more than 256MB, use StageBlock (PutBlock) and CommitBlockList (PutBlockList) functions.
// 	// Following is commented out intentionally because we will instead use UploadFileToBlockBlob API to upload the blob
// 	// _, err = blobURL.Upload(ctx, file, azblob.BlobHTTPHeaders{ContentType: "text/plain"}, azblob.Metadata{}, azblob.BlobAccessConditions{})
// 	// handleErrors(err)

// 	// The high-level API UploadFileToBlockBlob function uploads blocks in parallel for optimal performance, and can handle large files as well.
// 	// This function calls StageBlock/CommitBlockList for files larger 256 MBs, and calls Upload for any file smaller
// 	fmt.Printf("Uploading the file with blob name: %s\n", fileName)
// 	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
// 		BlockSize:   4 * 1024 * 1024,
// 		Parallelism: 16})

// 	// List the container that we have created above
// 	fmt.Println("Listing the blobs in the container:")
// 	for marker := (azblob.Marker{}); marker.NotDone(); {
// 		// Get a result segment starting with the blob indicated by the current Marker.
// 		listBlob, _ := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})

// 		// ListBlobs returns the start of the next segment; you MUST use this to get
// 		// the next segment (after processing the current result segment).
// 		marker = listBlob.NextMarker

// 		// Process the blobs returned in this result segment (if the segment is empty, the loop body won't execute)
// 		for _, blobInfo := range listBlob.Segment.BlobItems {
// 			fmt.Print("	Blob name: " + blobInfo.Name + "\n")
// 		}
// 	}

// 	// Here's how to download the blob
// 	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)

// 	// NOTE: automatically retries are performed if the connection fails
// 	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

// 	// read the body into a buffer
// 	downloadedData := bytes.Buffer{}
// 	_, err = downloadedData.ReadFrom(bodyStream)
// 	log.Print(downloadResponse)
// 	blobURLd := containerURL.NewBlobURL("rich-menu.png")
// 	downloadResponse, err = blobURLd.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)

// 	// NOTE: automatically retries are performed if the connection fails
// 	bodyStream = downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})

// 	// read the body into a buffer
// 	downloadedData = bytes.Buffer{}
// 	_, err = downloadedData.ReadFrom(bodyStream)

// 	// The downloaded blob data is in downloadData's buffer. :Let's print it
// 	fmt.Printf("Downloaded the blob: " + downloadedData.String())
// 	fmt.Printf("Downloaded the blob: ", downloadedData.Bytes())
// 	// f, err := os.Create("")
// 	err = ioutil.WriteFile("rich-menu.png", downloadedData.Bytes(), 0700)

// }

func main() {
	e := routes.Init()
	db := database.Open()
	defer db.Close()
	// database.Migration(db)
	// database.SetData(db)
	e.Use(middleware.CORS())
	e.Logger.Fatal(e.Start(":8080"))
}
