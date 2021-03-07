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

type Params struct {
	Query      string `url:"q"`
	Resolution string `url:"resolutions"`
	Sorting    string `url:"sorting"`
}

type Image struct {
	Path string `json:"path"`
}

type Meta struct {
	PageSize int `json:"per_page"`
	Total    int `json:"total"`
}

type Response struct {
	Data []Image `json:"data"`
	Meta Meta    `json:"meta"`
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
