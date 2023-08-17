package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
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
func Upload2Cos(buf *bytes.Buffer, filename string) string {
	//新建client
	client := NewClient()
	//上传文件
	_, err := client.Object.Put(context.Background(), filename, buf, nil)
	if err != nil {
		log.Fatal("文件上传失败：", err)
		return "上传失败"
	}
	return global.Config.Cos.CosUrl + "/" + filename
}
