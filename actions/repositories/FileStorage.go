package repositories

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func ConnectFileStorage() *session.Session {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String(getEnv("STORAGE_REGION", "")),
		Credentials: credentials.NewStaticCredentials(
			getEnv("STORAGE_ID", ""),     // id
			getEnv("STORAGE_SECRET", ""), // secret
			"",
		),
	})

	if err != nil {
		log.Fatal(err)
	}
	return session
}

func Base64toByte(Base64, typeBase64 string) []byte {
	file := ""
	if typeBase64 == "image" {
		checkbase64 := string([]rune(Base64)[16:22])
		if checkbase64 == "base64" {
			file = string([]rune(Base64)[23:])
		} else {
			file = string([]rune(Base64)[22:])
		}
	}
	if typeBase64 == "json" {
		log.Print("base64 :", Base64)
		file = string([]rune(Base64)[29:])
		log.Print("json : ", file)
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
		Bucket: aws.String(getEnv("STORAGE_NAME", "")),
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

func DeleteFile(path, fileName string) error {
	sess := ConnectFileStorage()
	svc := s3.New(sess)
	log.Print("path ", path+"/"+fileName)
	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(getEnv("STORAGE_NAME", "")),
		Delete: &s3.Delete{
			Objects: []*s3.ObjectIdentifier{
				{
					Key: aws.String(path + "/" + fileName),
				},
			},
			Quiet: aws.Bool(false),
		},
	}
	_, err := svc.DeleteObjects(input)
	if err != nil {
		return err
	}
	return nil
}

func GetFile(path, fileName string) (string, error) {
	sess := ConnectFileStorage()
	downloader := s3manager.NewDownloader(sess)
	file, err := os.Create(fileName)
	if err != nil {
		log.Print("create file error")
		return "", err
	}
	defer file.Close()
	pathfile := path + "/" + fileName
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(getEnv("STORAGE_NAME", "")),
			Key:    aws.String(pathfile),
		})
	if err != nil {
		log.Print("dowload file error")
		return "", err
	}

	fmt.Println("Downloaded", file.Name(), numBytes, "bytes")

	return fileName, nil
}

func DownloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
	}
	defer response.Body.Close()

	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}
