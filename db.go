package main

import (
	"database/sql"
	"fmt"
)

//写了一个工厂模式的数据库选择

type Do interface {
	GetDB() interface{}
}
type DB struct {
	do Do
}

type Mysqld struct {
	path string
}

type SqlServerd struct {
	path string
}

//获取连接
func (d Mysqld) GetDB() interface{} {
	db, err := sql.Open("mysql", d.path)
	if err != nil {
		panic(fmt.Sprint("mysql connect err %s", err.Error()))
	}
	return db
}
func (d SqlServerd) GetDb() interface{} {

	return nil
}
