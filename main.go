package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"v6tool/core"
	"v6tool/global"
	"v6tool/utils"

	"go.uber.org/zap"

	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli/v2"
)

var dirfiles []string

type Watch struct {
	watch *fsnotify.Watcher
}

//监控目录
func (w *Watch) watchDir(dir string) {
	//通过Walk来遍历目录下的所有子目录
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//这里判断是否为目录，只需监控目录即可
		if info.IsDir() {
			path, err := filepath.Abs(path)
			if err != nil {
				return err
			}
			err = w.watch.Add(path)
			if err != nil {
				return err
			}
			fmt.Println("监控目录 : ", path)
		}
		return nil
	})
	go func() {
		for {
			select {
			case ev := <-w.watch.Events:
				{
					if ev.Op&fsnotify.Create == fsnotify.Create {
						global.V6TOOL_LOG.Info("创建文件" + "file:" + ev.Name)
						//这里获取新创建文件的信息，如果是目录，则加入监控中
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							w.watch.Add(ev.Name)
							global.V6TOOL_LOG.Info("添加监控" + "dir:" + ev.Name)
						}
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						global.V6TOOL_LOG.Info("写入文件-同步cos" + "file:" + ev.Name)

						cos := &utils.TencentCOS{}
						file := ev.Name
						if utils.IsFile(file) {
							//上传
							url, filekey, err := cos.Upload(file)
							if err != nil {
								global.V6TOOL_LOG.Error("upload file: key "+"filekey:"+filekey, zap.Any("err", err))
							} else {
								global.V6TOOL_LOG.Info("upload file: key "+"filekey:"+filekey, zap.Any("url", url))
							}
						}
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						global.V6TOOL_LOG.Info("删除文件" + "file:" + ev.Name)
						//如果删除文件是目录，则移除监控
						fi, err := os.Stat(ev.Name)
						if err == nil && fi.IsDir() {
							w.watch.Remove(ev.Name)
							global.V6TOOL_LOG.Info("删除监控" + "dir:" + ev.Name)
						}
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						global.V6TOOL_LOG.Info("重命名文件" + "file:" + ev.Name)

						//如果重命名文件是目录，则移除监控
						w.watch.Remove(ev.Name)
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						global.V6TOOL_LOG.Info("修改权限" + "file:" + ev.Name)

					}
				}
			case err := <-w.watch.Errors:
				{
					fmt.Println("error : ", err)
					return
				}
			}
		}
	}()
}

func watchAssets() {
	watch, _ := fsnotify.NewWatcher()
	w := Watch{
		watch: watch,
	}
	defer watch.Close()
	done := make(chan bool)

	w.watchDir(global.V6TOOL_CONFIG.Watcher.Dir)

	<-done

}

func main() {
	global.V6TOOL_VP = core.Viper() // 初始化Viper
	global.V6TOOL_LOG = core.Zap()  // 初始化zap日志库
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:    "cos-init",
				Aliases: []string{"ci"},
				Usage:   "初次同步全部文件",
				Action: func(c *cli.Context) error {
					files := utils.GetAllFile(global.V6TOOL_CONFIG.Watcher.Dir)
					//去除全部文件夹
					cos := &utils.TencentCOS{}
					for key, file := range files {
						if utils.IsFile(file) {
							//上传
							url, filekey, err := cos.Upload(file)
							if err != nil {
								global.V6TOOL_LOG.Error("upload file: key "+strconv.Itoa(key)+"filekey:"+filekey, zap.Any("err", err))
							} else {
								global.V6TOOL_LOG.Info("upload file: key "+strconv.Itoa(key)+"filekey:"+filekey, zap.Any("url", url))
							}
						}
					}
					return nil
				},
			},
			{
				Name:    "cos-watch",
				Aliases: []string{"cw"},
				Usage:   "cos文件监听同步",
				Action: func(c *cli.Context) error {
					fmt.Println("Hello friend!")
					watchAssets()
					return nil
				},
			},
			{
				Name:    "test",
				Aliases: []string{"t"},
				Usage:   "测试",
				Action: func(c *cli.Context) error {
					fmt.Println("hello")
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
