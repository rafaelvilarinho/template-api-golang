package azure

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"api.template.com.br/helpers"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

func UploadImageToStorageAccount(accountName, containerName, path, fileName, content, extension string) (bool, error) {
	environment, _ := helpers.GetEnvironment()
	fileKey := fmt.Sprintf("%s/%s%s", path, fileName, extension)

	cred, err := azblob.NewSharedKeyCredential(accountName, environment.AZURE_STORAGE_ACCESS_KEY)
	if err != nil {
		return false, err
	}

	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	if err != nil {
		return false, err
	}

	// containerCreateResp, err := client.CreateContainer(context.TODO(), containerName, nil)
	// if err != nil {
	// 	return false, err
	// }

	fileBase64, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return false, err
	}

	_, err = client.UploadStream(context.TODO(),
		containerName,
		fileKey,
		strings.NewReader(string(fileBase64)),
		&azblob.UploadStreamOptions{},
	)
	if err != nil {
		return false, err
	}

	return err == nil, nil
}

func RemoveImageFromStorageAccount(accountName, containerName, path, fileName string) (bool, error) {
	environment, _ := helpers.GetEnvironment()

	cred, err := azblob.NewSharedKeyCredential(accountName, environment.AZURE_STORAGE_ACCESS_KEY)
	if err != nil {
		return false, err
	}

	client, err := azblob.NewClientWithSharedKeyCredential(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), cred, nil)
	if err != nil {
		return false, err
	}

	_, err = client.DeleteBlob(context.TODO(), containerName, fmt.Sprintf("%s/%s", path, fileName), nil)
	if err != nil {
		return false, err
	}

	return true, nil
}
