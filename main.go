package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
	_, database := simpleUtil.Sheet2MapArray(*db, *sheetName)
	for _, item := range database {
		var keyValues []string
		for _, k := range keys {
			keyValues = append(keyValues, item[k])
		}
		mainKey := strings.Join(keyValues, ":")
		allDb[mainKey] = item
	}

	// load input
	anno, title := simpleUtil.File2MapArray(*input, "\t", nil)

	var titleHash = make(map[string]bool)
	for _, k := range title {
		titleHash[k] = true
	}
	for _, v := range keyMap {
		if !titleHash[v] {
			title = append(title, v)
		}
	}

	// create output
	outputFh, err := os.Create(*output)
	simpleUtil.CheckErr(err)
	defer simpleUtil.DeferClose(outputFh)

	_, err = fmt.Fprintln(outputFh, strings.Join(title, "\t"))
	simpleUtil.CheckErr(err)

	// annotation
	for _, item := range anno {
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
