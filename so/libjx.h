/* Code generated by cmd/cgo; DO NOT EDIT. */

/* package tvbox */


#line 1 "cgo-builtin-export-prolog"

#include <stddef.h>

#ifndef GO_CGO_EXPORT_PROLOGUE_H
#define GO_CGO_EXPORT_PROLOGUE_H

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef struct { const char *p; ptrdiff_t n; } _GoString_;
#endif

#endif

/* Start of preamble from import "C" comments.  */




/* End of preamble from import "C" comments.  */


/* Start of boilerplate cgo prologue.  */
#line 1 "cgo-gcc-export-header-prolog"

#ifndef GO_CGO_PROLOGUE_H
#define GO_CGO_PROLOGUE_H

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef size_t GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
#ifdef _MSC_VER
#include <complex.h>
typedef _Fcomplex GoComplex64;
typedef _Dcomplex GoComplex128;
#else
typedef float _Complex GoComplex64;
typedef double _Complex GoComplex128;
#endif

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt.
*/
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef _GoString_ GoString;
#endif
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

#endif

/* End of boilerplate cgo prologue.  */

#ifdef __cplusplus
extern "C" {
#endif


//=====================解析函数=====================
extern char* JsHomeContent(char* c_jsinit, char* c_jsStr);
extern char* JsCategoryContent(char* c_tid, int c_pg, char* c_jsinit, char* c_jsStr);
extern char* JsDetailContent(char* c_ids, char* c_jsinit, char* c_jsStr);
extern char* JsSearchContent(char* c_key, char* c_jsinit, char* c_jsStr);
extern char* JsPlayerContent(char* c_id, char* c_jsinit, char* c_jsStr);

//===================解析结束===================
//初始化文件
extern void GoInit();

//测试动态链接库是否连接成功
extern void GoTest();

//测试动态链接库是否连接成功,是否能够传参
extern void GoTeststr(char* c_str);

/*
*EXT_X_VERSION    int
*EXT_X_KEY_METHOD string
*EXT_X_KEY        string
*EXT_X_KEY_IV     string
*Time_List        []float64
*Ts_list          []string

*解析m3u8
 */
extern char* GoParseM3u8(char* c_url, char* c_header);

//下载Ts片段
extern char* GoDownloadTs(char* c_jstr, int c_index, char* c_path);

//获得uuid
extern char* GoGetUUID();

//递归删除文件
extern void GoRemoveFile(char* c_path);

//发送请求客户端
extern char* GoRequestClient(char* c_URL, char* c_METHOD, char* c_HEADER, char* c_DATA);

#ifdef __cplusplus
}
#endif
