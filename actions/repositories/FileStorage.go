package repositories

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"

	"github.com/Azure/azure-storage-blob-go/azblob"
)

func ConnectFileStorage() (azblob.ContainerURL, context.Context) {
	// From the Azure portal, get your storage account name and key and set environment variables.
	accountName, accountKey := "sqlvafao4cvoaektc2", "SW/QzVdcTUK5xusO3X6Wjmg45XwBj7POHnfRy8lMpUtkILuS/aoFd6fQTSOq63+c1fTnBNPjnNsSQMXXtgz17A=="
	if len(accountName) == 0 || len(accountKey) == 0 {
		log.Fatal("Either the AZURE_STORAGE_ACCOUNT or AZURE_STORAGE_ACCESS_KEY environment variable is not set")
	}

	// Create a default request pipeline using your storage account name and account key.
	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		log.Fatal("Invalid credentials with error: " + err.Error())
	}
	p := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	// Create a random string for the quick start container
	containerName := "images"

	// From the Azure portal, get your storage account blob service URL endpoint.
	URL, _ := url.Parse(
		fmt.Sprintf("https://%s.blob.core.windows.net/%s", accountName, containerName))

	// Create a ContainerURL object that wraps the container URL and a request
	// pipeline to make requests.
	containerURL := azblob.NewContainerURL(*URL, p)

	// Create the container
	// fmt.Printf("Creating a container named %s\n", containerName)
	ctx := context.Background() // This example uses a never-expiring context
	// _, err = containerURL.Create(ctx, azblob.Metadata{}, azblob.PublicAccessNone)
	// handleErrors(err)
	return containerURL, ctx
}

func CreateFile(data []byte, fileName string) error {
	containerURL, ctx := ConnectFileStorage()
	// Create a file to test the upload and download.
	err := ioutil.WriteFile(fileName, data, 0700)
	if err != nil {
		return err
	}
	// Here's how to upload a blob.
	blobURL := containerURL.NewBlockBlobURL(fileName)
	file, err := os.Open(fileName)
	if err != nil {
		return nil
	}
	fmt.Printf("Uploading the file with blob name: %s\n", fileName)
	_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
		BlockSize:   4 * 1024 * 1024,
		Parallelism: 16})
	if err != nil {
		return nil
	}
	return nil
}

func GetFile(fileName string) (string, error) {
	containerURL, ctx := ConnectFileStorage()

	blobURL := containerURL.NewBlobURL(fileName)
	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
	if err != nil {
		return "", err
	}
	// NOTE: automatically retries are performed if the connection fails
	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
	// read the body into a buffer
	downloadedData := bytes.Buffer{}
	_, err = downloadedData.ReadFrom(bodyStream)
	if err != nil {
		return "", err
	}
	// The downloaded blob data is in downloadData's buffer. :Let's print it
	// fmt.Printf("Downloaded the blob: " + downloadedData.String())
	// fmt.Printf("Downloaded the blob: ", downloadedData.Bytes())
	// f, err := os.Create("")
	err = ioutil.WriteFile(fileName, downloadedData.Bytes(), 0700)
	if err != nil {
		return "", err
	}
	return fileName, nil
}
