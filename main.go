package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var changeChan = make(chan fsnotify.Event)

type FileNotify struct {
	watch *fsnotify.Watcher
}

func NewFileNotify() *FileNotify {
	fn := new(FileNotify)
	fn.watch, _ = fsnotify.NewWatcher()
	return fn
}

func (this FileNotify) watchDir(dir string) {
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = this.watch.Add(path)
			if err != nil {
				return err
			}
			log.Println("开始监控目录", path)
		}
		return nil
	})
	go this.watchEvent()
}

func (this FileNotify) watchEvent() {
	for {
		select {
		case event := <-this.watch.Events:
			changeChan <- event
			fmt.Println(changeChan)
			fmt.Println("event.Op =>", event.Op.String())
			if event.Op&fsnotify.Create == fsnotify.Create {
				file, err := os.Stat(event.Name)
				if err == nil && file.IsDir() {
					this.watch.Add(event.Name)
					log.Println("增加监控文件", event.Name)
				}
				log.Println("创建文件 : ", event.Name)
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Println("写入文件 : ", event.Name)
			}

			if event.Op&fsnotify.Remove == fsnotify.Remove {
				file, err := os.Stat(event.Name)
				if err == nil && file.IsDir() {
					this.watch.Remove(event.Name)
					log.Println("删除监控文件", event.Name)
				}
				log.Println("删除文件 : ", event.Name)
			}

			if event.Op&fsnotify.Rename == fsnotify.Rename {
				this.watch.Remove(event.Name)
				log.Println("重命名文件 : ", event.Name)
			}

			if event.Op&fsnotify.Chmod == fsnotify.Chmod {
				log.Println("授权文件 : ", event.Name)
			}
		case err := <-this.watch.Errors:
			{
				log.Println("error ：", err)
				return
			}
		}
	}
}

func handlerEvent(client MyClient) {
	for {
		c := <-changeChan
		fmt.Println("从队列中取出的数据为：", c)
		if c.Op == fsnotify.Create || c.Op == fsnotify.Write {
			file, _ := os.Stat(c.Name)
			if file.IsDir() == false {
				client.UploadFile(c.Name)
			}
		}
	}
}

func main() {
	log.Println("...................欢迎使用berserk.....................")
	config := GetConfig()
	log.Println("load config：", config)
	client := MyClient{
		Config: config,
	}
	client.init()

	// 启动处理
	go handlerEvent(client)

	//启动监听
	NewFileNotify().watchDir(config.Watch.Path)

	http.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		response.Write([]byte("欢迎使用berserk"))
	})
	http.ListenAndServe(":8080", nil)
}
