package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/liserjrqlxue/anno2xlsx/v2/anno"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/liserjrqlxue/goUtil/textUtil"
	"github.com/xuri/excelize/v2"
)

// os
var (
	ex, _  = os.Executable()
	exPath = filepath.Dir(ex)
	dbPath = filepath.Join(exPath, "db")
)

// \n -> <br/>
var isLF = regexp.MustCompile(`\n`)

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
		path.Join(dbPath, "ACMGSF.xlsx"),
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

	var titleHash, titleOrder, err = textUtil.File2MapOrder(*titleMap, "\t", true)
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

	var (
		scanner = bufio.NewScanner(file)
		sep     = "\t"
		skip    = regexp.MustCompile(`^##`)
		title   []string
		header  = true
	)

	for scanner.Scan() {
		line := scanner.Text()
		if skip.MatchString(line) {
			continue
		}
		array := strings.Split(line, sep)
		if header {
			header = false
			title = array

			var inTitle = make(map[string]bool)
			for _, k := range title {
				inTitle[k] = true
			}
			for _, k := range titleOrder {
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
			var target, ok = anno.GetFromMultiKeys(
				allDb,
				anno.GetKeys(item["Transcript"], item["cHGVS"]),
			)
			if ok {
				for k, v := range titleHash {
					item[k] = target[v]
				}
			}
			var row []string
			for _, k := range title {
				row = append(row, isLF.ReplaceAllString(item[k], "<br/>"))
			}
			_, err = fmt.Fprintln(outputFh, strings.Join(row, "\t"))
			simpleUtil.CheckErr(err)
		}
	}
	simpleUtil.CheckErr(scanner.Err())
}
