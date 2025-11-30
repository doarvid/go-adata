package header

func SinaCHeaders() map[string]string {
	return map[string]string{
		"Host":            "hq.sinajs.cn",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0",
		"Accept":          "*/*",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
		"Accept-Encoding": "gzip, deflate, br",
		"Referer":         "http://vip.stock.finance.sina.com.cn/",
		"Connection":      "keep-alive",
		"Sec-Fetch-Dest":  "script",
		"Sec-Fetch-Mode":  "no-cors",
		"Sec-Fetch-Site":  "cross-site",
	}
}
