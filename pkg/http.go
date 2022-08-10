package pkg

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	// Demo: pkg.Http.PostJson("http://xxx.com/yy", map[string]string{"a": "1"}).Do()
	Http myHttp
	// 日志记录器
	HttpLogger io.Writer = os.Stdout
	// 默认超时时间
	httpDefaultTimeout = time.Second * 3
)

type myHttp struct{}

func (myHttp) Get(url string, params map[string]string) *httpBuilder {
	var r httpBuilder
	r.Method = http.MethodGet
	if len(params) > 0 {
		r.Url = url + "?" + string(r.BuildQuery(params))
	} else {
		r.Url = url
	}
	return &r
}

func (myHttp) PostJson(url string, bodyData interface{}) *httpBuilder {
	var r httpBuilder
	r.Method = http.MethodPost
	r.Url = url
	r.Header = map[string]string{"Content-Type": "application/json"}
	b, _ := json.Marshal(bodyData)
	r.Body = bytes.NewReader(b)
	r.RequestBodyJson = b
	return &r
}

func (myHttp) PostForm(url string, bodyData map[string]string) *httpBuilder {
	var r httpBuilder
	r.Method = http.MethodPost
	r.Url = url
	r.Header = map[string]string{"Content-Type": "application/x-www-form-urlencoded"}
	bodyDataFormat := r.BuildQuery(bodyData)
	r.Body = bytes.NewReader(bodyDataFormat)
	r.RequestBodyString = string(bodyDataFormat)
	return &r
}

type httpBuilder struct {
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Timeout time.Duration     `json:"timeout"`
	Header  map[string]string `json:"header,omitempty"`
	Body    io.Reader         `json:"-"`
	httpLog
}

type httpLog struct {
	Cost              int64           `json:"cost"`
	Error             error           `json:"error,omitempty"`
	RequestBodyString string          `json:"req_string,omitempty"`
	RequestBodyJson   json.RawMessage `json:"req_json,omitempty"`
	ResponseStatus    int             `json:"res_status"`
	ResponseBody      json.RawMessage `json:"res_body"`
}

func (r *httpBuilder) BuildQuery(params map[string]string) []byte {
	if len(params) == 0 {
		return nil
	}

	var buffer bytes.Buffer
	var i int
	for k := range params {
		buffer.WriteString(k)
		buffer.WriteByte('=')
		buffer.WriteString(params[k])
		if i != len(params)-1 {
			buffer.WriteByte('&')
		}
		i++
	}

	return buffer.Bytes()
}

func (r *httpBuilder) SetTimeout(timeout time.Duration) *httpBuilder {
	r.Timeout = timeout
	return r
}

func (r *httpBuilder) AddHeader(key, val string) *httpBuilder {
	if len(r.Header) == 0 {
		r.Header = make(map[string]string)
	}
	r.Header[key] = val
	return r
}

func (r *httpBuilder) Do() ([]byte, error) {
	var (
		req    *http.Request
		client *http.Client
		resp   *http.Response
		start  = time.Now()
	)

	req, r.Error = http.NewRequest(r.Method, r.Url, r.Body)
	if r.Error != nil {
		r.writeLog(start)
		return nil, r.Error
	}

	for k, v := range r.Header {
		req.Header.Set(k, v)
	}
	if r.Timeout == 0 {
		r.Timeout = httpDefaultTimeout
	}

	// 发送请求
	client = &http.Client{Timeout: r.Timeout}
	resp, r.Error = client.Do(req)
	if r.Error != nil {
		r.writeLog(start)
		return nil, r.Error
	}
	defer resp.Body.Close()

	r.ResponseStatus = resp.StatusCode
	r.ResponseBody, r.Error = ioutil.ReadAll(resp.Body)
	r.writeLog(start)

	return r.ResponseBody, r.Error
}

func (r *httpBuilder) DoTo(to interface{}) error {
	b, err := r.Do()
	if err != nil {
		return err
	}
	if to != nil {
		return json.Unmarshal(b, to)
	}
	return nil
}

func (r *httpBuilder) writeLog(start time.Time) {
	r.Cost = time.Since(start).Milliseconds()
	info, _ := json.Marshal(r)
	HttpLogger.Write(info)
	HttpLogger.Write([]byte("\t\n"))
}
