package curl

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	ContentTypeJSON           = "application/json"
	ContentTypeFormData       = "application/form-data"
	ContentTypeFormUrlencoded = "application/x-www-form-urlencoded"
)

var (
	Logger io.Writer = os.Stdout
)

type request struct {
	Time              string            `json:"time"`
	Cost              int64             `json:"cost"`
	Error             error             `json:"error,omitempty"`
	Url               string            `json:"url"`
	Method            string            `json:"method"`
	Header            map[string]string `json:"header,omitempty"`
	RequestBody       io.Reader         `json:"-"`
	RequestBodyString string            `json:"request_string,omitempty"`
	RequestBodyJson   json.RawMessage   `json:"request_json,omitempty"`
	ResponseBody      json.RawMessage   `json:"response"`
}

func Get(url string, params map[string]string) *request {
	if len(params) > 0 {
		url = url + "?" + string(BuildQuery(params))
	}
	var r request
	r.Method = http.MethodGet
	r.Url = url
	return &r
}

func PostJson(url string, bodyData interface{}) *request {
	b, _ := json.Marshal(bodyData)
	var r request
	r.Method = http.MethodPost
	r.Url = url
	r.Header = map[string]string{"Content-Type": ContentTypeJSON}
	r.RequestBody = bytes.NewReader(b)
	r.RequestBodyJson = b
	return &r
}

func PostForm(url string, bodyData map[string]string) *request {
	bodyDataFormat := BuildQuery(bodyData)
	var r request
	r.Method = http.MethodPost
	r.Url = url
	r.Header = map[string]string{"Content-Type": ContentTypeFormUrlencoded}
	r.RequestBody = bytes.NewReader(bodyDataFormat)
	r.RequestBodyString = string(bodyDataFormat)
	return &r
}

func BuildQuery(params map[string]string) []byte {
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

func (r *request) AddHeader(key, val string) *request {
	if len(r.Header) == 0 {
		r.Header = make(map[string]string)
	}
	r.Header[key] = val
	return r
}

func (r *request) Do() ([]byte, error) {
	start := time.Now()
	r.Time = start.Format("2006-01-02 15:04:05")

	req, err := http.NewRequest(r.Method, r.Url, r.RequestBody)
	if err == nil {
		for k, v := range r.Header {
			req.Header.Set(k, v)
		}
		var resp *http.Response
		resp, err = http.DefaultClient.Do(req)

		if err == nil {
			r.ResponseBody, err = ioutil.ReadAll(resp.Body)
			if err == nil {
				resp.Body.Close()
			}
		}
	}
	if err != nil {
		r.Error = err
		logrus.Error(err)
	}

	r.Cost = time.Since(start).Milliseconds()
	info, _ := json.Marshal(r)
	Logger.Write(info)
	Logger.Write([]byte("\t\n"))

	return r.ResponseBody, err
}

func (r *request) DoTo(to interface{}) error {
	b, err := r.Do()
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, to)
	return err
}
