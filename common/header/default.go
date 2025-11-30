package header

import browser "github.com/EDDYCJY/fake-useragent"

func DefaultHeaders() map[string]string {
	return map[string]string{
		"User-Agent":      browser.Random(),
		"Accept":          "*/*",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Accept-Encoding": "gzip, deflate, br",
	}
}
