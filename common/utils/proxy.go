package utils

import (
	"context"
	"net"
	"net/http"
	"net/url"

	"github.com/go-resty/resty/v2"
	"golang.org/x/net/proxy"
)

func ApplyProxyResty(c *resty.Client, proxyURL string) {
	if proxyURL == "" {
		return
	}
	u, err := url.Parse(proxyURL)
	if err != nil || u == nil {
		c.SetProxy(proxyURL)
		return
	}
	if u.Scheme == "socks5" || u.Scheme == "socks5h" {
		var auth *proxy.Auth
		if u.User != nil {
			pw, _ := u.User.Password()
			auth = &proxy.Auth{User: u.User.Username(), Password: pw}
		}
		d, err := proxy.SOCKS5("tcp", u.Host, auth, proxy.Direct)
		if err != nil {
			return
		}
		tr := &http.Transport{}
		tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return d.Dial(network, addr)
		}
		c.SetTransport(tr)
		return
	}
	c.SetProxy(proxyURL)
}
