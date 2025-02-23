package curl_test

import (
	"fmt"
	"github.com/tianlin0/go-plat-curl/curl"
	"github.com/tianlin0/go-plat-utils/conf"
	"github.com/tianlin0/go-plat-utils/goroutines"
	"net/http"
	"testing"
	"time"
)

const localUrl = "https://static.json.cn/r/json/search_list.json"

var data = map[string]interface{}{
	"name":    "HttpRequest",
	"version": "v1.0",
}

var defaultClient = curl.NewClient()

func TestGetResponseWithCache(t *testing.T) {
	conf.SetEnv(conf.EnvLoc)
	_ = defaultClient.NewRequest(&curl.Request{
		Url:    localUrl,
		Data:   data,
		Method: http.MethodGet,
		Header: nil,
	}).SetCacheTime(5 * time.Second).SetCacheCheckFunc(func(resp *curl.Response) bool {
		t.Log("check func")
		return true
	}).Submit(nil)

	goroutines.GoAsync(func(params ...any) {
		time.Sleep(2 * time.Second)
		_ = defaultClient.NewRequest(&curl.Request{
			Url:    localUrl,
			Data:   data,
			Method: http.MethodGet,
			Header: nil,
		}).SetCache(5*time.Second, func(resp *curl.Response) bool {
			t.Log("check func")
			return true
		}).Submit(nil)
	})

	time.Sleep(7 * time.Second)

	_ = defaultClient.NewRequest(&curl.Request{
		Url:    localUrl,
		Data:   data,
		Method: http.MethodGet,
		Header: nil,
	}).SetCacheTime(5 * time.Second).Submit(nil)

	time.Sleep(3 * time.Second)

	//t.Log(conv.String(resp))
}

func TestGetResponseWithRetry(t *testing.T) {
	conf.SetEnv(conf.EnvLoc)

	_ = defaultClient.NewRequest(&curl.Request{
		Url:    localUrl,
		Data:   data,
		Method: http.MethodGet,
		Header: nil,
	}).SetCacheTime(10*time.Second).SetRetryPolicy(&curl.RetryPolicy{
		RetryCondFunc: func(resp *curl.Response) error {
			return fmt.Errorf("get error")
		},
		Attempts: 3,
	}).SetRetry(3, func(resp *curl.Response) error {
		return fmt.Errorf("get error")
	}).Submit(nil)

	//t.Log(conv.String(resp))
}
