Decode XLS to Go struct
=====

decxls is a module designed to assist in the decoding of xls based spreadsheets into go slice * structs.  Based upon the initial work w/ CSV parsing from [https://github.com/gocarina/gocsv] wrapping the excel reader from [https://github.com/360EntSecGroup-Skylar/excelize].

Finding that too frequently, we are pulling everything from dictionaries to data, we used basic encoding styled libraries with struct tags "xls:" to map header rows to structs.

This module supports embedded structs and pointers as well as typical "unmarshalling" methodes TypeXLSUnmarshall to allow non-native type treatment for underlying data.

Installation
=====
```go get -u github.com/dhushon/decxls
```

shift to directory...

Basic Spreadsheets
=====
in the test directory, there is an Excel spreadsheet used by the go test functions.  The spreadsheet has two sheets and a small number of rows and columns.  By and large, Sheet1 and 2 are the same, but using out of order headers on the second sheet to test tag<->field mapping.

![sample spreadhseet](https://github.com/dhushon/decxls/blob/main/doc/images/ExcelGrab.jpg?raw=true "Test1.xlsx")



Unmarshalling with GoLang
=====
There is an example of Unmarshalling in the [xls_test.go](./xls_test.go) file to show the basics.

```package main

import (
	"fmt"
	"github/dhushon/decxls"
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

func main {
	ts := test1{}
	err := UnmarshalFile(simplefile,sheet,&ts))
	fmt.Printf("Filename based, sheet selection %s test: %v\n",sheet,ts)
}
```
