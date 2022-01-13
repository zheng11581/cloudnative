package main

import (
	"cncamp/utils"
)

func main() {
	indexes := utils.GetIndex()
	for _, index := range indexes {
		utils.DeleteIndex(index)
	}
}
