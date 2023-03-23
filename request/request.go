package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/yezige/tools.liu.app/logx"
)

type RequestObj struct {
	u    string         // url
	p    url.Values     // params
	resp *http.Response // response
	err  error          // error
}

func New(u string) *RequestObj {
	return &RequestObj{
		u: u,
		p: url.Values{},
	}
}

func (r *RequestObj) SetParams(params map[string]interface{}) *RequestObj {
	for key, val := range params {
		switch val := val.(type) {
		case []string:
			for i, v := range val {
				r.p.Add(key+"["+strconv.Itoa(i)+"]", v)
			}
		default:
			r.p.Add(key, val.(string))
		}
	}
	return r
}

func (r *RequestObj) Get() *RequestObj {
	logx.LogError.Debugln("Request:" + r.u + "-GETIN:" + r.p.Encode())
	r.resp, r.err = http.Get(r.u + "?" + r.p.Encode())
	if r.err != nil {
		return r
	}
	return r
}

func (r *RequestObj) Post() *RequestObj {
	paramsJson, _ := json.Marshal(r.p)
	logx.LogError.Debugln("Request:" + r.u + "-POSTIN:" + string(paramsJson))
	r.resp, r.err = http.PostForm(r.u, r.p)
	if r.err != nil {
		return r
	}
	return r
}

func (r *RequestObj) PostJSON() *RequestObj {
	paramsJson, err := json.Marshal(r.p)
	if err != nil {
		r.err = err
		return r
	}
	logx.LogError.Debugln("Request:" + r.u + "-POSTIN:" + string(paramsJson))
	r.resp, r.err = http.Post(r.u, "application/json", bytes.NewBuffer(paramsJson))
	if r.err != nil {
		return r
	}
	return r
}

func (r *RequestObj) GetBody() ([]byte, error) {
	if r.err != nil {
		logx.LogError.Debugln("Request:" + r.u + "-OUT:" + r.err.Error())
		return nil, r.err
	}
	defer r.resp.Body.Close()
	res, err := io.ReadAll(r.resp.Body)
	if err != nil {
		logx.LogError.Debugln("Request:" + r.u + "-OUT:" + err.Error())
		return nil, err
	}
	logx.LogError.Debugln("Request:" + r.u + "-OUT:" + string(res))
	return res, nil
}
