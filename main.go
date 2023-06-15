package main

import (
	"C"
	"fmt"
	"os"
)
import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/asters1/goquery"
	"github.com/robertkrimen/otto"
	uuid "github.com/satori/go.uuid"
	"github.com/tidwall/gjson"
)

type m3u8 struct {
	EXT_X_VERSION    int
	EXT_X_KEY_METHOD string
	//KEY_METHOD====> 'AES-128' or 'NONE'
	//如果加密方法为 NONE，则 URI 和 IV 属性不得存在
	EXT_X_KEY    string
	EXT_X_KEY_IV string
	HEADER       string
	IsLive       bool
	Time_List    []float64
	Ts_list      []string
}

// func C.CString(string) *C.char              //go字符串转化为char*
// func C.CBytes([]byte) unsafe.Pointer        // go 切片转化为指针
// func C.GoString(*C.char) string             //C字符串 转化为 go字符串
// func C.GoStringN(*C.char, C.int) string
// func C.GoBytes(unsafe.Pointer, C.int) []byte

//=====================解析=====================
//otto Js初始化
func JsInit(jsStr string) *otto.Otto {
	//func JsInit(c_jsStr *C.char) {

	//初始化用户自定义的js
	//	jsStr := C.GoString(c_jsStr)
	vm := otto.New()

	//随机数,4位小数
	vm.Set("go_random", func(call otto.FunctionCall) otto.Value {

		rand := int(time.Now().UnixNano() % 10000)

		frand := float32(rand) * 0.0001
		value, _ := strconv.ParseFloat(fmt.Sprintf("%.4f", frand), 64)
		vm_random, _ := vm.ToValue(value)
		return vm_random

	})
	//获取md5
	vm.Set("go_md5", func(call otto.FunctionCall) otto.Value {
		str, _ := call.Argument(0).ToString()
		data := []byte(str) //切片
		has := md5.Sum(data)
		md5str := fmt.Sprintf("%x", has) //将[]byte转成16进制
		vm_md5, _ := vm.ToValue(md5str)
		return vm_md5

	})
	//获取时间戳
	vm.Set("go_getTime", func(call otto.FunctionCall) otto.Value {
		i, _ := call.Argument(0).ToInteger()
		if i > 19 {
			i = 19
		}
		timeUnixNano := time.Now().UnixNano()
		str_time := strconv.FormatInt(timeUnixNano, 10)
		s, _ := vm.ToValue(str_time[:i])
		return s
	})

	//发送请求
	vm.Set("go_RequestClient", func(call otto.FunctionCall) otto.Value {
		FormatStr := func(jsonstr string) map[string]string {
			DataMap := make(map[string]string)
			Nslice := strings.Split(jsonstr, "\n")
			for i := 0; i < len(Nslice); i++ {
				if strings.Index(Nslice[i], ":") != -1 {
					a := Nslice[i][:strings.Index(Nslice[i], ":")]
					b := Nslice[i][strings.Index(Nslice[i], ":")+1:]
					DataMap[a] = b
				}
			}
			return DataMap

		}

		URL, _ := call.Argument(0).ToString()
		METHOD, _ := call.Argument(1).ToString()
		HEADER, _ := call.Argument(2).ToString()
		DATA, _ := call.Argument(3).ToString()

		URL = strings.TrimSpace(URL)
		METHOD = strings.TrimSpace(METHOD)
		HEADER = strings.TrimSpace(HEADER)
		DATA = strings.TrimSpace(DATA)
		if URL == "" || METHOD == "" {
			return otto.Value{}
		}

		HeaderMap := FormatStr(HEADER)
		DataMap := FormatStr(DATA)
		client := &http.Client{}
		if METHOD == "get" {
			METHOD = http.MethodGet
		} else if METHOD == "post" {
			METHOD = http.MethodPost

		}
		FormatData := ""
		for i, j := range DataMap {
			FormatData = FormatData + i + "=" + j + "&"
		}
		if FormatData != "" {
			FormatData = FormatData[:len(FormatData)-1]
		}
		requset, _ := http.NewRequest(
			METHOD,
			URL,
			strings.NewReader(FormatData),
		)
		if METHOD == http.MethodPost && requset.Header.Get("Content-Type") == "" {
			requset.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
		requset.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.71 Safari/537.36")
		for i, j := range HeaderMap {
			requset.Header.Set(i, j)
		}
		resp, _ := client.Do(requset)
		body_bit, _ := ioutil.ReadAll(resp.Body)
		headerMap := resp.Header
		jsonByte, err := json.Marshal(headerMap)
		if err != nil {
			fmt.Printf("Marshal with error: %+v\n", err)
		}
		header := string(jsonByte)

		defer resp.Body.Close()
		status := strconv.Itoa(resp.StatusCode)
		body := string(body_bit)
		res_str := make(map[string]string)

		res_str["status"] = status
		res_str["header"] = header
		res_str["body"] = body
		result, _ := vm.ToValue(res_str)

		return result
	})
	vm.Set("go_FindJsonKey", func(call otto.FunctionCall) otto.Value {
		JSON, _ := call.Argument(0).ToString()
		Key, _ := call.Argument(1).ToString()
		value := gjson.Get(JSON, Key).String()
		result, _ := vm.ToValue(value)
		return result

	})
	//	vm.Set("go_FindJsonKeyArray", func(call otto.FunctionCall) otto.Value {
	//		JSON, _ := call.Argument(0).ToString()
	//		Key, _ := call.Argument(1).ToString()
	//		value := gjson.Get(JSON, Key).Array()
	//		fmt.Println(value)
	//		result, _ := vm.ToValue(value)
	//		return result
	//
	//	})
	vm.Set("go_FindHtml", func(call otto.FunctionCall) otto.Value {
		HTML, _ := call.Argument(0).ToString()
		CSS, _ := call.Argument(1).ToString()

		// 加载 HTML document对象
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(HTML))
		if err != nil {
			fmt.Println("加载HTML失败")
			os.Exit(0)
		}
		var res []string
		// Find the review items
		doc.Find(CSS).Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			value, err := s.String()

			if err != nil {
				fmt.Println("获取Html列表出错")
			}

			res = append(res, value)
		})

		result, _ := vm.ToValue(res)
		return result
	})
	vm.Set("go_FindText", func(call otto.FunctionCall) otto.Value {
		HTML, _ := call.Argument(0).ToString()
		CSS, _ := call.Argument(1).ToString()

		// 加载 HTML document对象
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(HTML))
		if err != nil {
			fmt.Println("加载节点失败")
			os.Exit(0)
		}
		// Find the review items
		res := doc.Find(CSS).Text()
		result, _ := vm.ToValue(res)
		return result
	})
	vm.Set("go_FindAttr", func(call otto.FunctionCall) otto.Value {
		HTML, _ := call.Argument(0).ToString()
		CSS, _ := call.Argument(1).ToString()
		KEY, _ := call.Argument(2).ToString()

		// 加载 HTML document对象
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(HTML))
		if err != nil {
			fmt.Println("加载节点失败")
			os.Exit(0)
		}
		// Find the review items
		// For each item found, get the band and title
		res := ""
		bl := true
		if CSS == "" {
			res, bl = doc.Attr(KEY)

		} else {
			res, bl = doc.Find(CSS).Attr(KEY)
		}

		if !bl {
			fmt.Println("没有获取到" + KEY)
		}
		result, _ := vm.ToValue(res)
		return result
	})
	//	fmt.Println(jsStr)
	_, err := vm.Run(jsStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	//	fmt.Println(a)

	return vm
}

