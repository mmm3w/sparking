package resb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"mmm3w/sparking/support"
)

func getData(md string, sbd string) (string, error) {
	if md == "" {
		return "", errors.New("main directory is empty")
	}

	targetFolder := path.Join(md, sbd)

	files := []map[string]interface{}{}
	dirs := []map[string]interface{}{}

	dir, err := ioutil.ReadDir(targetFolder)
	if err == nil {
		for _, fi := range dir {
			if fi.IsDir() {
				relativePath := path.Join(sbd, fi.Name())
				absolutePath := path.Join(targetFolder, fi.Name())

				dirNode := make(map[string]interface{})
				dirNode["name"] = fi.Name()
				dirNode["relative"] = relativePath
				dirNode["absolute"] = absolutePath
				dirNode["thumb"] = support.Exists(path.Join(absolutePath, ".thumb"))
				dirs = append(dirs, dirNode)
			} else {
				relativePath := path.Join(sbd, fi.Name())
				absolutePath := path.Join(targetFolder, fi.Name())
				thumbPath := path.Join(targetFolder, ".thumbcache", fi.Name()+".jpg")

				fileNode := make(map[string]interface{})
				fileNode["name"] = fi.Name()
				fileNode["relative"] = relativePath
				fileNode["absolute"] = absolutePath
				fileNode["suffix"] = path.Ext(fi.Name())
				fileNode["size"] = fi.Size()
				if support.Exists(thumbPath) {
					fileNode["thumb"] = path.Join(sbd, ".thumbcache", fi.Name()+".jpg")
				} else {
					fileNode["thumb"] = ""
				}

				files = append(files, fileNode)
			}
		}
	}

	data := make(map[string][]map[string]interface{})
	data["files"] = files
	data["dirs"] = dirs

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func Find(w http.ResponseWriter, r *http.Request) {
	var code int
	var message string
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)

		targetFolder := support.GetValue(r.PostForm, "folder", "")
		targetPath := support.GetValue(r.PostForm, "path", "")

		data, err := getData(targetFolder, targetPath)
		if err != nil {
			code, message = 500, err.Error()
		} else {
			w.Header().Set("content-type", "text/json")
			code, message = 200, data
		}
	} else {
		code, message = 403, "Error"
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	var code int
	var message string
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)

		dir := support.GetValue(r.PostForm, "dir", "")
		file, handler, err := r.FormFile("file")

		if err != nil {
			code, message = 500, err.Error()
		} else if dir == "" {
			defer file.Close()
			code, message = 500, "dir is empty"
		} else {
			defer file.Close()
			support.EnsureDir(dir)

			filePath := path.Join(dir, handler.Filename)
			if support.Exists(filePath) {
				os.Remove(filePath)
			}

			f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				code, message = 500, err.Error()
			} else {
				defer f.Close()
				io.Copy(f, file)
				code, message = 200, "Success"
			}
		}
	} else {
		code, message = 403, "Error"
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	fmt.Fprint(w, message)
}
