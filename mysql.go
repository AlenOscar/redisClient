/**
 * Created by Alen on 2019-05-20 14:02
 */

package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"time"
)

var Db *DbWorker

type DbWorker struct {
	db *sql.DB
}

const (
	DbUsername = "root"      // 用户名
	DbPassword = "123456"    // 密码
	DbHost     = "127.0.0.1" // 数据库主机
	DbPort     = 3308        // 端口
	DbName     = "testdb"    // 数据库名
)

func init() {
	/* DSN数据源名称
	[username[:password]@][protocol[(address)]]/dbname[?param1=value1¶mN=valueN]
	user@unix(/path/to/socket)/dbname
	user:password@tcp(localhost:5555)/dbname?charset=utf8&autocommit=true
	user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname?charset=utf8mb4,utf8
	user:password@/dbname
	无数据库: user:password@/
	*/
	dataSourceName := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8", DbUsername, DbPassword, DbHost, DbPort, DbName)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		fmt.Printf("Db open error: %v", err)
		return
	}

	Db = &DbWorker{
		db: db,
	}
}

func (w *DbWorker) Insert(gid uint64, data string) {
	// 方式1 insert
	// strconv,int转string:strconv.Itoa(i)
	start := time.Now()
	for i := 1001; i <= 1100; i++ {
		// 每次循环内部都会去连接池获取一个新的连接，效率低下
		w.db.Exec("INSERT INTO role(uid,username,age) values(?,?,?)", i, "user"+strconv.Itoa(i), i-1000)
	}
	end := time.Now()
	fmt.Println("方式1 insert total time:", end.Sub(start).Seconds())

	// 方式2 insert
	start = time.Now()
	for i := 1101; i <= 1200; i++ {
		// Prepare函数每次循环内部都会去连接池获取一个新的连接，效率低下
		stm, _ := w.db.Prepare("INSERT INTO role(uid,username,age) values(?,?,?)")
		stm.Exec(i, "user"+strconv.Itoa(i), i-1000)
		stm.Close()
	}
	end = time.Now()
	fmt.Println("方式2 insert total time:", end.Sub(start).Seconds())

	// 方式3 insert
	start = time.Now()
	stm, _ := w.db.Prepare("INSERT INTO role(uid,username,age) values(?,?,?)")
	for i := 1201; i <= 1300; i++ {
		// Exec内部并没有去获取连接，为什么效率还是低呢？
		stm.Exec(i, "user"+strconv.Itoa(i), i-1000)
	}
	stm.Close()
	end = time.Now()
	fmt.Println("方式3 insert total time:", end.Sub(start).Seconds())

	// 方式4 insert
	start = time.Now()
	// Begin函数内部会去获取连接
	tx, _ := w.db.Begin()
	for i := 1301; i <= 1400; i++ {
		// 每次循环用的都是tx内部的连接，没有新建连接，效率高
		tx.Exec("INSERT INTO role(uid,username,age) values(?,?,?)", i, "user"+strconv.Itoa(i), i-1000)
	}
	// 最后释放tx内部的连接
	tx.Commit()

	end = time.Now()
	fmt.Println("方式4 insert total time:", end.Sub(start).Seconds())

	// 方式5 insert
	start = time.Now()
	for i := 1401; i <= 1500; i++ {
		// Begin函数每次循环内部都会去连接池获取一个新的连接，效率低下
		tx, _ := w.db.Begin()
		tx.Exec("INSERT INTO role(uid,username,age) values(?,?,?)", i, "user"+strconv.Itoa(i), i-1000)
		// Commit执行后连接也释放了
		tx.Commit()
	}
	end = time.Now()
	fmt.Println("方式5 insert total time:", end.Sub(start).Seconds())
}

