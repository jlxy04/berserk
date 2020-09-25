package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Oss   Oss   `yaml:"oss"`
	Watch Watch `yaml:"watch"`
}

type Oss struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	BucketName      string `yaml:"bucketName"`
	UploadPath      string `yaml:"uploadPath"`
}

type Watch struct {
	Path string `yaml:"path"`
}

func GetConfig() Config {
	dir, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	configPath := dir + string(os.PathSeparator) + "conf.yml"
	fmt.Println("配置文件路径：", configPath)
	b, err := ioutil.ReadFile(configPath)
	config := Config{}
	yaml.Unmarshal(b, &config)
	return config
}
