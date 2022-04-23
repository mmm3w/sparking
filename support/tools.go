package support

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"net/url"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func EnsureDir(path string) {
	if Exists(path) {
		if !IsDir(path) {
			os.Remove(path)
			os.MkdirAll(path, os.ModePerm)
		}
	} else {
		os.MkdirAll(path, os.ModePerm)
	}
}

func GetValue(values url.Values, key string, def string) string {
	if vs := values[key]; len(vs) > 0 {
		return vs[0]
	} else {
		return def
	}
}

func Write(path string, data string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	_, err = w.WriteString(data)
	if err != nil {
		w.Flush()
		return err
	}
	return w.Flush()
}

func Read(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func NoErrorBase64(source string) string {
	temp, err := base64.RawURLEncoding.DecodeString(source)
	if err != nil {
		return ""
	}
	return string(temp)
}
