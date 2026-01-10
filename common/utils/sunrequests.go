package utils

import (
    "strings"
    "time"

    "github.com/go-resty/resty/v2"
)

type SunProxy struct{}

var sunProxyData = map[string]string{}

func (SunProxy) Set(key, value string) { sunProxyData[key] = value }
func (SunProxy) Get(key string) string { return sunProxyData[key] }
func (SunProxy) Delete(key string) { delete(sunProxyData, key) }

type SunRequests struct{ proxy SunProxy }

func NewSunRequests(p SunProxy) *SunRequests { return &SunRequests{proxy: p} }

func (s *SunRequests) Request(method, url string, times int, retryWaitTimeMs int, proxies map[string]string, waitTimeMs int, setup func(*resty.Request)) (*resty.Response, error) {
    if times <= 0 { times = 1 }
    if retryWaitTimeMs <= 0 { retryWaitTimeMs = 1588 }
    client := resty.New()
    isProxy := strings.ToLower(s.proxy.Get("is_proxy")) == "true"
    ip := s.proxy.Get("ip")
    proxyURL := s.proxy.Get("proxy_url")
    if isProxy && ip == "" && proxyURL != "" {
        client.SetRetryCount(1)
        client.SetRetryWaitTime(500 * time.Millisecond)
    }
    if isProxy && ip != "" {
        client.SetProxy("http://" + ip)
    }
    var lastResp *resty.Response
    var lastErr error
    for i := 0; i < times; i++ {
        if waitTimeMs > 0 { time.Sleep(time.Duration(waitTimeMs) * time.Millisecond) }
        req := client.R()
        if setup != nil { setup(req) }
        switch strings.ToLower(method) {
        case "post":
            lastResp, lastErr = req.Post(url)
        default:
            lastResp, lastErr = req.Get(url)
        }
        if lastErr == nil && lastResp != nil && (lastResp.StatusCode() == 200 || lastResp.StatusCode() == 404) {
            return lastResp, nil
        }
        time.Sleep(time.Duration(retryWaitTimeMs) * time.Millisecond)
    }
    return lastResp, lastErr
}

var SunRequestsClient = NewSunRequests(SunProxy{})
