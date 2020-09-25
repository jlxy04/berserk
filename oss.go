package main

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"os"
	"path/filepath"
)

type MyClient struct {
	Config Config
	Client *oss.Client
	Bucket *oss.Bucket
}

func (this *MyClient) init() error {
	cli, err := oss.New(this.Config.Oss.Endpoint, this.Config.Oss.AccessKeyId, this.Config.Oss.AccessKeySecret)
	if err != nil {
		log.Println(err)
		return err
	}
	this.Client = cli
	buc, err := cli.Bucket(this.Config.Oss.BucketName)
	if err != nil {
		log.Println(err)
		return err
	}
	this.Bucket = buc
	return nil
}

func (this MyClient) UploadFile(filePath string) error {
	fd, err := os.Open(filePath)
	if err != nil {
		log.Println(err)
		return err
	}
	path := this.Config.Oss.UploadPath + "/" + filepath.Base(fd.Name())
	log.Println("要上传的文件路径为：", path)
	err = this.Bucket.PutObject(path, fd)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
