package info

import (
	"awheel/base"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"awheel/data"
)

/*
调用http://whois.chaicp.com/home_whois/cha?ym=qqq.qq.com接口查询Whois信息
*/

// 传入域名查询whois信息
func Whois(domain string) (*data.WhoisInfo, error) {
	url := "http://whois.chaicp.com/home_whois/cha?ym=" + domain

	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8",
		"Accept-Language": "ja,zh-CN;q=0.8,zh;q=0.6",
		"Accept-Encoding": "gzip",
		"Connection":      "Close",
		"User-Agent":      "Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0",
	}

	response, err := base.Get(url, headers, true)
	if err != nil {
		return nil, err
	}

	body := response.Body

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	whoisInfo := new(data.WhoisInfo)

	doc.Find(".whois-list").Each(func(i int, s1 *goquery.Selection) {
		s1.Find("li").Each(func(i int, s2 *goquery.Selection) {
			name := strings.TrimSpace(s2.Find("div").Eq(0).Text())
			value := strings.TrimSpace(s2.Find("div").Eq(1).Text())
			switch name {
			case "域名：":
				whoisInfo.Domain = value
			case "注册商：":
				whoisInfo.Registrar = value
			case "联系人：":
				whoisInfo.Registrant = value
			case "联系邮箱：":
				whoisInfo.Emali = value
			case "联系电话：":
				whoisInfo.Phone = value
			case "创建时间：":
				whoisInfo.CreateDate = value
			case "过期时间：":
				whoisInfo.ExpireDate = value
			case "域名服务器：":
				whoisInfo.DomainServer = value
			case "DNS：":
				var values []string
				s2.Find("p").Each(func(i int, s3 *goquery.Selection) {
					values = append(values, strings.TrimSpace(s3.Text()))
				})
				whoisInfo.DNS = values
			case "状态：":
				var values []string
				s2.Find("p").Each(func(i int, s3 *goquery.Selection) {
					values = append(values, strings.TrimSpace(s3.Text()))
				})
				whoisInfo.Status = values
			}
		})
	})

	return whoisInfo, nil
}
