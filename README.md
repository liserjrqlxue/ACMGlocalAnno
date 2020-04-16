# 需求

1.  给定数据库（excel文件）

2.  注释变异添加数据库给定列（匹配变异）

3.  生成excel展示给定列

4.  解读系统读取并展示给定列（IT工作）

# 实现

1.  撰写软件`ACMGlocalAnno`实现解析数据库excel，并对注释结果文件添加给定列内容（数据库不包含的变异列内容为空）

2.  修改`anno2xlsx`输出tier1 excel
    `filter_variants`表头格式，添加给定列输出（`-filter_variants`控制表头，已在源码修改）

# 生信操作

1.  原流程`anno2xlsx` `-snv`\参数的输入文件记为`input.tsv`

2.  运行`ACMGlocalAnno -input input.tsv -output
    output.tsv`，`output.tsv`为`ACMGlocalAnno`输出

3.  原流程`anno2xlsx` `-snv`参数替换为`output.tsv`

# 细节

1.  通过`Transcript:cHGVS`匹配变异

2.  新增表头：

  注释表表头                                   |ACMG定点数据库表头
  --------------------------------------------| --------------------
  SecondaryFinding\_Var\_证据项                |证据项
  SecondaryFinding\_Var\_致病等级              |致病等级
  SecondaryFinding\_Var\_参考文献              |参考文献
  SecondaryFinding\_Var\_Phenotype\_OMIM\_ID   |关联疾病表型OMIM号
  SecondaryFinding\_Var\_DiseaseNameEN         |关联疾病英文名称
  SecondaryFinding\_Var\_DiseaseNameCH         |关联疾病中文名称
  SecondaryFinding\_Var\_updatetime            |数据库时间
