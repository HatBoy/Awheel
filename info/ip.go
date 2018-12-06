package info

import (
	"awheel/base"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"net/url"
	"regexp"
	"encoding/json"
	"awheel/data"
)

/*
收集域名的IP，判断CDN，IP归属地，机房等信息
*/

// 调用https://ping.aizhan.com/接口查询超级ping结果
func superPing(domain string) ([]string, error) {
	response, err := base.Get("https://ping.aizhan.com/", nil, true)
	if err != nil {
		return nil, err
	}

	body := response.Body
	headers := response.Headers
	cookie := headers["Set-Cookie"]

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	// 获取Cookie和token
	token, _ := doc.Find("meta[name=csrf-token]").Attr("content")

	postData := "type=ping&domain=DOMAIN&_csrf=CSRF"
	postData = strings.Replace(postData, "DOMAIN", domain, -1)
	postData = strings.Replace(postData, "CSRF", url.QueryEscape(token), -1)

	postHeaders := map[string]string{
		"X-Requested-With": "XMLHttpRequest",
		"Content-Type":     "application/x-www-form-urlencoded; charset=UTF-8",
		"Cookie":           cookie,
	}

	// 获取结果数据
	postResponse, err := base.Post("https://ping.aizhan.com/api/ping?callback=flightHandler", postHeaders, postData, true)

	if err != nil {
		return nil, err
	}

	responseBody := postResponse.Body
	reg := regexp.MustCompile(`(?U)"ip":".*",`)
	ips := reg.FindAllString(responseBody, -1)
	var ipList []string

	for _, ip := range ips {
		ip = ip[6 : len(ip)-2]
		check := true
		for _, i := range ipList {
			if ip == i {
				check = false
			}
		}

		if check && ip != "-" {
			ipList = append(ipList, ip)
		}

	}
	return ipList, nil
}

type DataS struct {
	Ip         string `json:ip`
	Country    string `json:country`
	Area       string `json:area`
	Region     string `json:region`
	City       string `json:city`
	County     string `json:county`
	Isp        string `json:isp`
	Country_id string `json:country_id`
	Areaid     string `json:area_id`
	Regionid   string `json:region_id`
	Cityid     string `json:city_id`
	Countyid   string `json:county_id`
	Ispid      string `json:isp_id`
}

type TaobaoData struct {
	Code int   `json:code`
	Data DataS `json:data`
}

// 收集IP归属地信息，http://ip.taobao.com/service/getIpInfo.php?ip=XXX
func ipSearch(ip string) (*data.IpInfo, error) {
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "gzip, deflate, br",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	url := "http://ip.taobao.com/service/getIpInfo.php?ip=" + ip
	response, err := base.Get(url, headers, true)

	if err != nil {
		return nil, err
	}
	body := response.Body

	var taobaoData TaobaoData
	err = json.Unmarshal([]byte(body), &taobaoData)
	if err != nil {
		return nil, err
	}
	ipInfo := new(data.IpInfo)
	ipInfo.Ip = ip
	ipInfo.Country = taobaoData.Data.Country
	ipInfo.Region = taobaoData.Data.Region
	ipInfo.City = taobaoData.Data.City
	ipInfo.Isp = taobaoData.Data.Isp

	return ipInfo, nil
}

// 调用该函数获得将域名转为IP的全部数据
func Domain2Ip(domain string) (*data.IpInfo, error) {
	// 220.250.64.20神奇的IP
	//ips1, err := net.LookupIP(domain)
	resolver := base.New([]string{"223.5.5.5", "114.114.114.114", "119.29.29.29", "182.254.116.116", "180.76.76.76"})
	resolver.RetryTimes = 10

	ips1, err := resolver.LookupHost(domain)

	if err != nil {
		return nil, err
	}

	var newIps1 []string
	for _, ip := range ips1 {
		newIps1 = append(newIps1, ip.String())
	}

	ipInfo := new(data.IpInfo)
	ipInfo.Domain = domain

	ips2, err := superPing(domain)
	if err != nil {
		return nil, err
	}
	var ipList []string
	for _, ip := range ips2 {
		newIps1 = append(newIps1, ip)
	}

	for _, ip := range newIps1 {
		check := true
		for _, i := range ipList {
			if ip == i {
				check = false
			}
		}

		if check {
			ipList = append(ipList, ip)
		}
	}

	ipInfo.IpList = ipList

	if len(ipList) == 1 {
		ipInfo.Ip = ipList[0]
		ipInfo.IsCDN = false
		ipMore, err := ipSearch(ipList[0])
		if err != nil {
			return ipInfo, nil
		}
		ipInfo.Country = ipMore.Country
		ipInfo.Region = ipMore.Region
		ipInfo.City = ipMore.City
		ipInfo.Isp = ipMore.Isp
		return ipInfo, nil
	} else {
		ipInfo.IsCDN = true
		return ipInfo, nil
	}
}
