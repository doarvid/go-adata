package indexmarket

import (
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
	"github.com/doarvid/go-adata/common/utils"
	"github.com/go-resty/resty/v2"
)

type HTTPClientConfig struct {
	Timeout   time.Duration
	Proxy     string
	UserAgent string
}

var pkgClient *resty.Client
var pkgCfg = HTTPClientConfig{
	Timeout:   15 * time.Second,
	UserAgent: "go-adata/indexmarket",
}

func SetHTTPClient(c *resty.Client) {
	pkgClient = c
}

func SetHTTPClientConfig(cfg HTTPClientConfig) {
	pkgCfg = cfg
	pkgClient = nil
}

func getHTTPClient() *resty.Client {
	if pkgClient != nil {
		return pkgClient
	}
	c := resty.New()
	c.SetTimeout(pkgCfg.Timeout)
	ua := pkgCfg.UserAgent
	if ua == "" {
		ua = browser.Random()
	}
	c.SetHeader("User-Agent", ua)
	if pkgCfg.Proxy != "" {
		utils.ApplyProxyResty(c, pkgCfg.Proxy)
	}
	pkgClient = c
	return pkgClient
}
