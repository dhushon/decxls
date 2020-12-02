package main

import (
	"fmt"
	"go-decxls/decxls"
	"go-decxls/livd"
	"github.com/360EntSecGroup-Skylar/excelize"
)

func main() {
	const sheet = "LOINC Mapping"
	f, err := excelize.OpenFile("./test/LIVD-SARS-CoV-2-2020-11-18.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	if !(func(m map[int]string, v string) bool {
		for _, x := range m {
			if x == v {
				return true
			}
		}
		return false
	} (f.GetSheetMap() , sheet)) {
		fmt.Printf("Sheet %s not found in file: %v\n",sheet,f.GetSheetMap())
	}
	l := []*livd.Livd{}
	err = decxls.UnmarshalExcelize(f,sheet,&l)
	fmt.Println("complete")
}
