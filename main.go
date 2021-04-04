package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/valyala/fasthttp"
)

var Owner = ""       // Owner
var Repository = ""  // github repository name
var Path = "storage" // folder in the repository
var Token = ""       // github token: https://github.com/settings/tokens/new

var BaseUrl = "https://api.github.com/repos/" + Owner + "/" + Repository + "/contents/" + Path

func Upload(srcFilePath string, dstFilePath string, message string, sha string) bool {
	if !strings.HasPrefix(dstFilePath, "/") {
		dstFilePath = "/" + dstFilePath
	}

	url := BaseUrl + "/" + dstFilePath

	content, err := ioutil.ReadFile(srcFilePath)
	if err != nil {
		log.Println("Failed to read src", err)
		return false
	}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.Header.SetMethod("PUT")
	req.Header.SetBytesKV([]byte("Content-Type"), []byte("application/json"))
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/vnd.github.v3+json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+Token))

	req.SetRequestURI(url)

	args := make(map[string]string)
	args["content"] = b64.StdEncoding.EncodeToString(content)
	if len(message) > 0 {
		args["message"] = message
	} else {
		args["message"] = "Upload " + dstFilePath
	}
	if len(sha) != 0 {
		args["sha"] = sha
	}

	jsonBytes, _ := json.Marshal(args)
	req.SetBodyRaw(jsonBytes)

	if err := fasthttp.Do(req, resp); err != nil {
		log.Println("Failed to upload", url, err.Error())
	}

	body := resp.Body()

	var mapResult map[string]interface{}

	err = json.Unmarshal(body, &mapResult)
	if err != nil {
		log.Println("Failed to parse response", err)
		return false
	}

	return true
}

func RetrieveFile(filePath string) map[string]interface{} {
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}

	url := BaseUrl + filePath

	body := GetFiles(url)

	var mapResult map[string]interface{}

	err := json.Unmarshal(body, &mapResult)
	if err != nil {
		log.Println("Failed to parse response", err)
	}

	return mapResult
}

func RetrieveFiles() []map[string]interface{} {
	body := GetFiles(BaseUrl)

	var mapResult []map[string]interface{}

	err := json.Unmarshal(body, &mapResult)
	if err != nil {
		log.Println("Failed to parse responses", err)
	}

	return mapResult
}

func GetFiles(url string) []byte {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.Header.SetMethod("GET")
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/vnd.github.v3+json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+Token))

	req.SetRequestURI(url)

	if err := fasthttp.Do(req, resp); err != nil {
		log.Println("Failed to get files", url, err.Error())
	}

	return resp.Body()
}

func DeleteFile(filePath, sha, message string) string {
	if !strings.HasPrefix(filePath, "/") {
		filePath = "/" + filePath
	}

	url := BaseUrl + filePath

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer func() {
		fasthttp.ReleaseResponse(resp)
		fasthttp.ReleaseRequest(req)
	}()

	req.Header.SetMethod("DELETE")
	req.Header.SetBytesKV([]byte("Content-Type"), []byte("application/json"))
	req.Header.SetBytesKV([]byte("Accept"), []byte("application/vnd.github.v3+json"))
	req.Header.SetBytesKV([]byte("Authorization"), []byte("token "+Token))

	req.SetRequestURI(url)

	args := make(map[string]string)
	args["sha"] = sha
	if len(message) > 0 {
		args["message"] = message
	} else {
		args["message"] = "Delete " + filePath
	}

	jsonBytes, _ := json.Marshal(args)
	req.SetBodyRaw(jsonBytes)

	if err := fasthttp.Do(req, resp); err != nil {
		log.Println("Failed to delete file", url, err.Error())
	}

	return string(resp.Body())
}

func interfaceToString(input interface{}) string {
	return fmt.Sprintf("%v", input)
}

func main() {
	uploadPtr := flag.Bool("upload", false, "upload a file to github")
	retrievePtr := flag.Bool("retrieve", false, "retrieve a file or files from github")
	deletePtr := flag.Bool("delete", false, "delete a file in the repository")
	updatePtr := flag.Bool("update", false, "update a file in the repository")

	messagePtr := flag.String("message", "", "git commit message")
	srcFilePathPtr := flag.String("src", "", "the filepath to be uploaded/retrieved/deleted")
	destFilePathPtr := flag.String("dst", "", "the dest filepath in the repository")

	flag.Parse()

	if *uploadPtr {
		if len(*srcFilePathPtr) == 0 {
			log.Println("src parameter is missing")
			return
		}

		log.Println("Uploading")
		uploaded := Upload(*srcFilePathPtr, *destFilePathPtr, *messagePtr, "")
		log.Println("Uploaded", uploaded)
	} else if *retrievePtr && len(*destFilePathPtr) == 0 {
		log.Println("Retrieving")
		results := RetrieveFiles()
		for _, f := range results {
			//for k, v := range f {
			//    log.Println(k, v)
			//}
			//log.Println(f["html_url"], f["sha"])
			log.Println(f["name"])
		}
		log.Println("Retrieved")
	} else {
		if len(*destFilePathPtr) == 0 {
			log.Println("dst parameter is missing")
			return
		}

		result := RetrieveFile(*destFilePathPtr)
		if result == nil {
			log.Println("File not found", *destFilePathPtr)
			return
		}

		if *retrievePtr {
			log.Println("Retrieving")
			//log.Println(result["html_url"], result["sha"])
			content, _ := b64.StdEncoding.DecodeString(interfaceToString(result["content"]))
			log.Println(string(content))
			log.Println("Retrieved")
		} else if *deletePtr {
			log.Println("Deleting")
			resp := DeleteFile(*destFilePathPtr, interfaceToString(result["sha"]), "")
			log.Println("Deleted", resp)
		} else if *updatePtr {
			if len(*srcFilePathPtr) == 0 {
				log.Println("src parameter is missing")
				return
			}

			log.Println("Updating")
			uploaded := Upload(*srcFilePathPtr, *destFilePathPtr, *messagePtr, interfaceToString(result["sha"]))
			log.Println("Uploaded", uploaded)
		} else {
			log.Println("Wrong parameters")
		}
	}
}
