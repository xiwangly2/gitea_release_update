package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func validateFiles(downloadDir string) error {
	files, err := os.ReadDir(downloadDir)
	if err != nil {
		return fmt.Errorf("error reading download directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() || strings.HasSuffix(file.Name(), ".sha256") {
			continue
		}

		filePath := filepath.Join(downloadDir, file.Name())
		sha256Path := filePath + ".sha256"

		// 检查是否有对应的 sha256 文件
		if _, err := os.Stat(sha256Path); os.IsNotExist(err) {
			return fmt.Errorf("SHA256 file not found for %s", file.Name())
		}

		// 校验文件的 SHA256 散列值
		sha256File, err := os.ReadFile(sha256Path)
		if err != nil {
			return fmt.Errorf("error reading SHA256 file: %v", err)
		}

		// 计算文件的 SHA256 散列值
		fileHash, err := calculateSHA256(filePath)
		if err != nil {
			return fmt.Errorf("error calculating SHA256 hash: %v", err)
		}

		// 将文件的 SHA256 散列值转换为字符串，并移除文件名
		expectedHash := strings.Fields(string(sha256File))[0]

		// 检查文件的 SHA256 散列值是否与预期值匹配
		if !strings.EqualFold(fileHash, expectedHash) {
			return fmt.Errorf("SHA256 hash mismatch for %s, expected: %s, actual: %s", file.Name(), expectedHash, fileHash)
		}
	}

	return nil
}
