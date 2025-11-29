package http

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type ProxyConfig struct {
	Enabled bool
	IP      string
}

var proxyCfg ProxyConfig

func SetProxy(enabled bool, ip string) {
	proxyCfg = ProxyConfig{Enabled: enabled, IP: ip}
}

func NewClient() *resty.Client {
	client := resty.New()
	client.SetTimeout(15 * time.Second)
	if proxyCfg.Enabled && proxyCfg.IP != "" {
		client.SetProxy("http://" + proxyCfg.IP)
	}
	return client
}
