package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

func main() {
	// 定义命令行参数
	fileDir := flag.String("fileDir", "/home/workspace/golang/data-transmission/wget_server", "Directory to serve files from")
	port := flag.String("port", "4567", "Port to listen on")

	// 解析命令行参数
	flag.Parse()

	// 处理文件下载请求
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		fileName := strings.TrimPrefix(r.URL.Path, "/download/")
		filePath := filepath.Join(*fileDir, fileName)

		// 设置响应头
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", "application/octet-stream")

		// 发送文件
		http.ServeFile(w, r, filePath)
	})

	// 启动服务器
	log.Printf("服务器启动，监听端口: %s", *port)
	if err := http.ListenAndServe("0.0.0.0:"+*port, nil); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

// ./test -fileDir=/home/workspace/project/test/wget -port=4567
// wget http://192.168.8.244:4567/download/baihetan_fengdong.sql -O source/baihetan_fengdong.sql
// /opt/backup/suanfa
// /home/workspace/project/test/wget
