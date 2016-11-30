package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/etgryphon/stringUp"
	"github.com/ogier/pflag"
	"github.com/xwb1989/sqlparser"

	"./templates"
)

//flags
var (
	sql string
)

func check(e error) {
	if e != nil {
		log.Print(e)
	}
}

func main() {
	pflag.Parse()

	fmt.Println("SQL to Realm.swift Parser")

	files := strings.Split(sql, ",")
	fmt.Printf("Searching files(s): %s\n", files)

	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println("Error 1")
			log.Fatal(err)
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			rawSql := strings.Replace(scanner.Text(), "\"", "", -1)
			rawSql = strings.Replace(rawSql, "autoincrement", "", -1)
			rawSql = strings.Replace(rawSql, "index", "sortIndex", -1)
			rawSql = strings.Replace(rawSql, ";", "", -1)

			result, err := sqlparser.Parse(rawSql)

			if err == nil {
				extractClass(result)
			} else {
				fmt.Println(rawSql)
				check(err)
				continue
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
			fmt.Println("Error 3")
		}

	}

}

func extractClass(sqlNode sqlparser.SQLNode) {
	node, ok := sqlNode.(*sqlparser.CreateTable)

	if !ok {
		return
	}

	tableName := string(node.Name)

	var fields []string
	var fieldTypes []string

	for _, col := range node.ColumnDefinitions {
		fields = append(fields, col.ColName)
		fieldTypes = append(fieldTypes, col.ColType)
	}
	buildTemplate(tableName, fields, fieldTypes)
}

func init() {
	pflag.StringVarP(&sql, "sql", "s", "", "SQL file to Parse")
}

func buildTemplate(tableName string, fields []string, fieldTypes []string) {
	var buf bytes.Buffer
	templates.WriteBuild(&buf, tableName, fields, fieldTypes)
	fmt.Printf("%s", buf.Bytes())
	writeToFile(tableName, buf.Bytes())
}

func writeToFile(tableName string, d1 []byte) {
	fileNameRaw := strings.Title(stringUp.CamelCase(tableName))
	fileName := fmt.Sprintf("/tmp/%s.swift", fileNameRaw)
	err := ioutil.WriteFile(fileName, d1, 0644)
	check(err)
}

func testTemplate() {
	fmt.Printf("%s\n", templates.Hello("Foo"))
	fmt.Printf("%s\n", templates.Hello("Bar"))
}
