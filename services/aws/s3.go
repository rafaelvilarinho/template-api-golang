package aws

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getSession() (*session.Session, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("sa-east-1")},
	)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func UploadImageToS3(bucket, path, fileName, fileBase64Content, extension string) (bool, error) {
	// Create a new session with default session credentials
	sess, err := getSession()
	if err != nil {
		return false, err
	}

	// Create a new S3 service client
	svc := s3.New(sess)

	fileKey := fmt.Sprintf("%s/%s%s", path, fileName, extension)

	decodedData, err := base64.StdEncoding.DecodeString(fileBase64Content)
	if err != nil {
		return false, err
	}

	// Configure the S3 object input parameters
	input := &s3.PutObjectInput{
		Body:          bytes.NewReader(decodedData),
		Bucket:        aws.String(bucket),
		Key:           aws.String(fileKey),
		ContentType:   aws.String(http.DetectContentType(decodedData)),
		ContentLength: aws.Int64(int64(len(decodedData))),
	}

	// Upload the image to S3
	_, err = svc.PutObject(input)
	if err != nil {
		return false, err
	}

	return true, nil
}

func RemoveImageFromS3(bucket, path, fileName string) (bool, error) {
	// Create a new session with default session credentials
	sess, err := getSession()
	if err != nil {
		return false, err
	}

	// Create a new S3 service client
	svc := s3.New(sess)

	fileKey := fmt.Sprintf("%s/%s", path, fileName)

	// Configure the S3 object input parameters
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	}

	// Upload the image to S3
	_, err = svc.DeleteObject(input)
	if err != nil {
		return false, err
	}

	return true, nil
}

// func UploadImagesToS3(bucketName string, imagePaths []string) bool {
// 	// Create a new session with default session credentials
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String("us-east-1")},
// 	)
// 	if err != nil {
// 		return false
// 	}

// 	// Create a channel to communicate image upload status
// 	statusChan := make(chan bool)

// 	// Launch a new Goroutine for each image upload
// 	for _, imagePath := range imagePaths {
// 		go func(imagePath string) {
// 			// Upload the image to S3
// 			ok, _ := UploadImageToS3(bucketName, imagePath, sess)
// 			// Send the upload status to the status channel
// 			statusChan <- ok
// 		}(imagePath)
// 	}

// 	// Wait for all uploads to complete and collect status
// 	for range imagePaths {
// 		ok := <-statusChan
// 		if !ok {
// 			return false
// 		}
// 	}

// 	return true
// }
