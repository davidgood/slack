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

const uri = "YOUR DOMAIN HERE.slack.com/api/%s?token=%s"

var authURL = "auth.test?token="

var userId = ""

type auth struct {
	Ok     bool   `json:"ok"`
	URL    string `json:"url"`
	Team   string `json:"team"`
	User   string `json:"user"`
	TeamID string `json:"team_id"`
	UserID string `json:"user_id"`
}

type fileDelete struct {
	Ok bool `json:"ok"`
}

type file struct {
	ID                 string        `json:"id"`
	Created            int           `json:"created"`
	Timestamp          int           `json:"timestamp"`
	Name               string        `json:"name"`
	Title              string        `json:"title"`
	Mimetype           string        `json:"mimetype"`
	Filetype           string        `json:"filetype"`
	PrettyType         string        `json:"pretty_type"`
	User               string        `json:"user"`
	Editable           bool          `json:"editable"`
	Size               int           `json:"size"`
	Mode               string        `json:"mode"`
	IsExternal         bool          `json:"is_external"`
	ExternalType       string        `json:"external_type"`
	IsPublic           bool          `json:"is_public"`
	PublicURLShared    bool          `json:"public_url_shared"`
	DisplayAsBot       bool          `json:"display_as_bot"`
	Username           string        `json:"username"`
	URLPrivate         string        `json:"url_private"`
	URLPrivateDownload string        `json:"url_private_download"`
	Thumb64            string        `json:"thumb_64"`
	Thumb80            string        `json:"thumb_80"`
	Thumb360           string        `json:"thumb_360"`
	Thumb360W          int           `json:"thumb_360_w"`
	Thumb360H          int           `json:"thumb_360_h"`
	Thumb480           string        `json:"thumb_480"`
	Thumb480W          int           `json:"thumb_480_w"`
	Thumb480H          int           `json:"thumb_480_h"`
	Thumb160           string        `json:"thumb_160"`
	Thumb720           string        `json:"thumb_720"`
	Thumb720W          int           `json:"thumb_720_w"`
	Thumb720H          int           `json:"thumb_720_h"`
	Thumb800           string        `json:"thumb_800"`
	Thumb800W          int           `json:"thumb_800_w"`
	Thumb800H          int           `json:"thumb_800_h"`
	Thumb960           string        `json:"thumb_960"`
	Thumb960W          int           `json:"thumb_960_w"`
	Thumb960H          int           `json:"thumb_960_h"`
	Thumb1024          string        `json:"thumb_1024"`
	Thumb1024W         int           `json:"thumb_1024_w"`
	Thumb1024H         int           `json:"thumb_1024_h"`
	ImageExifRotation  int           `json:"image_exif_rotation"`
	OriginalW          int           `json:"original_w"`
	OriginalH          int           `json:"original_h"`
	Permalink          string        `json:"permalink"`
	PermalinkPublic    string        `json:"permalink_public"`
	Channels           []string      `json:"channels"`
	Groups             []interface{} `json:"groups"`
	Ims                []interface{} `json:"ims"`
	CommentsCount      int           `json:"comments_count"`
	InitialComment     struct {
		ID        string `json:"id"`
		Created   int    `json:"created"`
		Timestamp int    `json:"timestamp"`
		User      string `json:"user"`
		IsIntro   bool   `json:"is_intro"`
		Comment   string `json:"comment"`
	} `json:"initial_comment"`
}

type fileList struct {
	Ok     bool `json:"ok"`
	Files  []*file
	Paging struct {
		Count int `json:"count"`
		Total int `json:"total"`
		Page  int `json:"page"`
		Pages int `json:"pages"`
	} `json:"paging"`
}

// This is somewhat unsafe - any errors with generic inteface will
// only show up at runtime
func getApi(url string, r interface{}) {
	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Unable to GET: %s", err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(r); err != nil {
		log.Println(err)
	}
}

func authenticate() bool {
	url := fmt.Sprintf(uri, "auth.test", token)

	a := auth{}

	getApi(url, &a)

	userId = a.UserID

	return a.Ok
}

func getFileList() fileList {
	url := fmt.Sprintf(uri+"&user="+userId, "files.list", token)

	files := fileList{}

	getApi(url, &files)

	return files
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

	fd := &fileDelete{}

	if err := json.NewDecoder(resp.Body).Decode(fd); err != nil {
		log.Println(err)
	}

	return fd.Ok
}

func main() {

	if !authenticate() {
		log.Fatal("Unable to authenticate")
	}

	files := getFileList()
	if len(files.Files) == 0 {
		log.Fatal("File list is empty")
	}

	fmt.Printf("Found %d files\n", len(files.Files))

	for _, f := range files.Files {
		fmt.Printf("Deleting: %s", f.Name)
		deleteFile(f.ID)
	}

}
