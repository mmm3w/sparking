package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/rs/xid"
)

/********************************************/
/** 组列的增删改查 ***************************/
func saveList(params url.Values, db *sql.DB) (int, string) {
	if len(params["tag"]) == 0 {
		return 500, "No tag"
	}

	id := xid.New().String()
	_, err := db.Exec("INSERT INTO list(id, tag) values(?,?)", id, params["tag"][0])
	if err != nil {
		return 500, err.Error()
	}

	return 200, id
}

func deleteList(params url.Values, db *sql.DB) (int, string) {
	if len(params["id"]) == 0 {
		return 500, "No id"
	}

	_, err := db.Exec("DELETE FROM list where id=?", params["id"][0])
	if err != nil {
		return 500, err.Error()
	}

	return 200, ""
}

func editList(data string, db *sql.DB) (int, string) {
	var item List
	err := json.Unmarshal([]byte(data), &item)
	if err != nil {
		return 500, err.Error()
	}

	_, err = db.Exec("UPDATE music list tag=? WHERE id=?", item.Tag, item.Id)
	if err != nil {
		return 500, err.Error()
	}

	fmt.Println(item)
	return 200, ""
}

func queryList(params url.Values, db *sql.DB) (int, string) {

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM list").Scan(count)
	if err != nil {
		return 500, err.Error()
	}

	listRows, err := db.Query("SELECT id, tag FROM list")
	if err != nil {
		return 500, err.Error()
	}
	defer listRows.Close()

	listArray := make([]List, count)
	listCounter := 0
	for listRows.Next() {
		var item List
		err = listRows.Scan(item.Id, item.Tag)
		if err != nil {
			return 500, err.Error()
		}
		listArray[listCounter] = item
		listCounter++
	}

	err = listRows.Err()
	if err != nil {
		return 500, err.Error()
	}

	fmt.Println(listArray)

	jsonBytes, err := json.Marshal(listArray)
	if err != nil {
		return 500, err.Error()
	}
	return 200, string(jsonBytes)
}
