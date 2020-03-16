package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	simpleUtil "github.com/liserjrqlxue/simple-util"
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
		"",
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
)

var keyMap = map[string]string{
	"自动化+人工证据":      "MUTATION_TYPE",
	"自动化+人工致病等级":    "LITERATURE",
	"PP_References": "DISEASE_PHENOTYPE",
}

func main() {
	flag.Parse()
	if *input == "" || *output == "" || *db == "" || *sheetName == "" {
		flag.Usage()
		os.Exit(1)
	}

	keys := strings.Split(*key, ":")

	// load excel
	var allDb = make(map[string]map[string]string)
	xlsxFh, err := excelize.OpenFile(*db)
	simpleUtil.CheckErr(err)
	rows, err := xlsxFh.GetRows(*sheetName)
	simpleUtil.CheckErr(err)
	for i, row := range rows {
		if i == 0 {
			keys = row
		} else {
			var item = make(map[string]string)
			for j, cell := range row {
				item[keys[j]] = cell
			}
			var keyValues []string
			for _, k := range keys {
				keyValues = append(keyValues, item[k])
			}
			mainKey := strings.Join(keyValues, ":")
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

			var titleHash = make(map[string]bool)
			for _, k := range title {
				titleHash[k] = true
			}
			for _, v := range keyMap {
				if !titleHash[v] {
					title = append(title, v)
				}
			}
			if title == nil {
				log.Fatal("title == nil")
			}
			_, err = fmt.Fprintln(outputFh, strings.Join(title, "\t"))
			simpleUtil.CheckErr(err)
		} else {
			var item = make(map[string]string)
			for j, k := range array {
				item[title[j]] = k
			}

			// annotation
			key := item["Transcript"] + ":" + item["cHGVS"]
			for k, v := range keyMap {
				item[v] = allDb[key][k]
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
