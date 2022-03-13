package resb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"

	"mmm3w/sparking/support"
)

func findThumb(folder string, parent string) string {
	dir, err := ioutil.ReadDir(folder)
	if err != nil {
		return ""
	}
	for _, fi := range dir {
		if !fi.IsDir() {
			return path.Join(parent, fi.Name())
		}
	}
	return ""
}

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
				dirNode["thumb"] = findThumb(absolutePath, relativePath)
				dirs = append(dirs, dirNode)
			} else {
				relativePath := path.Join(sbd, fi.Name())
				absolutePath := path.Join(targetFolder, fi.Name())

				fileNode := make(map[string]interface{})
				fileNode["name"] = fi.Name()
				fileNode["relative"] = relativePath
				fileNode["absolute"] = absolutePath
				fileNode["suffix"] = path.Ext(fi.Name())
				fileNode["size"] = fi.Size()
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
