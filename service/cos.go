package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/Crazypointer/simple-tok/global"
	"github.com/tencentyun/cos-go-sdk-v5"
)

// NewClient 创建 Client
func NewClient() *cos.Client {
	u, _ := url.Parse(global.Config.Cos.CosUrl)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  global.Config.Cos.SecretID,
			SecretKey: global.Config.Cos.SecretKey,
		},
	})
	return client
}

// FindBucket 查询 Bucket
func FindBucket(c *cos.Client) (*cos.Bucket, error) {
	s, _, err := c.Service.Get(context.Background())
	if err != nil {
		return nil, err
	}
	for _, b := range s.Buckets {
		fmt.Printf("%#v\n", b)
	}
	return &s.Buckets[0], nil
}

// 上传文件到cos
func Upload2Cos(file *multipart.FileHeader, filename string) string {
	//新建client
	client := NewClient()
	//读取文件
	var buf bytes.Buffer
	f, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return "文件打开失败"
	}
	io.Copy(&buf, f)
	//上传文件
	_, err = client.Object.Put(context.Background(), filename, &buf, nil)
	if err != nil {
		fmt.Println(err)
		return "上传失败"
	}
	return global.Config.Cos.CosUrl + "/" + filename
}
