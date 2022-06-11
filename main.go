package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path"
	"sort"
)

// model

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

func (a *Article) readDetail(root string) {
	file, err := os.Open(path.Join(root, fmt.Sprintf("%d.md", a.AID)))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	a.Name = scanner.Text()[2:]
}

func (c *Category) sortArticleByName() {
	sort.Slice(c.Article, func(i, j int) bool {
		return c.Article[i].Name > c.Article[j].Name
	})
}

// converter

func convertTreeToMarkdown(cat *Category, deep int, buff *bytes.Buffer) {
	space := ""
	for i := 0; i < deep; i++ {
		space = space + "  "
	}
	buff.WriteString(fmt.Sprintf("%s- %s\n", space, cat.Name))
	for _, article := range cat.Article {
		buff.WriteString(fmt.Sprintf("%s  - [%s](./docs/%d.md)\n", space, article.Name, article.AID))
	}
	for _, category := range cat.SubCategory {
		convertTreeToMarkdown(category, deep+1, buff)
	}
}

func convertDataToTree(root *Category, input []*Category) []*Category {
	var otherNode []*Category
	for _, category := range input {
		if category.PID == root.UUID {
			root.SubCategory = append(root.SubCategory, category)
		} else {
			otherNode = append(otherNode, category)
		}
	}
	for _, category := range root.SubCategory {
		otherNode = convertDataToTree(category, otherNode)
	}
	return otherNode
}

func bindArticleToCategory(cat []*Category, art []*Article) *Category {
	catMap := map[uint64]*Category{}
	var root *Category
	for _, category := range cat {
		catMap[category.UUID] = category
		if category.PID == 0 {
			root = category
		}
	}
	for _, article := range art {
		article.readDetail(path.Join(lib, "docs"))
		if c, ok := catMap[article.RID]; ok {
			c.Article = append(c.Article, article)
		}
	}

	for _, category := range cat {
		category.sortArticleByName()
	}
	return root
}

// dao
func loadDatasource() (cat []*Category, art []*Article) {
	if lib == "" {
		log.Fatalf("path to MWebLibrary is empty")
	}
	sqlPath := path.Join(lib, "mainlib.db")
	db, err := sql.Open("sqlite", sqlPath)
	if err != nil {
		log.Fatalf("open sqlite3 failed: %s", err)
	}
	defer db.Close()

	cat, err = loadCategories(db)
	if err != nil {
		log.Fatalf("load categories failed: %s", err)
	}
	art, err = loadArticles(db)
	if err != nil {
		log.Fatalf("load articles failed: %s", err)
	}
	return cat, art
}

func loadCategories(db *sql.DB) ([]*Category, error) {
	row, err := db.Query("select pid, uuid, name from cat order by sort")
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

func loadArticles(db *sql.DB) ([]*Article, error) {
	row, err := db.Query("select rid, aid from cat_article;")
	// 	row, err := db.Query("select cat_article.rid, cat_article.aid from cat_article left join article on cat_article.aid = article.uuid order by  article.sort ;")
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

var (
	lib    string
	target string
	mode   string
	help   bool
)

func init() {
	home, _ := os.UserHomeDir()
	pwd, _ := os.Getwd()
	libDefaultPath := home + "/Library/Containers/com.coderforart.MWeb3/Data/Library/Application Support/MWebLibrary"

	flag.StringVar(&lib, "path", libDefaultPath, "path to MWebLibrary")
	flag.StringVar(&target, "target", pwd, "export README.md directory")
	flag.StringVar(&mode, "mode", "debug", "'save': save file, 'debug': print only")
	flag.BoolVar(&help, "help", false, "show usage")

	flag.Parse()
}

func main() {

	if help {
		flag.Usage()
		return
	}

	var buffer bytes.Buffer
	cat, art := loadDatasource()
	root := bindArticleToCategory(cat, art)

	buffer.WriteString("# NoteBook\n\n")
	convertDataToTree(root, cat)
	convertTreeToMarkdown(root, 0, &buffer)
	switch mode {
	case "save":
		err := ioutil.WriteFile(path.Join(target, "README.md"), buffer.Bytes(), 0644)
		if err != nil {
			log.Fatalf("Write file fail: %v", err)
			return
		}
	case "debug":
		fmt.Printf("%s", buffer.String())
	default:
		log.Fatalf("Unknown mode: %s", mode)
	}

}
