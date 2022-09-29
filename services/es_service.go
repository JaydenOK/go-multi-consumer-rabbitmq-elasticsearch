package services

import (
	"app/libs/elasticsearchlib"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type EsService struct {
}

func (esService *EsService) IndexLists(ctx *gin.Context) interface{} {
	indexLists := elasticsearchlib.IndexLists()
	return indexLists
}

func (esService *EsService) IndexExist(ctx *gin.Context) interface{} {
	index := ctx.Query("index")
	return elasticsearchlib.IndexExist(index)
}

func (esService *EsService) IndexCreate(ctx *gin.Context) interface{} {
	index := ctx.Query("index")
	rawData, _ := ctx.GetRawData()
	mapping := map[string]interface{}{}
	_ = json.Unmarshal(rawData, &mapping)
	if index == "" {
		return "index不能为空"
	}
	if elasticsearchlib.IndexExist(index) {
		return "Index已存在：" + index
	}
	ret := elasticsearchlib.IndexCreate(index, mapping)
	return ret
}

func (esService *EsService) IndexGetMapping(ctx *gin.Context) interface{} {
	index := ctx.Query("index")
	if index == "" {
		return "index不能为空"
	}
	mapping := elasticsearchlib.IndexGetMapping(index)
	return mapping
}

//创建mapping文档结构

func (esService *EsService) IndexPutMapping(ctx *gin.Context) interface{} {
	index := ctx.Query("index")
	rawData, _ := ctx.GetRawData()
	mapping := map[string]interface{}{}
	_ = json.Unmarshal(rawData, &mapping)
	if index == "" || rawData == nil {
		return "参数错误"
	}
	ret := elasticsearchlib.IndexPutMapping(index, mapping)
	return ret
}

// es中索引的字段类型是不可修改的，只能是重新创建一个索引并设置好mapping，然后再将老索引的数据复制过去
func (esService *EsService) IndexReindex(ctx *gin.Context) interface{} {
	sourceIndex := ctx.Query("source_index")
	destIndex := ctx.Query("dest_index")
	ret := elasticsearchlib.IndexReindex(sourceIndex, destIndex)
	return ret
}

func (esService *EsService) IndexDelete(ctx *gin.Context) interface{} {
	index := ctx.Query("index")
	ret := elasticsearchlib.IndexDelete(index)
	return ret
}

func (esService *EsService) IndexAliasLists(ctx *gin.Context) interface{} {
	index := ctx.Query("index")
	ret := elasticsearchlib.IndexAliasLists(index)
	return ret
}

func (esService *EsService) IndexAlias(ctx *gin.Context) interface{} {
	index := ctx.Query("index")
	alias := ctx.Query("alias")
	ret := elasticsearchlib.IndexAlias(index, alias)
	return ret
}
