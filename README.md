# ACMGlocalAnno
ACMG定点筛查

## 注释
在原来注释信息的基础上添加`自动化+人工证据`，`自动化+人工致病等级`，`PP_References`三列内容信息

## ACMG local db
1. directly load excel
2. convert excel to json data
3. main key of item is 'transcript:cHGVS'
4. extra columns `自动化+人工证据`, `自动化+人工致病等级` and `PP_References`

## input
`input.tsv` has columns `transcript` and `cHGVS`

## output
`output.tsv` append columns `自动化+人工证据`, `自动化+人工致病等级` and `PP_References`

## how to
```
ACMGlocalAnno -db acmg.local.db.xlsx -sheetName Sheet1 -input input.tsv -output output.tsv
```
or  
```
ACMGlocalAnno -db acmg.local.db.json -sheetName Sheet1 -input input.tsv -output output.tsv
```
