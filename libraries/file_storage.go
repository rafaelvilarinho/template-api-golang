package libraries

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"api.template.com.br/helpers"
	"api.template.com.br/services/azure"
)

func ConvertBase64ToFile(fileString string) (*[]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(fileString)
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		return nil, err
	}

	return &decodedData, nil
}

func getFileExtensionFromContentType(contentType string) (string, error) {
	switch contentType {
	case "image/jpeg":
		return ".jpeg", nil
	case "image/jpg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	default:
		return "", fmt.Errorf("unsupported content type: %s", contentType)
	}
}

func GetFileTypeBase64(fileString string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(fileString)
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		return "", err
	}

	contentType := http.DetectContentType(decodedData)

	return contentType, nil
}

func getFileExtensionFromBase64(fileString string) (string, error) {
	decodedData, err := base64.StdEncoding.DecodeString(fileString)
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		return "", err
	}

	contentType := http.DetectContentType(decodedData)
	extension, err := getFileExtensionFromContentType(contentType)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return extension, nil
}

func ValidateFileExtensionBase64(fileString, ext string) (bool, error) {
	decodedData, err := base64.StdEncoding.DecodeString(fileString)
	if err != nil {
		fmt.Println("Error decoding Base64:", err)
		return false, err
	}

	contentType := http.DetectContentType(decodedData)
	extension, err := getFileExtensionFromContentType(contentType)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return fmt.Sprintf(".%s", ext) == extension, nil
}

func ValidateFileExtension(fileString, ext string) (bool, error) {
	extension, err := getFileExtensionFromBase64(fileString)
	if err != nil {
		return false, err
	}

	return fmt.Sprintf(".%s", ext) == extension, nil
}

func UploadFile(path, fileName, fileBase64Content, forceExtension string) (bool, error) {
	environment, _ := helpers.GetEnvironment()

	var extension string
	var err error
	if forceExtension != "" {
		extension = fmt.Sprintf(".%s", forceExtension)
	} else {
		extension, err = getFileExtensionFromBase64(fileBase64Content)
		if err != nil {
			return false, err
		}
	}

	return azure.UploadImageToStorageAccount(environment.AZURE_STORAGE_NAME, "path", path, fileName, fileBase64Content, extension)
}

func RemoveFile(path, fileName string) (bool, error) {
	environment, _ := helpers.GetEnvironment()

	return azure.RemoveImageFromStorageAccount(environment.AZURE_STORAGE_NAME, "path", path, fileName)
}
