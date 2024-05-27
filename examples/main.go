package main

import (
	"fmt"

	"github.com/taojy93/gofiledb"
)

func main() {

	// 创建数据库实例
	db, err := gofiledb.NewDBClient("example_db")
	if err != nil {
		fmt.Println("Error creating database:", err)
		return
	}

	// 创建表
	if err := db.CreateTable("example_table"); err != nil {
		fmt.Println("Error creating table:", err)
		return
	}

	// 添加记录
	if err := db.AddRecord("example_table", gofiledb.Record{Data: "First Record"}); err != nil {
		fmt.Println("Error adding record:", err)
		return
	}

	if err := db.AddRecord("example_table", gofiledb.Record{Data: "Second Record"}); err != nil {
		fmt.Println("Error adding record:", err)
		return
	}

	// 获取记录
	record, err := db.GetRecord("example_table", 1)
	if err != nil {
		fmt.Println("Error getting record:", err)
		return
	}
	fmt.Printf("Got record: ID: %d, Data: %v\n", record.ID, record.Data)

	// 更新记录
	if err := db.UpdateRecord("example_table", 1, "Updated First Record"); err != nil {
		fmt.Println("Error updating record:", err)
		return
	}

	// 获取更新后的记录
	record, err = db.GetRecord("example_table", 1)
	if err != nil {
		fmt.Println("Error getting record:", err)
		return
	}
	fmt.Printf("Got updated record: ID: %d, Data: %v\n", record.ID, record.Data)

	// 删除记录
	if err := db.DeleteRecord("example_table", 1); err != nil {
		fmt.Println("Error deleting record:", err)
		return
	}

	// 尝试获取已删除的记录
	record, err = db.GetRecord("example_table", 1)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Got record: ID: %d, Data: %v\n", record.ID, record.Data)
	}
}
