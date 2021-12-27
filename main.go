package main

import (
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

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

func (a *Article) update(root string) {
	b, err := ioutil.ReadFile(path.Join(root, fmt.Sprintf("%d.md", a.AID)))
	if err != nil {
		log.Fatalf("Read %d file error", a.AID)
	}
	f := string(b)
	line := strings.Split(f, "\n")
	a.Name = line[0][2:]
	//reg := regexp.MustCompile("!\\[[^\\]]*\\]\\((.*?)\\s*(\"(?:.*[^\"])\")?\\s*\\)")
	//imgs := reg.FindAllString(f, -1)
	//for _, img := range imgs {
	//	fmt.Printf("find image: %s\n", img)
	//}
}

func tree(cat *Category, deep int, buff *bytes.Buffer) {
	space := ""
	for i := 0; i < deep; i++ {
		space = space + "  "
	}
	buff.WriteString(fmt.Sprintf("%s- %s\n", space, cat.Name))
	for _, article := range cat.Article {
		buff.WriteString(fmt.Sprintf("%s  - [%s](./docs/%d.md)\n", space, article.Name, article.AID))
	}
	for _, category := range cat.SubCategory {
		tree(category, deep+1, buff)
	}
}

func categories(db *sql.DB) ([]*Category, error) {
	row, err := db.Query("select pid, uuid, name from cat")
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

func article(db *sql.DB) ([]*Article, error) {
	row, err := db.Query("select rid, aid from cat_article")
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
	lib := flag.String("path", "", "Path to MWebLibrary")
	flag.Parse()

	if *lib == "" {
		log.Fatalf("You must set MWebLibrary path")
	}

	sqlPath := path.Join(*lib, "mainlib.db")
	log.Printf("Open lib: %s", sqlPath)
	db, dErr := sql.Open("sqlite3", sqlPath)
	if dErr != nil {
		log.Fatalf("Open database  fail: %v", dErr)
	}

	cat, cErr := categories(db)
	if cErr != nil {
		log.Fatalf("Read categories fail: %v", cErr)
	}

	art, aErr := article(db)
	if aErr != nil {
		log.Fatalf("Read article fail: %v", dErr)
	}

	catMap := map[uint64]*Category{}
	var root *Category
	var buffer bytes.Buffer

	for _, category := range cat {
		catMap[category.UUID] = category
		if category.PID == 0 {
			root = category
		}
	}
	for _, article := range art {
		article.update(path.Join(*lib, "docs"))
		if c, ok := catMap[article.RID]; ok {
			c.Article = append(c.Article, article)
		}
	}
	makeCategoryTree(root, cat)
	buffer.WriteString("# NoteBook\n\n")
	tree(root, 0, &buffer)
	ioutil.WriteFile(path.Join(*lib, "README.md"), buffer.Bytes(), 0644)
}
