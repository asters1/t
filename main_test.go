package main

import (
	"fmt"
	"strings"
	"testing"
)

type Cs struct {
	a string
	b int
	C []string
	D []int
	E []float64
}

func TestMain(t *testing.T) {
	//	GoInit()
	a := "http://baidu.com/a/c"
	a = a[7:]
	host := a[:strings.Index(a, "/")]
	fmt.Println(host)

	//	a := Cs{a: "ahshhsh", b: 4689, C: []string{"dd", "kk", "ll"}, D: []int{7, 9, 4}, E: []float64{3.66, 7.98, 4.88}}
	//	res := GetJson(a)
	//	key := `#EXT-X-KEY:METHOD=AES-128,URI="/m3u8key/eHlxVFVzZDc5MVdrYTBuMnV0c2I2NThKUzZWM3JWN0o=.key"`
	//	fmt.Println(res)
	//	m := parseLineParameters(key)
	//	fmt.Println(m)
	//	u := completionUrl("http://a.c/kk/aa.m3u8", "://hh")
	//	fmt.Println(u)

	//	h := `User-Agent: Wget/1.21.3`
	//	res := ParseM3u8("http://m35.grelighting.cn/m3u8/N2t1endXZ0F0K3o2aTZRNWNqNGNWTE1iaDRuTmlTUkE=.m3u8", h)
	//	res := ParseM3u8("http://a.zhaojiuwanwu.top/api/GetDownUrlMu/3bb24322f78b47dfb8723c13d46d45ee/a8bc6c3deccd4f4ab99ecea367d6a2d2.m3u8?sign=2b4715f5158f4612323318a5284ba8a9&t=1686208863", h)
	//	fmt.Println(res)
	//	r := DownloadTs(res, 0, "./cache/ts/cs.ts")
	//	fmt.Println("r:" + r)
	//	time := gjson.Get(res, "Time_List").String()
	//	fmt.Println(time)
	//	fmt.Println(time)
	//	time := `[0.2,1.2,3.7,5.7,6.1,6.2,7.3,7.5,8.2,8.3,8.5,9.1,9.4,9.8,10,11]`
	//	TwoFind(time, "9")

	///jstr_byte, _ := ioutil.ReadFile("./source.json")
	//jstr := string(jstr_byte)
	//a := gjson.Get(jstr, "init").Str
	//b := gjson.Get(jstr, "homeContent").Str
	////	c := gjson.Get(jstr, "categoryContent").Str
	////	d := gjson.Get(jstr, "detailContent").Str
	////	e := gjson.Get(jstr, "playerContent").Str
	//s := gjson.Get(jstr, "searchContent").Str
	//JsInit(a)
	//res := JsHomeContent(a, b)
	//fmt.Println(res)
	////	tid := "/vodshow/4--------"
	////	pg := 1
	////	res = JsCategoryContent(tid, pg, a, c)
	////	fmt.Println(res)
	////	ids := "/video/130539.html"
	////	res = JsDetailContent(ids, a, d)
	////	fmt.Println(res)
	////	id := "/vplay/130539-1-1.html"
	////	res = JsPlayerContent(id, a, e)
	////	fmt.Println(res)
	//key := "斗罗大陆"
	//res = JsSearchContent(key, a, s)

	//fmt.Println(res)
}
