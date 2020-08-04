package repositories

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func ConnectFileStorage() *session.Session {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-1"),
		Credentials: credentials.NewStaticCredentials(
			"AKIASCHC264LQFS6QQ2A",                     // id
			"q3DNPz3+BMf9qxy34gvhS5NgPgCvPzRjH+r12wN0", // secret
			"",
		),
	})

	if err != nil {
		log.Fatal(err)
	}
	return session
}

func Base64toByte(imageBase64 string) []byte {
	file := ""
	checkbase64 := string([]rune(imageBase64)[16:22])
	if checkbase64 == "base64" {
		file = string([]rune(imageBase64)[23:])
	} else {
		file = string([]rune(imageBase64)[22:])
	}
	dec, err := base64.StdEncoding.DecodeString(file)
	if err != nil {
		return nil
	}
	return dec
}

func CreateFile(sess *session.Session, imageByte []byte, fileName, pathStorage string) error {
	// Create a file to test the upload and download.
	err := ioutil.WriteFile(fileName, imageByte, 0700)
	if err != nil {
		return err
	}

	file, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open file %q, %v", fileName, err)
	}
	uploader := s3manager.NewUploader(sess)
	path := pathStorage + "/" + fileName
	log.Print("path", path)
	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("announcer-project"),
		Key:    aws.String(path),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file, %v", err)
	}
	fmt.Printf("file uploaded to, %s\n", result.Location)
	file.Close()
	log.Print("filename", fileName)
	if err := os.Remove(fileName); err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// func GetFile(fileName string) (string, error) {
// 	containerURL, ctx := ConnectFileStorage()

// 	blobURL := containerURL.NewBlobURL(fileName)
// 	downloadResponse, err := blobURL.Download(ctx, 0, azblob.CountToEnd, azblob.BlobAccessConditions{}, false)
// 	if err != nil {
// 		return "", err
// 	}
// 	// NOTE: automatically retries are performed if the connection fails
// 	bodyStream := downloadResponse.Body(azblob.RetryReaderOptions{MaxRetryRequests: 20})
// 	// read the body into a buffer
// 	downloadedData := bytes.Buffer{}
// 	_, err = downloadedData.ReadFrom(bodyStream)
// 	if err != nil {
// 		return "", err
// 	}
// 	// The downloaded blob data is in downloadData's buffer. :Let's print it
// 	// fmt.Printf("Downloaded the blob: " + downloadedData.String())
// 	// fmt.Printf("Downloaded the blob: ", downloadedData.Bytes())
// 	// f, err := os.Create("")
// 	err = ioutil.WriteFile(fileName, downloadedData.Bytes(), 0700)
// 	if err != nil {
// 		return "", err
// 	}
// 	return fileName, nil
// }
