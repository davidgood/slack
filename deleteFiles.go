package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

var token = url.QueryEscape("YOUR TOKEN HERE")

const uri = "https://{YOURDOMAIN}.slack.com/api/%s?token=%s"

var authURL = "auth.test?token="

var userID = ""

func getAPI(url string) (result map[string]interface{}, err error) {
	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Unable to GET: %s\n", err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println(err)
		return nil, err
	}

	return
}

func authenticate() (bool, error) {
	url := fmt.Sprintf(uri, "auth.test", token)

	auth, err := getAPI(url)
	if err != nil {
		return false, err
	}

	ok := auth["ok"].(bool)
	if ok == false {
		return false, nil
	}

	userID := auth["user_id"].(string)
	if userID == "" {
		return false, nil
	}

	return true, nil
}

func getFileList() (map[string]interface{}, error) {
	url := fmt.Sprintf(uri+"&user="+userID, "files.list", token)

	files, err := getAPI(url)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func deleteFile(id string) bool {
	loc := fmt.Sprintf(uri, "files.delete", token)

	client := &http.Client{Timeout: time.Second * 10}

	v := url.Values{}
	v.Set("token", token)
	v.Set("file", id)

	resp, err := client.PostForm(loc, v)
	if err != nil {
		log.Printf("Unable to POST: %s", err)
	}
	defer resp.Body.Close()

	fd := make(map[string]interface{})
	if err := json.NewDecoder(resp.Body).Decode(&fd); err != nil {
		log.Println(err)
	}

	return fd["ok"].(bool)
}

func main() {

	ok, err := authenticate()
	if err != nil || ok == false {
		log.Fatal("Unable to authenticate\n")
	}

	fileList, err := getFileList()
	if err != nil {
		log.Fatal("Unable to retrieve file list")
	}
	if len(fileList) == 0 {
		log.Fatal("File list is empty")
	}

	fmt.Printf("Found %d files\n", len(fileList))

	files := fileList["files"].([]interface{})
	for i := 0; i < len(files); i++ {
		deleteFile(files[i].(map[string]interface{})["id"].(string))
	}

}
