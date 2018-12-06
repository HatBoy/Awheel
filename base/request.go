package base

import (
	"net/http"
	"io/ioutil"
	"strings"
	"compress/gzip"
	"errors"
	"net"
	"time"
	"sync"
	"crypto/tls"
)

/*
将请求进行包装简化调用过程
*/

// 响应结构体
type Response struct {
	StatusCode int
	Url        string
	Headers    map[string]string
	Body       string
}

func stopRedirect(req *http.Request, via []*http.Request) error {
	return errors.New("stopped redirects")
}

func httpDo(method string, url string, headers map[string]string, data string, redirect bool) (*Response, error) {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := new(http.Client)
	if redirect {
		client = &http.Client{
			Timeout:   30 * time.Second,
			Transport: transport,
		}
	} else {
		client = &http.Client{
			Timeout:       30 * time.Second,
			Transport:     transport,
			CheckRedirect: stopRedirect,
		}
	}

	var reqest *http.Request
	var err error
	if method == "GET" {
		reqest, err = http.NewRequest(method, url, nil)
	} else if method == "POST" {
		reqest, err = http.NewRequest(method, url, strings.NewReader(data))
		reqest.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	} else {
		return nil, errors.New("Error HTTP Method")
	}

	if err != nil {
		return nil, err
	}
	//添加header请求头
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:62.0) Gecko/20100101 Firefox/62.0")

	for h, v := range headers {
		reqest.Header.Add(h, v)
	}

	resp, err := client.Do(reqest) //提交

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 重新组装响应头
	rawHeaders := resp.Header
	newHeaders := make(map[string]string)

	for k, v := range rawHeaders {
		newHeaders[k] = strings.Join(v, "")
	}

	// 根据收到的数据判断是否进行解压缩
	_, ok := newHeaders["Content-Encoding"]

	var bodyString string
	if !ok {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		bodyString = string(body)
	} else {
		compressedReader, err := gzip.NewReader(resp.Body)
		body, err := ioutil.ReadAll(compressedReader)
		if err != nil {
			return nil, err
		}
		bodyString = string(body)
	}

	// 组装Response并返回，目前仅支持UTF-8编码，日后考虑添加GBK编码
	response := new(Response)
	response.StatusCode = resp.StatusCode
	response.Url = url
	response.Headers = newHeaders
	response.Body = bodyString
	return response, nil
}

// GET请求包装
func Get(url string, headers map[string]string, redirect bool) (*Response, error) {
	response, err := httpDo("GET", url, headers, "", redirect)
	return response, err
}

// POST请求包装
func Post(url string, headers map[string]string, data string, redirect bool) (*Response, error) {
	response, err := httpDo("POST", url, headers, data, redirect)
	return response, err
}

var hwg sync.WaitGroup

func getIndex(proto, domain string, ch chan string) {
	defer hwg.Done()
	url := proto + "://" + domain
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	response, err := Get(url, headers, true)
	if err != nil {
		return
	}
	if response.StatusCode == 200 {
		ch <- proto
	}
}

//传入域名通过请求判断是HTTP协议还是HTTPS协议
func HttpOrHttps(domain string) string {
	hwg.Add(2)
	ch := make(chan string, 2)

	go getIndex("http", domain, ch)
	go getIndex("https", domain, ch)

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Second * 30) // 等待30秒钟
		timeout <- true
	}()
	select {
	case proto := <-ch:
		return proto
	case <-timeout:
		return ""
	}

	hwg.Wait()
	close(ch)
	return ""
}

var twg sync.WaitGroup

func testIndex(url string, t string, ch chan string) {
	defer twg.Done()
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	response, err := Get(url, headers, true)
	if err != nil {
		return
	}
	if response.StatusCode == 200 || response.StatusCode == 302 {
		ch <- t
	}
}

//判断Web类型，是ASP还是JSP还是PHP
func TestWebType(url string) string {
	twg.Add(4)
	ch := make(chan string, 4)

	go testIndex(url+"/index.asp", "asp", ch)
	go testIndex(url+"/index.aspx", "aspx", ch)
	go testIndex(url+"/index.jsp", "jsp", ch)
	go testIndex(url+"/index.php", "php", ch)

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Second * 30) // 等待30秒钟
		timeout <- true
	}()
	select {
	case t := <-ch:
		return t
	case <-timeout:
		return ""
	}

	twg.Wait()
	close(ch)
	return ""
}

func WebIsAlive(url string) bool {
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	_, err := Get(url, headers, true)
	if err != nil {
		return false
	} else {
		return true
	}
}