//=====================解析函数=====================
//export JsHomeContent
func JsHomeContent(c_jsinit *C.char, c_jsStr *C.char) *C.char {

	jsinit := C.GoString(c_jsinit)
	jsStr := C.GoString(c_jsStr)
	vm := JsInit(jsinit)
	res, err := vm.Run(jsStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	return C.CString(res.String())
}

//export JsCategoryContent
func JsCategoryContent(c_tid *C.char, c_pg C.int, c_jsinit *C.char, c_jsStr *C.char) *C.char {
	jsStr := C.GoString(c_jsStr)
	jsinit := C.GoString(c_jsinit)
	tid := C.GoString(c_tid)
	pg := int(c_pg)
	vm := JsInit(jsinit)
	vm.Set("tid", tid)
	vm.Set("pg", pg)
	res, err := vm.Run(jsStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	return C.CString(res.String())

}

//export JsDetailContent
func JsDetailContent(c_ids *C.char, c_jsinit *C.char, c_jsStr *C.char) *C.char {
	jsStr := C.GoString(c_jsStr)
	ids := C.GoString(c_ids)
	jsinit := C.GoString(c_jsinit)
	vm := JsInit(jsinit)
	vm.Set("ids", ids)
	res, err := vm.Run(jsStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	return C.CString(res.String())

}

//export JsSearchContent
func JsSearchContent(c_key *C.char, c_jsinit *C.char, c_jsStr *C.char) *C.char {
	key := C.GoString(c_key)
	jsinit := C.GoString(c_jsinit)
	jsStr := C.GoString(c_jsStr)

	vm := JsInit(jsinit)
	vm.Set("key", key)
	res, err := vm.Run(jsStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	return C.CString(res.String())
}

//export JsPlayerContent
func JsPlayerContent(c_id *C.char, c_jsinit *C.char, c_jsStr *C.char) *C.char {
	id := C.GoString(c_id)
	jsinit := C.GoString(c_jsinit)
	jsStr := C.GoString(c_jsStr)

	vm := JsInit(jsinit)
	vm.Set("id", id)
	res, err := vm.Run(jsStr)
	if err != nil {
		fmt.Println(err.Error())
	}
	return C.CString(res.String())
}

//===================解析结束===================
//初始化文件
//export GoInit
func GoInit() {
	RemoveFile("cache")
	url := "https://jihulab.com/asters1/config/-/raw/master/"
	dir_list := []string{"cache/ts", "so", "config", "bin", "cache/img", "default"}
	file_list := []string{"so/json.hpp", "so/json.cpp", "default/pic.jpg"}
	//检查Dir
	for i := 0; i < len(dir_list); i++ {
		if !IsExists(dir_list[i]) {
			os.MkdirAll(dir_list[i], 0777)
		}
	}
	for i := 0; i < len(file_list); i++ {
		if !IsExists(file_list[i]) {

			if len(file_list[i]) > 3 && file_list[i][:3] == "so/" {
				GetFile(url+file_list[i][3:], file_list[i])

			} else if len(file_list[i]) > 7 && file_list[i][:7] == "config/" {
				GetFile(url+file_list[i][7:], file_list[i])

			} else if len(file_list[i]) > 8 && file_list[i][:8] == "default/" {
			}
			//			GetFile(url+file_list[i][8:], file_list[i])

		}

	}
	//检查json.hpp

}
func FormatJStr(jstr string) string {
	str_map := make(map[string]string)
	//	str_map[`"`] = `\"`
	//	str_map["\t"] = ``
	str_map[`\n`] = `\\n`
	//	str_map[` `] = `ddd`
	for k, v := range str_map {

		jstr = strings.ReplaceAll(jstr, k, v)
	}
	return jstr

}

//测试动态链接库是否连接成功
//export GoTest
func GoTest() {
	fmt.Println("外部链接库，链接成功！")
}

//测试动态链接库是否连接成功,是否能够传参
//export GoTeststr
func GoTeststr(c_str *C.char) {
	str := C.GoString(c_str)
	fmt.Print("链接库测试函数：")
	fmt.Println(str)

}

/*获得一个结构体的json串
 *不能嵌套结构体
 *支持得类型int,string,float64,bool
 *结构体成员为切片类型是时,首字母必须大写,否则会报错
 */
func GetJson(x any) string {

	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)

	//	fmt.Println(t.Kind())
	result := ""
	res := []string{}
	for i := 0; i < t.NumField(); i++ {
		//		fmt.Println(t.Field(i).Type.String())
		switch t.Field(i).Type.String() {

		case "string":

			res = append(res, `"`+t.Field(i).Name+`":"`+v.Field(i).String()+`"`)
		case "int":
			res = append(res, `"`+t.Field(i).Name+`":`+strconv.Itoa(int(v.Field(i).Int())))
		case "bool":
			res = append(res, `"`+t.Field(i).Name+`":`+strconv.FormatBool(v.Field(i).Bool()))
		case "float64":
			res = append(res, `"`+t.Field(i).Name+`":`+strconv.FormatFloat(v.Field(i).Float(), 'f', 2, 64))
		case "[]string":
			obj := v.Field(i).Interface()
			str_list := obj.([]string)
			//			data := "["
			for i := 0; i < len(str_list); i++ {
				str_list[i] = `"` + str_list[i] + `"`
			}
			data := `[` + strings.Join(str_list, ",") + "]"
			res = append(res, `"`+t.Field(i).Name+`":`+data)
		case "[]int":
			obj := v.Field(i).Interface()
			int_list := obj.([]int)
			str_list := []string{}
			//			data := "["
			for i := 0; i < len(int_list); i++ {
				str_list = append(str_list, strconv.Itoa(int_list[i]))
			}
			data := `[` + strings.Join(str_list, ",") + "]"
			res = append(res, `"`+t.Field(i).Name+`":`+data)

		case "[]float64":
			obj := v.Field(i).Interface()
			float64_list := obj.([]float64)
			str_list := []string{}
			//			data := "["
			for i := 0; i < len(float64_list); i++ {
				str_list = append(str_list, strconv.FormatFloat(float64_list[i], 'f', 2, 64))

			}
			data := `[` + strings.Join(str_list, ",") + "]"
			res = append(res, `"`+t.Field(i).Name+`":`+data)

			//		default:
			//			fmt.Println(t.Field(i).Type.String())
		}
	}
	result = "{" + strings.Join(res, ",") + "}"
	result = strings.ReplaceAll(result, "\n", "\\n")
	//	fmt.Println(strings.Index(result, `\`))

	return result
}

/*
*EXT_X_VERSION    int
*EXT_X_KEY_METHOD string
*EXT_X_KEY        string
*EXT_X_KEY_IV     string
*Time_List        []float64
*Ts_list          []string

*解析m3u8
 */
//export GoParseM3u8
func GoParseM3u8(c_url *C.char, c_header *C.char) *C.char {
	url := C.GoString(c_url)
	header := C.GoString(c_header)

	m := m3u8{}
	m.EXT_X_KEY_METHOD = "NONE"
	m.HEADER = header
	m.EXT_X_KEY = ""
	m.EXT_X_KEY_IV = ""
	m.IsLive = true
	time := 0.0

	res := RequestClient(url, "get", header, "")
	status_str := res["status"]
	status, _ := strconv.Atoi(status_str)
	if !(status > 199 && status < 300) {
		return C.CString("")
	}

	list := strings.Split(res["body"], "\n")
	//	fmt.Println(list)
	if strings.TrimSpace(list[0]) != `#EXTM3U` {
		return C.CString("")
	}
	for i := 1; i < len(list); i++ {
		if !m.IsLive {
			break
		}
		line := strings.TrimSpace(list[i])
		switch {
		case strings.HasPrefix(line, "#EXT-X-VERSION:"):
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			str_version := list[i][strings.Index(list[i], `:`)+1:]
			//			fmt.Println("str:" + str_version)
			m.EXT_X_VERSION, _ = strconv.Atoi(str_version)
		case strings.HasPrefix(line, "#EXT-X-KEY"):
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			params := parseLineParameters(line)
			m.EXT_X_KEY_METHOD = params["METHOD"]
			KEY_URL := completionUrl(url, params["URI"])
			r := RequestClient(KEY_URL, "get", header, "")
			m.EXT_X_KEY = r["body"]

			m.EXT_X_KEY_IV = params["IV"]
		case strings.HasPrefix(line, "#EXTINF:"):
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			list[i] = strings.TrimSpace(list[i])
			str_time := list[i][strings.Index(list[i], `:`)+1 : len(list[i])-1]
			f_time, _ := strconv.ParseFloat(str_time, 32)
			//			fmt.Println(str_time)

			time = time + f_time
			//			fmt.Println(i)
			t, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", time), 64)

			m.Time_List = append(m.Time_List, t)
			//			fmt.Println(t)
		case !strings.HasPrefix(line, "#"):
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			ts := completionUrl(url, line)

			m.Ts_list = append(m.Ts_list, ts)
		case line == "#EXT-X-ENDLIST":
			m.IsLive = false

		}
	}
	if len(m.Ts_list) != len(m.Time_List) {
		return C.CString("")
	}
	return C.CString(GetJson(m))
}

////export teststr
//func teststr(c_str *C.char) {
//	str := C.GoString(c_str)
//	fmt.Print("链接库测试函数：")
//	fmt.Println(str)
//
//}

//下载Ts片段
//export GoDownloadTs
func GoDownloadTs(c_jstr *C.char, c_index C.int, c_path *C.char) *C.char {
	jstr := C.GoString(c_jstr)
	index := int(c_index)
	path := C.GoString(c_path)
	//	fmt.Println(jstr)
	dir := path[:strings.LastIndex(path, "/")]
	//	fmt.Println(dir)
	if !IsExists(dir) {
		os.MkdirAll(dir, 0777)
	}
	list := gjson.Get(jstr, "Ts_list").Array()
	header := gjson.Get(jstr, "HEADER").String()
	key := gjson.Get(jstr, "EXT_X_KEY").String()
	iv := gjson.Get(jstr, "EXT_X_KEY_IV").String()

	//	fmt.Println(header)
	//	u := list[index].Str
	res := RequestClient(list[index].String(), "get", header, "")
	bytes := []byte(res["body"])
	if key != "" {
		b, err := AES128Decrypt(bytes, []byte(key), []byte(iv))
		if err != nil {
			return C.CString("")
		}
		bytes = b
		// Some TS files do not start with SyncByte 0x47,
		// 一些 ts 文件不以同步字节 0x47 开头，
		//	they can not be played after merging,
		// 合并后不能播放，
		// Need to remove the bytes before the SyncByte 0x47(71).
		// 需要删除同步字节 0x47(71) 之前的字节。
	}

	syncByte := uint8(71) //0x47
	bLen := len(bytes)
	for j := 0; j < bLen; j++ {
		if bytes[j] == syncByte {
			//			fmt.Println(bytes[:j])
			bytes = bytes[j:]
			break
		}
	}
	file, _ := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	defer file.Close()
	file.Write(bytes)

	return c_path

}
func parseLineParameters(line string) map[string]string {
	var linePattern = regexp.MustCompile(`([a-zA-Z-]+)=("[^"]+"|[^",]+)`)
	r := linePattern.FindAllStringSubmatch(line, -1)
	params := make(map[string]string)
	for _, arr := range r {
		params[arr[1]] = strings.Trim(arr[2], "\"")
	}
	return params
}
func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else {
		return false
	}

}

//获得uuid
//export GoGetUUID
func GoGetUUID() *C.char {

	u1 := uuid.NewV4()
	return C.CString(u1.String())
}
func FormatStr(jsonstr string) map[string]string {
	DataMap := make(map[string]string)
	Nslice := strings.Split(jsonstr, "\n")
	for i := 0; i < len(Nslice); i++ {
		if strings.Index(Nslice[i], ":") != -1 {
			a := Nslice[i][:strings.Index(Nslice[i], ":")]
			b := Nslice[i][strings.Index(Nslice[i], ":")+1:]
			DataMap[a] = b
		}
	}
	return DataMap
}

//递归删除文件
//export GoRemoveFile
func GoRemoveFile(c_path *C.char) {
	path := C.GoString(c_path)
	os.RemoveAll(path)
}
func RemoveFile(path string) {
	os.RemoveAll(path)
}

//发送请求客户端
//export GoRequestClient
func GoRequestClient(c_URL *C.char, c_METHOD *C.char, c_HEADER *C.char, c_DATA *C.char) *C.char {

	URL := C.GoString(c_URL)
	METHOD := C.GoString(c_METHOD)
	HEADER := C.GoString(c_HEADER)
	DATA := C.GoString(c_DATA)

	URL = strings.TrimSpace(URL)
	METHOD = strings.TrimSpace(METHOD)
	HEADER = strings.TrimSpace(HEADER)
	DATA = strings.TrimSpace(DATA)
	if URL == "" || METHOD == "" {
		fmt.Println("URL或者METHOD为空!")
		return nil
	}
	HeaderMap := FormatStr(HEADER)
	DataMap := FormatStr(DATA)
	client := &http.Client{}
	if METHOD == "get" {
		METHOD = http.MethodGet
	} else if METHOD == "post" {
		METHOD = http.MethodPost

	}
	FormatData := ""
	for i, j := range DataMap {
		FormatData = FormatData + i + "=" + j + "&"
	}
	if FormatData != "" {
		FormatData = FormatData[:len(FormatData)-1]
	}
	requset, _ := http.NewRequest(
		METHOD,
		URL,
		strings.NewReader(FormatData),
	)
	if METHOD == http.MethodPost && requset.Header.Get("Content-Type") == "" {
		requset.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	requset.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.71 Safari/537.36")
	for i, j := range HeaderMap {
		requset.Header.Set(i, j)
	}
	resp, _ := client.Do(requset)
	body_bit, _ := ioutil.ReadAll(resp.Body)
	headerMap := resp.Header
	jsonByte, err := json.Marshal(headerMap)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
	}
	header := string(jsonByte)

	defer resp.Body.Close()
	status := strconv.Itoa(resp.StatusCode)
	body := string(body_bit)
	res_str := `{"status":` + status + `,"header":"` + header + `","body":"` + body + `"}`

	return C.CString(res_str)
}

func RequestClient(URL string, METHOD string, HEADER string, DATA string) map[string]string {
	URL = strings.TrimSpace(URL)
	METHOD = strings.TrimSpace(METHOD)
	HEADER = strings.TrimSpace(HEADER)
	DATA = strings.TrimSpace(DATA)
	if URL == "" || METHOD == "" {
		fmt.Println("URL或者METHOD为空!")
		return nil
	}
	HeaderMap := FormatStr(HEADER)
	DataMap := FormatStr(DATA)
	client := &http.Client{}
	if METHOD == "get" {
		METHOD = http.MethodGet
	} else if METHOD == "post" {
		METHOD = http.MethodPost

	}
	FormatData := ""
	for i, j := range DataMap {
		FormatData = FormatData + i + "=" + j + "&"
	}
	if FormatData != "" {
		FormatData = FormatData[:len(FormatData)-1]
	}
	requset, _ := http.NewRequest(
		METHOD,
		URL,
		strings.NewReader(FormatData),
	)
	if METHOD == http.MethodPost && requset.Header.Get("Content-Type") == "" {
		requset.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	requset.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.5060.71 Safari/537.36")
	for i, j := range HeaderMap {
		requset.Header.Set(i, j)
	}
	resp, _ := client.Do(requset)
	body_bit, _ := ioutil.ReadAll(resp.Body)
	headerMap := resp.Header
	jsonByte, err := json.Marshal(headerMap)
	if err != nil {
		fmt.Printf("Marshal with error: %+v\n", err)
	}
	header := string(jsonByte)

	defer resp.Body.Close()
	status := strconv.Itoa(resp.StatusCode)
	body := string(body_bit)
	res_str := make(map[string]string)

	res_str["status"] = status
	res_str["header"] = header
	res_str["body"] = body
	return res_str
}
func GetFile(url string, path string) {
	res := RequestClient(url, "get", "", "")
	body := res["body"]
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("打开文件失败,错误:", err)
		return
	}
	defer file.Close()
	file.WriteString(body)

}

//自动补全url
func completionUrl(url string, path string) string {

	url = strings.TrimSpace(url)
	u := url[:strings.LastIndex(url, `/`)+1]
	host := ""
	switch {
	case strings.HasPrefix(url, "http://"):
		a := url[7:]
		host = "http://" + a[:strings.Index(a, "/")]
	case strings.HasPrefix(url, "https://"):
		a := url[8:]
		host = "https://" + a[:strings.Index(a, "/")]
	}
	//	fmt.Println(u)

	switch {
	case strings.HasPrefix(path, "http"):
		return path
	case strings.HasPrefix(path, "://"):
		return `http` + path
	case strings.HasPrefix(path, "//"):
		return `http:` + path
	case strings.HasPrefix(path, "/"):
		return host + path
	default:
		return u + path
	}
}

//二分法查找
//jstr 是数组的Json串,s_time是查找值
func TwoFind(jstr string, s_time string) int {

	res := gjson.Get(`[`+jstr+`]`, "0").Array()
	//	fmt.Println("res:" + res)
	f_time, _ := strconv.ParseFloat(s_time, 32)
	t, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", f_time), 64)
	//	fmt.Println("长度为")
	//	fmt.Println(len(res))
	low := 0
	high := len(res) - 1

	mid := 0
	for {
		if low > high {
			break
		}

		mid = (low + high) / 2
		if res[mid].Float() < t {
			low = mid + 1
		} else if res[mid].Float() > t {
			high = mid - 1
		} else {
			//			fmt.Println(mid)
			break
		}
	}
	if res[mid].Float() <= t {
		//		fmt.Println(mid)
		//		fmt.Println(res[mid].Float())
		return mid
	} else {
		if mid == 0 {
			return mid
		}
		return mid - 1
	}

}
func AES128Decrypt(crypted, key, iv []byte) ([]byte, error) {
	//	fmt.Println(crypted)
	//	fmt.Println(key)
	//	fmt.Println(iv)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(iv) == 0 {
		iv = key
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	length := len(origData)
	unPadding := int(origData[length-1])
	origData = origData[:(length - unPadding)]
	return origData, nil
}
func main() {

}
