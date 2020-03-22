package main

import (
	"fmt"
	"net/http"

	"example.com/m/v2/Astar"
	_ "example.com/m/v2/himompkg"
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

func hifive(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Hi five")
	w.Write([]byte("hi fiveeeee"))
}

func http_practive() {
	var m meomeo
	http.HandleFunc("/hifive", hifive)
	http.ListenAndServe(":9000", m)
}

type meomeo struct {
}

func (m meomeo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("meomeo")
}

func chandler(w http.ResponseWriter, req *http.Request) {
	fmt.Println(req.URL)
}

func emptyInterface() {
	var i struct {
		a int
	}
	i.a = 16

	fmt.Println(i.a)
}

func try() {

}

func main() {
	fmt.Println("hi mom i'm in main")
	// emptyInterface()
	// pratice()

	// var i int64 = 0
	// var loop int = 2
	// for loop = 2; loop <= 10; loop++ {
	// 	i += int64(math.Pow10(loop)*9/10) * int64(loop-1)
	// }

	// fmt.Println("i:", i)
	PrepareData()
	path, ok := Astar.Route(Astar.Nodes[6972965259], Astar.Nodes[6972965899])
	if ok {
		fmt.Println(path)
	}
	// http_practive()
}
