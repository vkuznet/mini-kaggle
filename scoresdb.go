package main

// ScoresDB module
//
// Copyright (c) 2019 - Valentin Kuznetsov <vkuznet@gmail.com>
//

// for Go database API: http://go-database-sql.org/overview.html
// tutorial: https://golang-basic.blogspot.com/2014/06/golang-database-step-by-step-guide-on.html
// Oracle drivers:
//   _ "gopkg.in/rana/ora.v4"
//   _ "github.com/mattn/go-oci8"
// MySQL driver:
//   _ "github.com/go-sql-driver/mysql"
// SQLite driver:
//  _ "github.com/mattn/go-sqlite3"
//

import (
	"database/sql"
	"errors"
	"log"
	"sort"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// global variable to keep pointer to ScoresDB
var ScoresDB *sql.DB

// InitScoresDB sets pointer to ScoresDB
func InitScoresDB(uri string) (*sql.DB, error) {
	dbAttrs := strings.Split(uri, "://")
	if len(dbAttrs) != 2 {
		return nil, errors.New("Please provide proper ScoresDB uri")
	}
	db, err := sql.Open(dbAttrs[0], dbAttrs[1])
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(100)

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS scores (id INTEGER PRIMARY KEY, name TEXT, score FLOAT, private FLOAT)")
	if err != nil {
		return nil, err
	}
	stmt.Exec()
	return db, err
}

// ScoreRecord represent score record
type ScoreRecord struct {
	Name  string
	Score float64
}

// ScoreRecordsList implement sort for []int type
type ScoreRecordsList []ScoreRecord

func (s ScoreRecordsList) Len() int           { return len(s) }
func (s ScoreRecordsList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ScoreRecordsList) Less(i, j int) bool { return s[i].Score < s[j].Score }

// GetScores fetches all scores from our DB
func GetScores(stype string) []ScoreRecord {
	stmt := "select name, score from scores"
	if stype == "private" {
		stmt = "select name, private from scores"
	}
	rows, err := ScoresDB.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var records ScoreRecordsList
	for rows.Next() {
		var name string
		var score float64
		err := rows.Scan(&name, &score)
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, ScoreRecord{Name: name, Score: score})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	sort.Sort(sort.Reverse(records))
	return records

}

// InsertScore insert given files into ScoresDB
func InsertScore(name string, score, private float64) error {
	// proceed with transaction operation
	tx, err := ScoresDB.Begin()
	if err != nil {
		log.Println("DB error", err)
		return err
	}
	defer tx.Rollback()

	var stmt string
	// insert main attributes
	stmt = "INSERT INTO scores (name, score, private) VALUES (?, ?, ?)"
	_, err = tx.Exec(stmt, name, score, private)
	if err != nil {
		tx.Rollback()
		return err
	}
	// commit whole workflow
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
