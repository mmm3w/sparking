package main

type Artist struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type List struct {
	Id  string `json:"id"`
	Tag string `json:"tag"`
}

type Music struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Lrc         string   `json:"lrc"`
	Meta        string   `json:"meta"`
	Is_lossless bool     `json:"is_lossless"`
	Artist      []Artist `json:"artist"`
	List        []List   `json:"list"`
}

type EditBean struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Lrc         string   `json:"lrc"`
	Meta        string   `json:"meta"`
	Is_lossless bool     `json:"is_lossless"`
	Artist      []string `json:"artist"`
	List        []string `json:"list"`
}

type Paging struct {
	Index     int `json:"index"`
	Page_size int `json:"page_size"`
}
