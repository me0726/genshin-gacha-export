package main

import (
	"database/sql"
	"log"
)

const DDL = `CREATE TABLE IF NOT EXISTS gacha
					(
					    "id"         varchar(100),
					    "uid"        varchar(100),
					    "gacha_type" varchar(100),
					    "item_id"    varchar(100),
					    "count"      varchar(100),
					    "time"       varchar(100),
					    "name"       varchar(100),
					    "lang"       varchar(100),
					    "item_type"  varchar(100),
					    "rank_type"  varchar(100)
					)`

func initSqlite() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./data.sqlite3")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(DDL)
	if err != nil {
		return nil, err
	}
	return db, nil
}

const INSERT = `INSERT INTO gacha(id, uid, gacha_type, item_id, count, time, name, lang, item_type, rank_type)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

func insert(db *sql.DB, data RealData) error {
	prepare, err := db.Prepare(INSERT)
	if err != nil {
		return err
	}
	_, err = prepare.Exec(data.Id,
		data.Uid,
		data.GachaType,
		data.ItemId,
		data.Count,
		data.Time,
		data.Name,
		data.Lang,
		data.ItemType,
		data.RankType)
	if err != nil {
		return err
	}
	return nil
}

const SELECT = `SELECT * FROM gacha WHERE gacha_type=? ORDER BY time `

func findAllByGachaType(db *sql.DB, gachaType int) (RealDataList, error) {
	var list = RealDataList{}
	prepare, err := db.Prepare(SELECT)
	if err != nil {
		return nil, err
	}
	res, err := prepare.Query(gachaType)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		var data = RealData{}
		err = res.Scan(&data.Id,
			&data.Uid,
			&data.GachaType,
			&data.ItemId,
			&data.Count,
			&data.Time,
			&data.Name,
			&data.Lang,
			&data.ItemType,
			&data.RankType)
		if err != nil {
			log.Print(err)
			continue
		}
		list = append(list, data)
	}
	return list, nil
}

const ExistsById = "SELECT * FROM gacha WHERE id=?"

func existsById(id string, db *sql.DB) (bool, error) {
	prepare, err := db.Prepare(ExistsById)
	if err != nil {
		return false, err
	}
	query, err := prepare.Query(id)
	if err != nil {
		return false, err
	}
	defer query.Close()
	return query.Next(), nil
}
