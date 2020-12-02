package decxls

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	excelize "github.com/360EntSecGroup-Skylar/excelize/v2"
)

type Test1 struct {
	Header1 string `xls:"Header 1"`
	Header2 int `xls:"Header 2"`
	Header3 float32 `xls:"Header 3"`
}

const simplefile = "./test/Test1.xlsx"

type test1 []*Test1

func (tests test1) String() string {
    s := "["
    for i, test := range tests {
        if i > 0 {
            s += ", "
        }
        s += fmt.Sprintf("%v", test)
    }
    return s + "]"
}

// TestNormalizeString serves to excercise the string normalizer used to compare key structures 
func TestParseXLSOneSheet(t *testing.T) {

	const sheet = "LOINC Mapping"
	f, err := excelize.OpenFile(simplefile)
	assert.NoError(t,err,"Error opening file")
	if err != nil {
		return
	}
	sm := f.GetSheetList()
	ts := test1{}
	assert.NoError(t,UnmarshalExcelize(f,sm[0],&ts))
	fmt.Printf("OneSheet test: %v",ts)
}
