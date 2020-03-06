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

var appendColumns = []string{"自动化+人工证据", "自动化+人工致病等级", "PP_References"}

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

	// create output
	outputFh, err := os.Create(*output)
	simpleUtil.CheckErr(err)
	defer simpleUtil.DeferClose(outputFh)

	title = append(title, appendColumns...)
	_, err = fmt.Fprintln(outputFh, strings.Join(title, "\t"))
	simpleUtil.CheckErr(err)

	// annotation
	for _, item := range anno {
		key := item["Transcript"] + ":" + item["cHGVS"]
		for _, k := range appendColumns {
			item[k] = allDb[key][k]
		}
		var row []string
		for _, k := range title {
			row = append(row, item[k])
		}
		_, err = fmt.Fprintln(outputFh, strings.Join(row, "\t"))
		simpleUtil.CheckErr(err)
	}
}
