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
	Header3 CustomFloatUnmarshal `xls:"Header 3"`
}

type CustomFloatUnmarshal struct {
	Value string
	Percent float64
}

func (c *CustomFloatUnmarshal)UnmarshalXLS(value string) error {
	c.Value = value
	f , err := toFloat(value)
	if err != nil {
		return err
	}
	c.Percent = f
	return nil
}

const simplefile = "./test/Test1.xlsx"
const sheet = "Sheet 2"

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
	
	f, err := excelize.OpenFile(simplefile)
	assert.NoError(t,err,"Error opening file")
	if err != nil {
		return
	}
	sm := f.GetSheetList()
	ts := test1{}
	assert.NoError(t,UnmarshalExcelize(f,sm[0],&ts))
	fmt.Printf("OneSheet test: %v\n",ts)
}

func TestParseXLSSelectSheet(t *testing.T) {
	f, err := excelize.OpenFile(simplefile)
	assert.NoError(t,err,"Error opening file")
	if err != nil {
		return
	}
	ts := test1{}
	assert.NoError(t,UnmarshalExcelize(f,sheet,&ts))
	fmt.Printf("Sheet Selection %s test: %v\n",sheet,ts)
}

func TestParseXLSFileNameOneSheet(t *testing.T) {
	ts := test1{}
	assert.NoError(t,UnmarshalFile(simplefile,sheet,&ts))
	fmt.Printf("Filename based, sheet selection %s test: %v\n",sheet,ts)
}
