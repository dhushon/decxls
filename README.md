Decode XLS to Go struct
=====

decxls is a module designed to assist in the decoding of xls based spreadsheets into go slice * structs.  Based upon the initial work w/ CSV parsing from https://github.com/gocarina/gocsv wrapping the excel reader from https://github.com/360EntSecGroup-Skylar/excelize.

Finding that too frequently, we are pulling everything from dictionaries to data, we used basic encoding styled libraries with struct tags "xls:" to map header rows to structs.

This module supports embedded structs and pointers as well as typical "unmarshalling" methodes TypeXLSUnmarshall to allow non-native type treatment for underlying data.

Installation
=====
```go get -u github.com/dhushon/decxls```
shift to directory...
```go get```

Basic Spreadsheets
=====
in the test directory, there is an Excel spreadsheet used by the go test functions.  The spreadsheet has two sheets and a small number of rows and columns.  By and large, Sheet1 and 2 are the same, but using out of order headers on the second sheet to test tag<->field mapping.


Unmarshalling with GoLang
=====
There is an example of Unmarshalling in the xls_test.go file to show the basics.


