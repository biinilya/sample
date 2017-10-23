package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/astaxie/beego"
)

func FromJson(data []byte, v interface{}) interface{} {
	if jsErr := json.Unmarshal(data, v); jsErr != nil {
		panic(jsErr)
	}
	return v
}

func ToJson(v interface{}) []byte {
	if data, jsErr := json.Marshal(v); jsErr != nil {
		panic(jsErr)
	} else {
		return data
	}
}

func BeegoCall(method, url string, body io.Reader) *httptest.ResponseRecorder {
	return BeegoCallWithHeader(method, url, body, make(http.Header))
}

func BeegoCallWithHeader(method, url string, body io.Reader, header http.Header) *httptest.ResponseRecorder {
	r, _ := http.NewRequest(method, url, body)
	for k, v := range header {
		r.Header.Del(k)
		for _, vi := range v {
			r.Header.Add(k, vi)
		}
	}
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	return w
}
