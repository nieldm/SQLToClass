package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/etgryphon/stringUp"
	"github.com/ogier/pflag"
	"github.com/xwb1989/sqlparser"
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
			rawSql = strings.Replace(rawSql, "index", "", -1)
			rawSql = strings.Replace(rawSql, ";", "", -1)

			result, err := sqlparser.Parse(rawSql)

			if err == nil {
				extractClass(result)
			} else {
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
	tableName := stringUp.CamelCase(string(node.Name))
	tableText := fmt.Sprintf("TableName: %s", tableName)
	fmt.Println(tableText)
	for _, col := range node.ColumnDefinitions {
		colCCName := stringUp.CamelCase(col.ColName)
		columnText := fmt.Sprintf("   ColumnName: %s ColType: %s -> %s", col.ColName, col.ColType, colCCName)
		fmt.Println(columnText)
	}
}

func init() {
	pflag.StringVarP(&sql, "sql", "s", "", "SQL file to Parse")
}
