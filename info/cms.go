package info

import (
	"awheel/base"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"encoding/json"
	"sync"
	"time"
)

/*
调用多个接口查询CMS类型
http://whatweb.bugscaner.com/look/
https://whatcms.org/?s=www.freebuf.com
http://www.iguoli.cn/cms.php
*/

var cmsWG sync.WaitGroup

// http://whatweb.bugscaner.com/look/查询
func bugscanerSearch(domain string, ch chan string) (string, error) {
	defer cmsWG.Done()
	hashUrl := "http://whatweb.bugscaner.com/look/"
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	hashResponse, err := base.Get(hashUrl, headers, true)
	if err != nil {
		return "", err
	}
	hashBody := hashResponse.Body
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(hashBody))
	if err != nil {
		return "", err
	}
	hash, ok := doc.Find("input[type=hidden]").Attr("value")
	if !ok {
		return "", err
	}

	url := "http://whatweb.bugscaner.com/what/"
	data := "url=URL&hash=HASH"
	data = strings.Replace(data, "URL", domain, -1)
	data = strings.Replace(data, "HASH", hash, -1)
	response, err := base.Post(url, headers, data, true)
	if err != nil {
		return "", err
	}

	type resp struct {
		Url   string `json:url`
		Md5   string `json:md5`
		Cms   string `json:cms`
		Error string `json:error`
	}
	var result resp
	err = json.Unmarshal([]byte(response.Body), &result)
	if err != nil {
		return "", err
	}
	cms := result.Cms
	cms = strings.TrimSpace(strings.ToLower(cms))
	ch <- cms
	return cms, nil
}

// https://whatcms.org/查询
func whatcmsSearch(domain string, ch chan string) (string, error) {
	defer cmsWG.Done()
	url := "https://whatcms.org/?s=" + domain
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	response, err := base.Get(url, headers, true)
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Body))
	if err != nil {
		return "", err
	}

	cms := doc.Find("div[class=panel-body]").Eq(1).Find("a").Eq(1).Text()
	cms = strings.TrimSpace(strings.ToLower(cms))
	ch <- cms
	return cms, nil
}

// http://www.iguoli.cn/cms.php查询
func iguoliSearch(domain string, ch chan string) (string, error) {
	defer cmsWG.Done()
	url := "http://www.iguoli.cn/cms.php"
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	response, err := base.Post(url, headers, "url="+domain, true)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(response.Body))
	if err != nil {
		return "", err
	}
	cms := doc.Find("font").Eq(1).Text()
	cms = strings.TrimSpace(strings.ToLower(cms))
	ch <- cms

	return cms, err
}

// CMS类型识别
func CMSDetect(domain string) string {
	cmsWG.Add(3)
	ch := make(chan string, 3)

	go bugscanerSearch(domain, ch)
	go iguoliSearch(domain, ch)
	go whatcmsSearch(domain, ch)

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Second * 30) // 等待30秒钟
		timeout <- true
	}()
	select {
	case cms := <-ch:
		return cms
	case <-timeout:
		return ""
	}

	cmsWG.Wait()
	close(ch)
	return ""
}
