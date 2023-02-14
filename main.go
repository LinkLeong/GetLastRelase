/*
 * @Author: a624669980@163.com a624669980@163.com
 * @Date: 2023-02-14 10:38:57
 * @LastEditors: a624669980@163.com a624669980@163.com
 * @LastEditTime: 2023-02-14 11:28:35
 * @FilePath: /GetLastTag/main.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type GithubRelease struct {
	TagName string   `json:"tag_name"`
	Name    string   `json:"name"`
	Assets  []Assets `json:"assets"`
}
type Assets struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func GetLatestTag(repo string) (GithubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases?per_page=1", repo)
	resp, err := http.Get(url)
	if err != nil {
		return GithubRelease{}, err
	}
	defer resp.Body.Close()

	var tag []GithubRelease
	if err := json.NewDecoder(resp.Body).Decode(&tag); err != nil {
		return GithubRelease{}, err
	}

	if len(tag) > 0 {
		return tag[0], nil
	} else {
		return GithubRelease{}, fmt.Errorf("no tags found for %s", repo)
	}
}

type Config struct {
	Repos []string `json:"repos"`
	Arch  []string `json:"arch"`
}

const CONFIG = "/etc/config/config.json"

func main() {

	// 检查文件是否存在
	if _, err := os.Stat(CONFIG); os.IsNotExist(err) {
		// 创建文件
		file, err := os.Create(CONFIG)
		if err != nil {
			fmt.Println(err)
			return
		}
		file.Write([]byte(`{
			"repos": ["IceWhaleTech/CasaOS","IceWhaleTech/CasaOS-AppManagement","IceWhaleTech/CasaOS-UI","IceWhaleTech/CasaOS-MessageBus","IceWhaleTech/CasaOS-LocalStorage","IceWhaleTech/CasaOS-UserService","IceWhaleTech/CasaOS-Gateway","IceWhaleTech/CasaOS-CLI"],
			"arch":["amd64","all"]
		}`))
		defer file.Close()
		fmt.Println("文件已创建")
	} else {
		fmt.Println("文件已存在")
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		configData, err := ioutil.ReadFile(CONFIG)
		if err != nil {
			fmt.Println("Error reading config file:", err)
			w.Write([]byte("Error reading config file:" + err.Error()))
			return
		}
		var config Config
		if err := json.Unmarshal(configData, &config); err != nil {
			fmt.Println("Error parsing config file:", err)
			w.Write([]byte("Error parsing config file:" + err.Error()))
			return
		}
		result := []GithubRelease{}
		for _, v := range config.Repos {
			r, err := GetLatestTag(v)
			if err != nil {
				fmt.Println(err)
				w.Write([]byte(err.Error()))
				return
			} else {
				tempAssets := []Assets{}
				for i := range r.Assets {
					for _, arch := range config.Arch {
						if !strings.Contains(r.Assets[i].Name, "migration") && strings.Contains(r.Assets[i].Name, arch) {
							tempAssets = append(tempAssets, r.Assets[i])
						}
					}
				}
				r.Assets = tempAssets
				result = append(result, r)
			}
		}
		jsonBytes, err := json.Marshal(result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(jsonBytes)
	})
	http.ListenAndServe(":8089", nil)
}

// ...
