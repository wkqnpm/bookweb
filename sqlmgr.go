/*
 * @Descripttion:
 * @version:
 * @Author: wkq
 * @Date: 2022-06-12 10:10:16
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2022-07-02 19:11:09
 */
package sqlmgr

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

var db *sqlx.DB

type Book struct {
	ID         string `json:"id" db:"id"`
	BookName   string `json:"bookname" db:"bookname"`
	Author     string `json:"author" db:"author"`
	CreateTime int64  `json:"create_time" db:"createtime"`
}

func InitDB(username, pwd, url, dbname string) (err error) {

	dsn := username + ":" + pwd + "@tcp(" + url + ")/" + dbname
	// dsn := "root:root12345@tcp(127.0.0.1:3306)/gostudy"
	db, err = sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Printf("open database falid %v\n", err)
		return
	}
	err = db.Ping() //校验
	if err != nil {
		return
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(20)
	return
}
func CreateBookTable() {
	var schema = `CREATE TABLE books (
		id VARCHAR(36) NOT NULL,
		bookname VARCHAR(50) DEFAULT '',
		author VARCHAR(50) DEFAULT '',
		createtime VARCHAR(50) DEFAULT '',
		PRIMARY KEY(id)
	);`
	res, err := db.Exec(schema)
	if err != nil {
		fmt.Println("create table faild", err)
		return
	}
	fmt.Println(res)
}
func QueryAllData(pagenum, pagesize int) (books []Book, total int, err error) {
	pageNum := (pagenum - 1) * pagesize
	sqlstr := "select * from books order by createtime limit ?,? ;"
	err = db.Select(&books, sqlstr, pageNum, pagesize)
	sqlstr1 := "select count(*) from books;"
	totalRow, err1 := db.Query(sqlstr1)
	if err1 != nil {
		fmt.Println("GetKnowledgePointListTotal error", err)
		return
	}
	total = 0
	for totalRow.Next() {
		err = totalRow.Scan(
			&total,
		)
		if err != nil {
			fmt.Println("GetKnowledgePointListTotal error", err)
			continue
		}
	}
	return
}
func DelData(id string) (err error) {
	sql := "delete from books where id=?"
	_, err = db.Exec(sql, id)
	if err != nil {
		fmt.Printf("update faild: %v", err)
		return
	}
	return
}
func UpdateData(id string, bookname string) error {
	sql := "update books set bookname=? where id=?"
	res, err := db.Exec(sql, bookname, id)
	if err != nil {
		fmt.Printf("update faild: %v", err)
		return err
	}
	_, err = res.RowsAffected() //新插入数据的ID
	if err != nil {
		fmt.Printf("update faild and no effect: %v", err)
		return err
	}
	return err
}
func AddData(bookname, author string) (err error) {
	sqlstr := "insert into books(id,bookname,author,createtime) values(?,?,?,?)"
	uid, _ := uuid.NewV4()
	id := uid.String()
	nowsec := time.Now().Unix()
	res := db.MustExec(sqlstr, id, bookname, author, nowsec)
	_, err = res.RowsAffected()
	return
}
