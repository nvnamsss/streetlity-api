package main

import (
	"fmt"

	_ "example.com/m/v2/src/himompkg"
)

var a int = c
var b int = a
var c int = 1

type Page struct {
	pageNum int
	title   string
	data    []byte
}

func init() {
	fmt.Println("hi mom i'm in init")
}

func pratice() {

	var pref *Page //like cpp, go separate reference and value by '*' symbol, it very similar to cpp by the way it initialize
	pref = new(Page)
	pref.pageNum = 2
	pref.title = "hi mom reference"
	var pvalue Page //create value for p, we do not need any constructor

	pvalue.pageNum = 1
	pvalue.title = "hi mom value"

	var samevalue Page
	samevalue = pvalue
	samevalue.title = "hi mom samevalue"

	var sameref *Page
	sameref = pref
	sameref.title = "hi mom sameref"

	fmt.Println(sameref)
	fmt.Println(samevalue)
	fmt.Println(pref.title)
	fmt.Println(pvalue.title)

}

func main() {
	fmt.Println("hi mom i'm in main")
	pratice()
}
