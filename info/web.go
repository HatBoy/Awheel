package info

import (
	"awheel/base"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"awheel/data"
)

/*
获取Web相关信息，包括包括服务器类型，版本，语言，网站的title，X-Powered-By等
*/

func serverInfo(url string) (string, string, string, error) {
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	response, err := base.Get(url, headers, true)
	if err != nil {
		return "", "", "", err
	}

	respHeaders := response.Headers
	poweredBy, _ := respHeaders["X-Powered-By"]
	server, _ := respHeaders["Server"]

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Body))
	if err != nil {
		return "", poweredBy, server, nil
	}

	title := doc.Find("title").Text()

	return title, poweredBy, server, nil
}

func robots(url string) (string, error) {
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	response, err := base.Get(url+"/robots.txt", headers, false)
	if err == nil && response.StatusCode == 200 {
		body := response.Body
		if strings.Contains(body, "User-agent") && strings.Contains(body, "Disallow") {
			return body, nil
		} else {
			return "", nil
		}
	} else {
		return "", nil
	}
}

func WebInfo(url string) (*data.WebInfoData, error) {
	webInfoData := new(data.WebInfoData)
	title, poweredBy, server, err := serverInfo(url)
	if err != nil {
		return nil, err
	}
	robot, _ := robots(url)
	webInfoData.Url = url
	webInfoData.Title = title
	webInfoData.PoweredBy = poweredBy
	webInfoData.Server = server
	webInfoData.Robots = robot
	return webInfoData, nil
}
