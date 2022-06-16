package main

import (
	"github.com/zheng11581/cloudnative/cncamp/utils"
)

func main() {
	indexes := utils.GetIndex()
	for _, index := range indexes {
		utils.DeleteIndex(index)
	}
}
