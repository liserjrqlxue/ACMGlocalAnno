package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/textUtil"
)

// os
var (
	ex, _  = os.Executable()
	exPath = filepath.Dir(ex)
)

var (
	input = flag.String(
		"input",
		"",
		"input file",
	)
	output = flag.String(
		"output",
		"",
		"output file",
	)
	db = flag.String(
		"db",
		path.Join(exPath, "ACMG59.db.xlsx"),
		"acmg local db file",
	)
	sheetName = flag.String(
		"sheetName",
		"Sheet1",
		"sheet name of db.xlsx",
	)
	key = flag.String(
		"key",
		"Transcript:cHGVS",
		"main key for hit",
	)
	titleMap = flag.String(
		"titleMap",
		path.Join(exPath, "title.txt"),
		"title map",
	)
)

func main() {
	flag.Parse()
	if *input == "" || *output == "" || *db == "" || *sheetName == "" {
		flag.Usage()
		os.Exit(1)
	}

	var titleHash, err = textUtil.File2Map(*titleMap, "\t", true)
	simpleUtil.CheckErr(err)

	var keys = strings.Split(*key, ":")
	var allDb = make(map[string]map[string]string)
	xlsxFh, err := excelize.OpenFile(*db)
	simpleUtil.CheckErr(err)
	rows, err := xlsxFh.GetRows(*sheetName)
	simpleUtil.CheckErr(err)
	var firstLine []string
	for i, row := range rows {
		if i == 0 {
			firstLine = row
		} else {
			var item = make(map[string]string)
			for j, cell := range row {
				if firstLine != nil {
					item[firstLine[j]] = cell
				}
			}
			var keyValues []string
			for _, k := range keys {
				keyValues = append(keyValues, item[k])
			}
			var mainKey = strings.Join(keyValues, ":")
			allDb[mainKey] = item
		}
	}

	// load input
	file, err := os.Open(*input)
	simpleUtil.CheckErr(err)
	defer simpleUtil.DeferClose(file)

	// create output
	outputFh, err := os.Create(*output)
	simpleUtil.CheckErr(err)
	defer simpleUtil.DeferClose(outputFh)

	scanner := bufio.NewScanner(file)
	sep := "\t"
	var title []string
	var header = true
	for scanner.Scan() {
		line := scanner.Text()
		array := strings.Split(line, sep)
		if header {
			header = false
			title = array

			var inTitle = make(map[string]bool)
			for _, k := range title {
				inTitle[k] = true
			}
			for k := range titleHash {
				if !inTitle[k] {
					title = append(title, k)
				}
			}
			if title == nil {
				log.Fatal("title == nil")
			}
			_, err = fmt.Fprintln(outputFh, strings.Join(title, "\t"))
			simpleUtil.CheckErr(err)
		} else {
			var item = make(map[string]string)
			for i, v := range array {
				if title != nil {
					item[title[i]] = v
				}
			}

			// annotation
			key := item["Transcript"] + ":" + item["cHGVS"]
			for k, v := range titleHash {
				item[k] = allDb[key][v]
			}
			var row []string
			for _, k := range title {
				row = append(row, item[k])
			}
			_, err = fmt.Fprintln(outputFh, strings.Join(row, "\t"))
			simpleUtil.CheckErr(err)
		}
	}
	simpleUtil.CheckErr(scanner.Err())
}
