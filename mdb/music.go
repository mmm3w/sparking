package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/rs/xid"
)

/********************************************/
/** 歌曲的增删改查 ***************************/
/********************************************/

func saveMusic(data string, db *sql.DB) (int, string) {
	var item EditBean
	err := json.Unmarshal([]byte(data), &item)
	if err != nil {
		return 500, err.Error()
	}

	tx, _ := db.Begin()

	//插入Music
	item.Id = xid.New().String()
	_, err = tx.Exec("INSERT INTO music(id, name, lrc, is_lossless, meta) values(?,?,?,?,?)",
		item.Id, item.Name, item.Lrc, item.Is_lossless, item.Meta)
	if err != nil {
		tx.Rollback()
		return 500, err.Error()
	}

	//插入artist
	for i := 0; i < len(item.Artist); i++ {
		_, err = tx.Exec("INSERT INTO artist_association(id, mid, aid) values(?,?,?)",
			xid.New().String(), item.Id, item.Artist[i])
		if err != nil {
			tx.Rollback()
			return 500, err.Error()
		}
	}

	//插入list
	for i := 0; i < len(item.List); i++ {
		_, err = tx.Exec("INSERT INTO list_association(id, mid, lid) values(?,?,?)",
			xid.New().String(), item.Id, item.List[i])
		if err != nil {
			tx.Rollback()
			return 500, err.Error()
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 500, err.Error()
	}
	fmt.Println(item)

	return 200, item.Id
}

func deleteMusic(params url.Values, db *sql.DB) (int, string) {
	if len(params["id"]) == 0 {
		return 500, "No id"
	}

	_, err := db.Exec("DELETE FROM music where id=?", params["id"][0])
	if err != nil {
		return 500, err.Error()
	}

	return 200, ""
}

func editMusic(data string, db *sql.DB) (int, string) {
	var item EditBean
	err := json.Unmarshal([]byte(data), &item)
	if err != nil {
		return 500, err.Error()
	}

	tx, _ := db.Begin()

	//更新核心内容
	_, err = tx.Exec("UPDATE music set name=? lrc=? is_lossless=? meta=? WHERE id=?",
		item.Name, item.Lrc, item.Is_lossless, item.Meta, item.Id)
	if err != nil {
		tx.Rollback()
		return 500, err.Error()
	}

	_, err = db.Exec("DELETE FROM artist_association where mid=?", item.Id)
	if err != nil {
		tx.Rollback()
		return 500, err.Error()
	}
	for i := 0; i < len(item.Artist); i++ {
		_, err = tx.Exec("INSERT INTO artist_association(id, mid, aid) values(?,?,?)",
			xid.New().String(), item.Id, item.Artist[i])
		if err != nil {
			tx.Rollback()
			return 500, err.Error()
		}
	}

	_, err = db.Exec("DELETE FROM list_association where mid=?", item.Id)
	if err != nil {
		tx.Rollback()
		return 500, err.Error()
	}
	for i := 0; i < len(item.List); i++ {
		_, err = tx.Exec("INSERT INTO list_association(id, mid, lid) values(?,?,?)",
			xid.New().String(), item.Id, item.List[i])
		if err != nil {
			tx.Rollback()
			return 500, err.Error()
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return 500, err.Error()
	}
	fmt.Println(item)

	return 200, ""
}

func queryMusic(params url.Values, db *sql.DB) (int, string) {
	/**
	模糊查询艺术家和名字，歌词过滤，无损过滤，有无源过滤，
	*/

	// var index int
	// var pageSize int
	// var err error

	// if len(params["index"]) == 0 {
	// 	index = 0
	// } else {
	// 	i := params["index"][0]
	// 	index, err = strconv.Atoi(i)
	// 	if err != nil {
	// 		index = 0
	// 	}
	// }

	// if len(params["size"]) == 0 {
	// 	pageSize = 0
	// } else {
	// 	i := params["size"][0]
	// 	pageSize, err = strconv.Atoi(i)
	// 	if err != nil {
	// 		pageSize = 0
	// 	}
	// }

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM music").Scan(&count)
	if err != nil {
		return 500, err.Error()
	}

	musicRows, err := db.Query("SELECT id, name, lrc, is_lossless, meta FROM music")
	if err != nil {
		return 500, err.Error()
	}
	defer musicRows.Close()

	musicArray := make([]Music, count)
	musicCounter := 0
	for musicRows.Next() {
		var item Music
		err = musicRows.Scan(&item.Id, &item.Name, &item.Lrc, &item.Is_lossless, &item.Meta)
		if err != nil {
			return 500, err.Error()
		}

		musicArray[musicCounter] = item
		musicCounter++
	}

	err = musicRows.Err()
	if err != nil {
		return 500, err.Error()
	}

	fmt.Println(musicArray)

	jsonBytes, err := json.Marshal(musicArray)
	if err != nil {
		return 500, err.Error()
	}
	return 200, string(jsonBytes)
}
