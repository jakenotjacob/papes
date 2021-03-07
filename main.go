package main

import (
	"encoding/json"
	"fmt"
	"io"
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
		fmt.Print("Need somethin to search for!")
	}

	tags := strings.Join(os.Args[1:], " ")

	params := Params{Query: tags, Resolution: "1920x1080", Sorting: "random", Page: 1}
	v, _ := query.Values(params)

	resp, err := http.Get(BaseURL + v.Encode())
	if err != nil {
		fmt.Printf("Failed to read URL: %v", BaseURL)
	}
	defer resp.Body.Close()

	// TODO To avoid spamming wallhaven we just save the json blob
	// into a local file for testing since the interface should be
	// a reader anyways, so we should be able to swap:
	//blob, err := os.ReadFile("/tmp/blob")
	var blob []byte
	blob, err = io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Unable to blob stream: %v", err)
	}

	var r Response
	if err := json.Unmarshal(blob, &r); err != nil {
		fmt.Printf("Couldn't parse anything meaningful: %v", err)
	}
	//For JSON arrays... []interface{}
	//For JSON objects.. map[string]interface{}
	if len(r.Data) <= 0 {
		fmt.Printf("No results, no work to do!")
		os.Exit(1)
	}
	// TODO concurrent downloads
	//for _, image := range r.Data {
	//	//fmt.Println(image.Path)
	//}

	// Fetch an image...
	res, err := http.Get(r.Data[0].Path)
	if err != nil {
		fmt.Printf("Failed to download image! %v ", err)
	}
	defer res.Body.Close()

	// Save temporarily...
	n := strings.Split(r.Data[0].Path, "/")
	imgpath := fmt.Sprintf("./%v", n[len(n)-1])
	tmpfile, err := os.Create(imgpath)
	_, err = io.Copy(tmpfile, res.Body)
	if err != nil {
		fmt.Printf("Failed to save image! %v", err)
	}
	// TODO This is not deferred on purpose because the
	// file write must complete before we can use it
	tmpfile.Close()

	// Finally, use $somecmd to set the background
	// TODO loop setting images from download set and save permanently
	// case <-finished/sig; when break: set() && save()
	ags := []string{"--bg-fill", imgpath}
	bgset := exec.Command("feh", ags...)
	if err := bgset.Run(); err != nil {
		fmt.Printf("Fucked up while tryna set the background: %v", err)
	}

}
