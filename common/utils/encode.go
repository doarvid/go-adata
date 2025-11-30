package utils

import (
	"bytes"
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GB2312ToUTF8(s []byte) (string, error) {
	reader := transform.NewReader(
		strings.NewReader(string(s)),
		simplifiedchinese.HZGB2312.NewDecoder(),
	)

	decoded, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func GBKToUTF8(s []byte) (string, error) {
	reader := transform.NewReader(
		bytes.NewReader(s),
		simplifiedchinese.GBK.NewDecoder(),
	)

	decoded, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}
