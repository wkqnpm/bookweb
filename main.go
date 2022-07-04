/*
 * @Descripttion:
 * @version:
 * @Author: wkq
 * @Date: 2022-06-12 10:05:33
 * @LastEditors: Please set LastEditors
 * @LastEditTime: 2022-07-02 19:11:59
 */
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	sqlmgr "Pro/book/sqlmgr"
)

type PageConf struct {
	pageNum  int `json:"pageNum"`
	pageSize int `json:"pageSize"`
}

type Ret struct {
	Code  int           `json:"code"`
	Msg   string        `json:"msg"`
	Data  []sqlmgr.Book `json:"data"`
	Total int           `json:"total"`
}

type Record struct {
	Id       string `json:"id"`
	Author   string `json:"author"`
	BookName string `json:"bookname"`
}

func cros(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// 允许访问所有域，可以换成具体url，注意仅具体url才能带cookie信息
		w.Header().Set("Access-Control-Allow-Origin", "*")
		//header的类型
		//w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		//设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		//允许请求方法
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		//返回数据格式是json
		//w.Header().Set("content-type", "application/json;charset=UTF-8")
		if req.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, req)
	}
}

func ResData(err error, W http.ResponseWriter, msg string) {
	var ret Ret
	if err != nil {
		fmt.Println(msg)
		ret.Code = 500
		ret.Msg = msg
		ret.Data = nil
		v, _ := json.Marshal(ret)
		W.Write(v)
		return
	}
	ret.Code = 200
	ret.Msg = "ok"
	v, _ := json.Marshal(ret)
	W.Write(v)
}

func handleList(W http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var pageConf PageConf
	//解析处理查询参数
	pageConf.pageNum, _ = strconv.Atoi(queryParams.Get("pageNum"))
	pageConf.pageSize, _ = strconv.Atoi(queryParams.Get("pageSize"))
	//查询数据
	res, total, err := sqlmgr.QueryAllData(pageConf.pageNum, pageConf.pageSize)
	if err != nil {
		fmt.Println(err)
		return
	}
	var temp Ret
	temp.Code = 200
	temp.Msg = "success"
	temp.Total = total
	for _, v := range res {
		temp.Data = append(temp.Data, v)
	}
	v, _ := json.Marshal(temp)
	W.Write(v)
}
func handleAdd(W http.ResponseWriter, r *http.Request) {
	params, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	var ad Record
	json.Unmarshal(params, &ad)
	err = sqlmgr.AddData(ad.Author, ad.BookName)
	ResData(err, W, "add data error")
}
func handleEdit(W http.ResponseWriter, r *http.Request) {
	params, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	var ed Record
	json.Unmarshal(params, &ed)
	err = sqlmgr.UpdateData(ed.Id, ed.BookName)
	ResData(err, W, "edit data error")
}
func handleDel(W http.ResponseWriter, r *http.Request) {
	params, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
	}
	var ed Record
	json.Unmarshal(params, &ed)
	err = sqlmgr.DelData(ed.Id)
	ResData(err, W, "del data error")
}
func main() {
	err := sqlmgr.InitDB("root", "123456", "127.0.0.1:3306", "books")
	if err != nil {
		fmt.Println("connect db faild\n", err)
		return
	}
	fmt.Println("database connect success")
	// sqlmgr.CreateBookTable()
	http.HandleFunc("/list", cros(handleList))
	http.HandleFunc("/add", cros(handleAdd))
	http.HandleFunc("/edit", cros(handleEdit))
	http.HandleFunc("/del", cros(handleDel))

	err = http.ListenAndServe("127.0.0.1:8088", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