func (w *DbWorker) Delete(gid uint64) {
	// 方式1 delete
	start := time.Now()
	for i := 1001; i <= 1100; i++ {
		w.db.Exec("DELETE FROM role WHERE uid=?", i)
	}
	end := time.Now()
	fmt.Println("方式1 delete total time:", end.Sub(start).Seconds())

	// 方式2 delete
	start = time.Now()
	for i := 1101; i <= 1200; i++ {
		stm, _ := w.db.Prepare("DELETE FROM role WHERE uid=?")
		stm.Exec(i)
		stm.Close()
	}
	end = time.Now()
	fmt.Println("方式2 delete total time:", end.Sub(start).Seconds())

	// 方式3 delete
	start = time.Now()
	stm, _ := w.db.Prepare("DELETE FROM role WHERE uid=?")
	for i := 1201; i <= 1300; i++ {
		stm.Exec(i)
	}
	stm.Close()
	end = time.Now()
	fmt.Println("方式3 delete total time:", end.Sub(start).Seconds())

	// 方式4 delete
	start = time.Now()
	tx, _ := w.db.Begin()
	for i := 1301; i <= 1400; i++ {
		tx.Exec("DELETE FROM role WHERE uid=?", i)
	}
	tx.Commit()

	end = time.Now()
	fmt.Println("方式4 delete total time:", end.Sub(start).Seconds())

	// 方式5 delete
	start = time.Now()
	for i := 1401; i <= 1500; i++ {
		tx, _ := w.db.Begin()
		tx.Exec("DELETE FROM role WHERE uid=?", i)
		tx.Commit()
	}
	end = time.Now()
	fmt.Println("方式5 delete total time:", end.Sub(start).Seconds())
}

func (w *DbWorker) Update(gid uint64, data string) {
	// 方式1 update
	start := time.Now()
	for i := 1001; i <= 1100; i++ {
		w.db.Exec("UPdate role set age=? where uid=? ", i, i)
	}
	end := time.Now()
	fmt.Println("方式1 update total time:", end.Sub(start).Seconds())

	// 方式2 update
	start = time.Now()
	for i := 1101; i <= 1200; i++ {
		stm, _ := w.db.Prepare("UPdate role set age=? where uid=? ")
		stm.Exec(i, i)
		stm.Close()
	}
	end = time.Now()
	fmt.Println("方式2 update total time:", end.Sub(start).Seconds())

	// 方式3 update
	start = time.Now()
	stm, _ := w.db.Prepare("UPdate role set age=? where uid=?")
	for i := 1201; i <= 1300; i++ {
		stm.Exec(i, i)
	}
	stm.Close()
	end = time.Now()
	fmt.Println("方式3 update total time:", end.Sub(start).Seconds())

	// 方式4 update
	start = time.Now()
	tx, _ := w.db.Begin()
	for i := 1301; i <= 1400; i++ {
		tx.Exec("UPdate role set age=? where uid=?", i, i)
	}
	tx.Commit()

	end = time.Now()
	fmt.Println("方式4 update total time:", end.Sub(start).Seconds())

	// 方式5 update
	start = time.Now()
	for i := 1401; i <= 1500; i++ {
		tx, _ := w.db.Begin()
		tx.Exec("UPdate role set age=? where uid=?", i, i)
		tx.Commit()
	}
	end = time.Now()
	fmt.Println("方式5 update total time:", end.Sub(start).Seconds())
}

func (w *DbWorker) Query(gid uint64) {
	// 方式1 query
	start := time.Now()
	rows, _ := w.db.Query("SELECT uid,username FROM role")
	defer rows.Close()
	for rows.Next() {
		var name string
		var id int
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("name:%s ,id:is %d\n", name, id)
	}
	end := time.Now()
	fmt.Println("方式1 query total time:", end.Sub(start).Seconds())

	// 方式2 query
	start = time.Now()
	stm, _ := w.db.Prepare("SELECT uid,username FROM role")
	defer stm.Close()
	rows, _ = stm.Query()
	defer rows.Close()
	for rows.Next() {
		var name string
		var id int
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("name:%s ,id:is %d\n", name, id)
	}
	end = time.Now()
	fmt.Println("方式2 query total time:", end.Sub(start).Seconds())

	// 方式3 query
	start = time.Now()
	tx, _ := w.db.Begin()
	defer tx.Commit()
	rows, _ = tx.Query("SELECT uid,username FROM role")
	defer rows.Close()
	for rows.Next() {
		var name string
		var id int
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("name:%s ,id:is %d\n", name, id)
	}
	end = time.Now()
	fmt.Println("方式3 query total time:", end.Sub(start).Seconds())
}
