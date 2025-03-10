package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ProgressWriter struct {
	TotalSize    int64
	BytesWritten int64
	FileName     string
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.BytesWritten += int64(n)
	fmt.Printf("\rDownloading %s: %d of %d bytes (%.2f%%)", pw.FileName, pw.BytesWritten, pw.TotalSize, float64(pw.BytesWritten)*100/float64(pw.TotalSize))
	return n, nil
}

func downloadFile(url, downloadDir, token string) error {
	// 解析文件名
	_, file := filepath.Split(url)

	// 拼接下载路径
	downloadPath := filepath.Join(downloadDir, file)

	// 创建目录
	err := os.MkdirAll(downloadDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating download directory: %v", err)
	}

	// 创建文件
	out, err := os.Create(downloadPath)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer func(out *os.File) {
		err := out.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(out)

	// 发送 HTTP GET 请求下载文件
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	req.Header.Set("Authorization", "token "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(resp.Body)

	// 创建 ProgressWriter
	pw := &ProgressWriter{
		TotalSize: resp.ContentLength,
		FileName:  file,
	}

	// 将响应体写入文件，同时更新进度
	_, err = io.Copy(out, io.TeeReader(resp.Body, pw))
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	fmt.Println("\nDownload completed.")
	return nil
}
