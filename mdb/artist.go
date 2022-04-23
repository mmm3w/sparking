package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/rs/xid"
)

/********************************************/
/** 艺术家的增删改查 *************************/
func saveArtist(params url.Values, db *sql.DB) (int, string) {
	if len(params["name"]) == 0 {
		return 500, "No name"
	}

	id := xid.New().String()
	_, err := db.Exec("INSERT INTO artist(id, name) values(?,?)", id, params["name"][0])
	if err != nil {
		return 500, err.Error()
	}

	return 200, id
}

func deleteArtist(params url.Values, db *sql.DB) (int, string) {
	if len(params["id"]) == 0 {
		return 500, "No id"
	}

	_, err := db.Exec("DELETE FROM artist where id=?", params["id"][0])
	if err != nil {
		return 500, err.Error()
	}

	return 200, ""
}

func editArtist(data string, db *sql.DB) (int, string) {
	var item Artist

	err := json.Unmarshal([]byte(data), &item)
	if err != nil {
		return 500, err.Error()
	}

	_, err = db.Exec("UPDATE music artist name=? WHERE id=?", item.Name, item.Id)
	if err != nil {
		return 500, err.Error()
	}

	fmt.Println(item)
	return 200, ""
}

func queryArtist(params url.Values, db *sql.DB) (int, string) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM artist").Scan(count)
	if err != nil {
		return 500, err.Error()
	}

	artistRows, err := db.Query("SELECT id, name FROM artist")
	if err != nil {
		return 500, err.Error()
	}
	defer artistRows.Close()

	artistArray := make([]Artist, count)
	artistCounter := 0
	for artistRows.Next() {
		var item Artist
		err = artistRows.Scan(item.Id, item.Name)
		if err != nil {
			return 500, err.Error()
		}
		artistArray[artistCounter] = item
		artistCounter++
	}

	err = artistRows.Err()
	if err != nil {
		return 500, err.Error()
	}

	fmt.Println(artistArray)

	jsonBytes, err := json.Marshal(artistArray)
	if err != nil {
		return 500, err.Error()
	}
	return 200, string(jsonBytes)
}
