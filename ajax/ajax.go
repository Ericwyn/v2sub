package ajax

import (
	"bytes"
	"encoding/json"
	"github.com/Ericwyn/v2sub/utils/log"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type Method string

//const
const POST Method = "POST"
const GET Method = "GET"

var client = &http.Client{} //客户端,被Get,Head以及Post使用

type Request struct {
	Url    string
	Method Method

	Data map[string]string // GET 和 POST url编码的数据
	Form map[string]string // POST 时候 FormData 编码的数据
	Json interface{}       // POST 时候 JSON 格式编码的数据

	Header  map[string]string
	Success Success
	Fail    Fail
	Always  Always
}

type Response struct {
	Code int
	Body string
}

type Success func(response *Response)
type Fail func(status int, errMsg string)
type Always func()

func Send(reqData Request) {
	urlTempArr := strings.Split(reqData.Url, "://")
	if len(urlTempArr) > 1 {
		urlStartTemp := urlTempArr[0] + "://"
		requestUrl := strings.Replace(reqData.Url, urlStartTemp, "", 1)
		requestUrl = strings.Replace(requestUrl, "//", "/", -1)
		reqData.Url = urlStartTemp + requestUrl
	}
	//reqData.Url = strings.Replace(reqData.Url, "//", "/", -1)
	var body io.Reader = nil

	if reqData.Data != nil {
		v := url.Values{}
		for key, value := range reqData.Data {
			v.Set(key, value)
		}
		body = ioutil.NopCloser(strings.NewReader(v.Encode()))
	}

	if reqData.Json != nil {
		byteTemp, err := json.Marshal(reqData.Json)
		if err != nil {
			log.E("请求参数JSON化时发生错误:", err)
			return
		}
		body = strings.NewReader(string(byteTemp))
	}

	//利用指定的method,url以及可选的body返回一个新的请求.如果body参数实现了io.Closer接口，Request返回值的Body 字段会被设置为body，并会被Client类型的Do、Post和PostFOrm方法以及Transport.RoundTrip方法关闭。

	var request *http.Request
	var err error
	if reqData.Form != nil {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		for key, val := range reqData.Form {
			_ = writer.WriteField(key, val)
		}
		writer.Close()
		//fmt.Println(body.String())
		//fmt.Println(		writer.Boundary())
		request, err = http.NewRequest(string(reqData.Method), reqData.Url, body)
		request.Header.Set("Content-Type", "multipart/form-data;boundary="+writer.Boundary())
	} else {
		request, err = http.NewRequest(string(reqData.Method), reqData.Url, body)
	}

	if err != nil {
		log.E(err)
		return
	}
	// 必须设定该参数,POST参数才能正常提交
	// request.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	// 设置 Header
	if reqData.Data != nil {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded;")
	}
	//if reqData.Form != nil {
	//}
	if reqData.Json != nil {
		request.Header.Set("Content-Type", "application/application/json;")
	}
	if reqData.Header != nil {
		for key, value := range reqData.Header {
			request.Header.Set(key, value)
		}
	}

	resp, err := client.Do(request) //发送请求

	if err != nil {
		log.E("请求错误 ", err.Error())
		return
		//if reqData.Fail != nil {
		//	reqData.Fail()
		//}
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.E("Body 读取错误 ", err.Error())
		return
	}

	defer resp.Body.Close() //一定要关闭resp.Body
	if resp.StatusCode == 200 {
		if reqData.Success != nil {
			reqData.Success(&Response{
				Code: resp.StatusCode,
				Body: string(content),
			})
		}
	} else {
		if reqData.Fail != nil {
			reqData.Fail(resp.StatusCode, string(content))
		}
	}
	if reqData.Always != nil {
		reqData.Always()
	}

}

func Get(data Request) {

}
