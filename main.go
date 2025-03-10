package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"
)

func main() {
	// 读取配置文件
	config, err := readConfig("config.yml")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	// 发送 HTTP GET 请求获取发布信息
	req, err := http.NewRequest("GET", config.URL, nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	req.Header.Set("Authorization", "token "+config.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("Error:", err)
		}
	}(resp.Body)

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 解析 JSON 响应
	var releases []Release
	err = json.Unmarshal(body, &releases)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 按版本号排序
	sort.Slice(releases, func(i, j int) bool {
		return versionCompare(releases[i].TagName, releases[j].TagName) > 0
	})

	// 获取最新版本
	if len(releases) > 0 {
		latestRelease := releases[0]

		// 检查当前版本文件是否存在，如果不存在则创建一个
		_, err = os.Stat("current_version.txt")
		if os.IsNotExist(err) {
			_, err = os.Create("current_version.txt")
			if err != nil {
				fmt.Println("Error creating current version file:", err)
				return
			}
		}

		// 读取当前版本号
		currentVersionBytes, err := os.ReadFile("current_version.txt")
		if err != nil {
			fmt.Println("Error reading current version:", err)
			return
		}

		// 剔除额外的换行和空格
		currentVersion := strings.TrimSpace(string(currentVersionBytes))

		// 如果版本号一致，那么就不进行下载操作
		if currentVersion == latestRelease.TagName {
			fmt.Println("Current version is up to date.")
			return
		}

		// 下载文件
		fmt.Printf("Downloading files for latest release %s\n", latestRelease.TagName)
		for _, asset := range latestRelease.Assets {
			// 如果下载文件列表为空，或者文件在下载文件列表中，那么就下载文件
			if len(config.DownloadFiles) == 0 {
				err := downloadFile(asset.BrowserDownloadURL, config.DownloadDir, config.Token)
				if err != nil {
					fmt.Println("Error downloading file:", err)
					return
				}
			} else {
				for _, fileToDownloads := range config.DownloadFiles {
					if asset.Name == fileToDownloads {
						err := downloadFile(asset.BrowserDownloadURL, config.DownloadDir, config.Token)
						if err != nil {
							fmt.Println("Error downloading file:", err)
							return
						}
						break
					}
				}
			}
		}

		fmt.Println("Download completed.")

		// 验证下载的文件
		err = validateFiles(config.DownloadDir)
		if err != nil {
			fmt.Println("Error validating files:", err)
			return
		}

		// 检查最新版本文件是否存在，如果不存在则创建一个
		_, err = os.Stat("latest_version.txt")
		if os.IsNotExist(err) {
			_, err = os.Create("latest_version.txt")
			if err != nil {
				fmt.Println("Error creating latest version file:", err)
				return
			}
		}

		// 将最新的版本号写入 latest_version.txt 文件
		err = os.WriteFile("latest_version.txt", []byte(latestRelease.TagName), 0644)
		if err != nil {
			fmt.Println("Error writing latest version to file:", err)
			return
		}
	}
}
