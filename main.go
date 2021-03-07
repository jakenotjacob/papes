package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/google/go-querystring/query"
	_ "io"
	"log"
	_ "net/http"
	"os"
)

// Query is the main search string entered, such as "forest"
// Resolution retricts the result set to images of dimension:  "1920x1080", etc.
// Sorting allows us to choose the result set ordering, so we don't have to shuffle :)
type Params struct {
	Query      string `url:"q"`
	Resolution string `url:"resolutions"`
	Sorting    string `url:"sorting"`
}

// Image contains identifying data and metadatas
// about the image, such as: id, url, views, favorites, dimension, file_size, colors, etc.
// TODO Fill out, we just don't need to yet
type Image struct {
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

	//params := Params{Query: "forest", Resolution: "1920x1080", Sorting: "random"}
	//v, _ := query.Values(params)

	// resp, err := http.Get(BaseURL + v.Encode())
	// if err != nil {
	// 	log.Fatalf("Failed to read URL: %v", BaseURL)
	// }
	// defer resp.Body.Close()

	// TODO To avoid spamming wallhaven we just save the json blob
	// into a local file for testing since the interface should be
	// a reader/writer anyways, so we should be able to read both ez.
	// images, err := io.ReadAll(resp.Body)
	blob, err := os.ReadFile("blob")
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
}
