package main

import (
	"fmt"
	"net/http"
)

type s1 struct {
	name string
	id   int
}

type s2 struct {
	name string
	id   int
}

type i1 interface {
	m1()
}

type i2 interface {
	m2()
}

type str struct {
	meo int
}

func (s str) m1() {
	fmt.Println("[Multiple Interface]", "hi mom i'm i1")
}

func (s str) m2() {
	fmt.Println("[Multiple Interface]", "hi mom i'm i2")
}

func multiple_returnvalue(s string, i int) (int, int, int) {
	return 1, 2, 3
}

func casting() {
	var ss s1
	ss.name = "abc"
	ss.id = 10
	fmt.Println(s2(ss))
}

func pointer() {
	a := 5
	b := &a
	var c *int = new(int)
	*c = 3
	*b = *c

	fmt.Println(a)
	fmt.Println(*b)
}

func multiple_dimension_array() {
	var a [][]int
	a[0] = []int{7, 2, 3, 4, 5}

	fmt.Println("[Multiple Dimension Array]", a[0])
	fmt.Println("[Multiple Dimension Array]", len(a))
}

func foreach() {

	var a []int = []int{7, 2, 3, 4, 5}

	for index, element := range a {
		// index is the index where we are
		// element is the element from someSlice for where we are
		fmt.Println("[Foreach]", index, element)
	}
}
func forloop() {
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
	}
	fmt.Println("[For loop]", sum)
}

func cast_interface() {
	var s str

	// var io i1 = s
	// var it i2 = s

	i1(s).m1()
	i2(s).m2()

}

func hashtable() {
	m2 := map[string]int{
		"a": 1,
		"b": 2,
	}

	for key, value := range m2 {
		fmt.Println("[Hashtable]", key, value)
	}
}

func pratice() {

	// var pref *Page //like cpp, go separate reference and value by '*' symbol, it very similar to cpp by the way it initialize
	// pref = new(Page)
	// pref.pageNum = 2
	// pref.title = "hi mom reference"
	// var pvalue Page //create value for p, we do not need any constructor

	// pvalue.pageNum = 1
	// pvalue.title = "hi mom value"

	// var samevalue Page
	// samevalue = pvalue
	// samevalue.title = "hi mom samevalue"

	// var sameref *Page
	// sameref = pref
	// sameref.title = "hi mom sameref"

	// fmt.Println(sameref)
	// fmt.Println(samevalue)
	// fmt.Println(pref.title)
	// fmt.Println(pvalue.title)

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

func init() {
	fmt.Println("multiple_iterface")

	var v1, v2, v3 int
	v1, v2, v3 = multiple_returnvalue("1", 1)
	cast_interface()
	foreach()
	hashtable()
	forloop()
	fmt.Println(v1, v2, v3)
	// io.m1()
	// it.m2()
}
