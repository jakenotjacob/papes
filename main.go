package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/google/go-querystring/query"
)

// Query is the main search string entered, such as "forest"
// Resolution retricts the result set to images of dimension:  "1920x1080", etc.
// Sorting allows us to choose the result set ordering, so we don't have to shuffle :)
type Params struct {
	Query      string `url:"q"`
	Resolution string `url:"resolutions"`
	Sorting    string `url:"sorting"`
	Page       int    `url:"page"`
}

// Image contains identifying data and metadatas
// about the image, such as: id, url, views, favorites, dimension, file_size, colors, etc.
// TODO Fill out, we just don't need to yet
type Image struct {
	Id   string `json:"id"`
	Path string `json:"path"`
}

// The Meta struct contains query response metadata such
// as page size, total pages, current page, etc.
// TODO Not using all of the objects returned in the response
// cause we only really want the URL.
type Meta struct {
	PageSize int `json:"per_page"`
	Total    int `json:"total"`
}

// Response is a union of both the data and metadatas about that datas
// Does your data data its datas data?
type Response struct {
	Data []Image `json:"data"`
	// ♫ Do dooo do doodoo... ♫
	Meta Meta `json:"meta"`
}

func main() {
	const BaseURL string = "https://wallhaven.cc/api/v1/search?"

	if len(os.Args) <= 1 {
		log.Fatal("Need somethin to search for!")
	}

	tags := strings.Join(os.Args[1:], " ")
	fmt.Printf("Searching for... %v", tags)

	params := Params{Query: tags, Resolution: "1920x1080", Sorting: "random", Page: 1}
	v, _ := query.Values(params)

	resp, err := http.Get(BaseURL + v.Encode())
	if err != nil {
		log.Fatalf("Failed to read URL: %v", BaseURL)
	}
	defer resp.Body.Close()

	// TODO To avoid spamming wallhaven we just save the json blob
	// into a local file for testing since the interface should be
	// a reader/writer anyways, so we should be able to read both ez.
	//blob, err := os.ReadFile("/tmp/blob")
	var blob []byte
	blob, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var r Response
	if err := json.Unmarshal(blob, &r); err != nil {
		log.Fatal(err)
	}
	//For JSON arrays... []interface{}
	//For JSON objects.. map[string]interface{}
	for _, image := range r.Data {
		fmt.Println(image.Path)
	}

	// Fetch an image...
	res, err := http.Get(r.Data[1].Path)
	if err != nil {
		log.Fatalf("Failed to download image! %v ", err)
	}
	defer res.Body.Close()

	// Save temporarily...
	n := strings.Split(r.Data[1].Path, "/")
	imgpath := fmt.Sprintf("/tmp/%v", n[len(n)-1])
	tmpfile, err := os.Create(imgpath)
	_, err = io.Copy(tmpfile, res.Body)
	if err != nil {
		log.Fatalf("Failed to save image! %v", err)
	}
	tmpfile.Close()

	// Finally, use $somecmd to set the background
	ags := []string{"--bg-fill", imgpath}
	bgset := exec.Command("feh", ags...)
	if err := bgset.Run(); err != nil {
		log.Fatal("Fucked up while tryna set the background")
	}

}
