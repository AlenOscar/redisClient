/**
 * Created by Alen on 2019-05-20 14:38
 */

package models

import (
	"testing"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"encoding/json"
)

func checkErr(err error){
	if err != nil {
		panic(err)
	}
}

func TestDbWorker_Insert(t *testing.T) {
	Db.Insert(uint64(1), "s")
}

func TestDbWorker_Delete(t *testing.T) {

}

func TestDbWorker_Update(t *testing.T) {

}

func TestDbWorker_Query(t *testing.T) {
	Db.Query(uint64(1))
}

func TestProcedure(t *testing.T) {
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", DbUsername, DbPassword, DbHost, DbPort, DbName)
	db, err := sql.Open("mysql", dataSourceName)
	checkErr(err)

	type mailData struct {
		Item uint32 `json:"item"`
		Count uint32 `json:"count"`
		Desc string `json:"desc"`
	}

	tx, e := db.Begin()
	checkErr(e)

	// 最后释放tx内部的连接
	defer tx.Commit()

	data := mailData{
		Item: uint32(1000),
		Count: uint32(200),
		Desc: string("奖励邮件，请查收！"),
	}
	for i := 1; i <= 1000; i++ {
		jsonData, err := json.Marshal(data)
		checkErr(err)

		// 每次循环用的都是tx内部的连接，没有新建连接，效率高
		tx.Exec(fmt.Sprintf("CALL add_mail('%d','%d','%s')", i, 123456, jsonData))
	}

	rows, _ := tx.Query(fmt.Sprintf("CALL get_mail('%d','%d')", 0, 10))
	defer rows.Close()
	for rows.Next() {
		var gid, rolegid uint64
		var blobData string
		rows.Scan(&gid, &rolegid, &blobData)
		t.Log(gid)
		t.Log(rolegid)

		var dataS mailData
		checkErr(json.Unmarshal([]byte(blobData), &dataS))
		t.Logf("Item = %v, Count = %v, Desc = %v", dataS.Item, dataS.Count, dataS.Desc)
	}
}