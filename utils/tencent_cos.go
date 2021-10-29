package utils

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"v6tool/global"

	"github.com/tencentyun/cos-go-sdk-v5"
	"go.uber.org/zap"
)

type TencentCOS struct{}

// UploadFile upload file to COS
func (*TencentCOS) Upload(file string) (string, string, error) {
	fileobj, openError := os.Open(file)
	if openError != nil {
		global.V6TOOL_LOG.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	//windows情况下处理
	filename := strings.Replace(fileobj.Name(), "\\", "/", -1)
	//获取文件名
	fileKey := strings.Replace(filename, global.V6TOOL_CONFIG.Watcher.Dir, "", -1)

	fmt.Println(fileKey)

	client := NewClient()

	_, err := client.Object.PutFromFile(context.Background(), fileKey, file, nil)
	if err != nil {
		global.V6TOOL_LOG.Error("function file.Upload() Filed", zap.Any("err", err.Error()))
		return "", "", errors.New("function file.Upload() Filed, err:" + err.Error())
	}

	return global.V6TOOL_CONFIG.TencentCOS.BaseURL + "" + global.V6TOOL_CONFIG.TencentCOS.PathPrefix + "/" + fileKey, fileKey, nil
}

// UploadFile upload file to COS
func (*TencentCOS) UploadFile(file *multipart.FileHeader) (string, string, error) {
	client := NewClient()
	f, openError := file.Open()
	if openError != nil {
		global.V6TOOL_LOG.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	fileKey := fmt.Sprintf("%d%s", time.Now().Unix(), file.Filename)

	_, err := client.Object.Put(context.Background(), global.V6TOOL_CONFIG.TencentCOS.PathPrefix+"/"+fileKey, f, nil)
	if err != nil {
		panic(err)
	}
	return global.V6TOOL_CONFIG.TencentCOS.BaseURL + "/" + global.V6TOOL_CONFIG.TencentCOS.PathPrefix + "/" + fileKey, fileKey, nil
}

// DeleteFile delete file form COS
func (*TencentCOS) DeleteFile(key string) error {
	client := NewClient()
	name := global.V6TOOL_CONFIG.TencentCOS.PathPrefix + "/" + key
	_, err := client.Object.Delete(context.Background(), name)
	if err != nil {
		global.V6TOOL_LOG.Error("function bucketManager.Delete() Filed", zap.Any("err", err.Error()))
		return errors.New("function bucketManager.Delete() Filed, err:" + err.Error())
	}
	return nil
}

// NewClient init COS client
func NewClient() *cos.Client {
	urlStr, _ := url.Parse("https://" + global.V6TOOL_CONFIG.TencentCOS.Bucket + ".cos." + global.V6TOOL_CONFIG.TencentCOS.Region + ".myqcloud.com")
	baseURL := &cos.BaseURL{BucketURL: urlStr}
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  global.V6TOOL_CONFIG.TencentCOS.SecretID,
			SecretKey: global.V6TOOL_CONFIG.TencentCOS.SecretKey,
		},
	})
	return client
}
