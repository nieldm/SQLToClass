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

	testTemplate()
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
	prefix := "M"
	tableName := prefix + strings.Title(stringUp.CamelCase(string(node.Name)))

	fileName := fmt.Sprintf("/tmp/%s.swift", tableName)

	tableText := fmt.Sprintf("class %s : Object {", tableName)
	fmt.Println(tableText)
	tableText = "import RealmSwift\n\n" + tableText
	for _, col := range node.ColumnDefinitions {
		colCCName := stringUp.CamelCase(col.ColName)
		var colType string = col.ColType
		var colValue string = "\"\""
		switch col.ColType {
		case "integer":
			colType = "Int"
			colValue = "0"
		case "varchar", "text":
			colType = "String"
		case "datetime":
			colType = "Date"
			colValue = "Date()"
		case "float":
			colType = "Float"
			colValue = "0.0"
		case "tinyint(1)":
			colType = "Bool"
			colValue = "false"
		}
		columnText := fmt.Sprintf("\tdynamic var %s: %s = %s", colCCName, colType, colValue)
		tableText += "\n" + columnText
		fmt.Println(columnText)
	}
	tableText += "\n" + "}"
	fmt.Println("}")
	d1 := []byte(tableText)
	err := ioutil.WriteFile(fileName, d1, 0644)
	check(err)
}

func init() {
	pflag.StringVarP(&sql, "sql", "s", "", "SQL file to Parse")
}

func testTemplate() {
	fmt.Printf("%s\n", templates.Hello("Foo"))
	fmt.Printf("%s\n", templates.Hello("Bar"))
	fields := []string{"Kate", "Go", "John", "Brad"}
	fieldTypes := []string{"Int", "String", "String", "Int"}
	var buf bytes.Buffer
	templates.WriteBuild(&buf, "test", fields, fieldTypes)
	fmt.Printf("%s", buf.Bytes())
}
