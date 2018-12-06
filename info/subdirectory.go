package info

import (
	"awheel/base"
)

func SubDirCheck(url string) (bool, int) {
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	response, err := base.Get(url, headers, false)
	if err != nil{
		return false, 0
	}
	statusCode := response.StatusCode
	if statusCode == 200 || statusCode == 403 || statusCode == 302 {
		return true, statusCode
	} else {
		return false, 0
	}
}
