package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func currentFolder() string {
	folder, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return folder
}

func printError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func initCreate(db *sql.DB) {

	db.Exec("PRAGMA foreign_keys = ON")

	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS music(
		id			TEXT	PRIMARY KEY	NOT NULL,
		name		TEXT	NOT NULL,
		lrc			TEXT	NOT NULL,
		is_lossless	BOOLEAN	NOT NULL,
		meta		TEXT	NOT NULL
    );
	`)
	printError(err)

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS artist(
		id		TEXT	PRIMARY KEY	NOT NULL,
		name	TEXT	NOT NULL	UNIQUE
	);
	`)
	printError(err)

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS list(
		id		TEXT	PRIMARY KEY	NOT NULL,
		tag		TEXT	NOT NULL	UNIQUE
	);
	`)
	printError(err)

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS artist_association(
		id	TEXT	PRIMARY KEY,
		mid	TEXT	NOT NULL,
		aid	TEXT NOT NULL,
		FOREIGN KEY(mid) REFERENCES music(id) ON DELETE CASCADE,
		FOREIGN KEY(aid) REFERENCES artist(id) ON DELETE CASCADE
	);
	`)
	printError(err)

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS list_association(
		id	TEXT	PRIMARY KEY,
		mid	TEXT	NOT NULL,
		lid	TEXT NOT NULL,
		FOREIGN KEY(mid) REFERENCES music(id) ON DELETE CASCADE,
		FOREIGN KEY(lid) REFERENCES list(id) ON DELETE CASCADE
	);
	`)
	printError(err)
}

func main() {
	//确定目录
	workFolder := currentFolder()
	fmt.Println("Work folder : " + workFolder)

	// db, err := sql.Open("sqlite3", path.Join(workFolder, "music.db"))
	db, err := sql.Open("sqlite3", "./music.db")
	defer db.Close()
	printError(err)
	initCreate(db)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
	})

	http.HandleFunc("/music", func(w http.ResponseWriter, r *http.Request) {
		var code int
		var msg string

		if r.Method == "GET" {
			r.ParseForm()
			code, msg = queryMusic(r.Form, db)
		} else if r.Method == "POST" {
			defer r.Body.Close()
			con, err := ioutil.ReadAll(r.Body)
			if err != nil {
				code, msg = 500, err.Error()
			} else {
				code, msg = editMusic(string(con), db)
			}
		} else if r.Method == "PUT" {
			defer r.Body.Close()
			con, err := ioutil.ReadAll(r.Body)
			if err != nil {
				code, msg = 500, err.Error()
			} else {
				code, msg = saveMusic(string(con), db)
			}
		} else if r.Method == "DELETE" {
			r.ParseForm()
			code, msg = deleteMusic(r.Form, db)
		} else {
			code, msg = 403, ""
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(code)
		fmt.Fprintf(w, msg)
	})

	http.HandleFunc("/artist", func(w http.ResponseWriter, r *http.Request) {
		var code int
		var msg string

		if r.Method == "GET" {
			r.ParseForm()
			code, msg = queryArtist(r.Form, db)
		} else if r.Method == "POST" {
			defer r.Body.Close()
			con, err := ioutil.ReadAll(r.Body)
			if err != nil {
				code, msg = 500, err.Error()
			} else {
				code, msg = editList(string(con), db)
			}
		} else if r.Method == "PUT" {
			r.ParseForm()
			code, msg = saveArtist(r.Form, db)
		} else if r.Method == "DELETE" {
			r.ParseForm()
			code, msg = deleteArtist(r.Form, db)
		} else {
			code, msg = 403, ""
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(code)
		fmt.Fprintf(w, msg)
	})

	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		var code int
		var msg string

		if r.Method == "GET" {
			r.ParseForm()
			code, msg = queryList(r.Form, db)
		} else if r.Method == "POST" {
			defer r.Body.Close()
			con, err := ioutil.ReadAll(r.Body)
			if err != nil {
				code, msg = 500, err.Error()
			} else {
				code, msg = editList(string(con), db)
			}
		} else if r.Method == "PUT" {
			r.ParseForm()
			code, msg = saveList(r.Form, db)
		} else if r.Method == "DELETE" {
			r.ParseForm()
			code, msg = deleteList(r.Form, db)
		} else {
			code, msg = 403, ""
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(code)
		fmt.Fprintf(w, msg)
	})

	//监听相应端口
	err = http.ListenAndServe(":23323", nil)
	printError(err)
}
