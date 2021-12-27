package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

type Category struct {
	PID         uint64
	UUID        uint64
	Name        string
	SubCategory []*Category
	Article     []*Article
}

type Article struct {
	RID   uint64
	AID   uint64
	Name  string
	Media []string
}

func (d *Database) Open(path string) error {
	db, err := sql.Open("sqlite3", path)
	d.db = db
	return err
}

func (d *Database) Categories() ([]*Category, error) {
	row, err := d.db.Query("select pid, uuid, name from cat")
	if err != nil {
		return nil, err
	}
	var temp []*Category
	for row.Next() {
		var (
			pid  uint64
			uuid uint64
			name string
		)
		if err := row.Scan(&pid, &uuid, &name); err == nil {
			temp = append(temp, &Category{
				PID:  pid,
				UUID: uuid,
				Name: name,
			})
		}
	}
	return temp, nil
}

func (d *Database) Article() ([]*Article, error) {
	row, err := d.db.Query("select rid, aid from cat_article")
	if err != nil {
		return nil, err
	}
	var temp []*Article
	for row.Next() {
		var (
			rid uint64
			aid uint64
		)
		if err := row.Scan(&rid, &aid); err == nil {
			temp = append(temp, &Article{
				RID: rid,
				AID: aid,
			})
		}
	}
	return temp, nil
}

func makeCategoryTree(root *Category, input []*Category) []*Category {
	var otherNode []*Category
	for _, category := range input {
		if category.PID == root.UUID {
			root.SubCategory = append(root.SubCategory, category)
		} else {
			otherNode = append(otherNode, category)
		}
	}
	for _, category := range root.SubCategory {
		otherNode = makeCategoryTree(category, otherNode)
	}
	return otherNode
}

func main() {
	path := "/Users/tbxark/Desktop/Repos/Notebook/"
	db := &Database{}
	db.Open(path + "mainlib.db")
	cat, _ := db.Categories()
	var root *Category
	for _, category := range cat {
		if category.PID == 0 {
			root = category
			break
		}
	}
	makeCategoryTree(root, cat)
	fmt.Print(root)
	//art, _ := db.Article()
	//fmt.Printf("Cat: %v, Art: %v", cat, art)
}
