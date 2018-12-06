package info

import "awheel/data"

/*
所有的URL的基础信息汇总
*/

// 传入URL然后获取关于该URL的所有的基础信息，必须有HTTP或者HTTPS头
func BaseInfo(domain, url string, isAlive bool) *data.UrlBaseInfo {

	urlBaseInfo := new(data.UrlBaseInfo)
	urlBaseInfo.Domain = domain
	urlBaseInfo.Url = url

	ipinfo, err := Domain2Ip(domain)
	if err == nil {
		urlBaseInfo.Ip = ipinfo.Ip
		urlBaseInfo.IpList = ipinfo.IpList
		urlBaseInfo.IsCDN = ipinfo.IsCDN
		urlBaseInfo.Country = ipinfo.Country
		urlBaseInfo.Region = ipinfo.Region
		urlBaseInfo.City = ipinfo.City
		urlBaseInfo.Isp = ipinfo.Isp
	}

	if err == nil && !urlBaseInfo.IsCDN {
		portinfos, err := IpPort(urlBaseInfo.Ip)
		if err == nil {
			urlBaseInfo.PortInfos = portinfos
		}
	}

	if isAlive {
		webInfo, err := WebInfo(url)
		if err == nil {
			urlBaseInfo.Title = webInfo.Title
			urlBaseInfo.PoweredBy = webInfo.PoweredBy
			urlBaseInfo.Server = webInfo.Server
			urlBaseInfo.Via = webInfo.Via
			urlBaseInfo.Robots = webInfo.Robots
		}
	}

	whois, err := Whois(domain)
	if err == nil {
		urlBaseInfo.Registrar = whois.Registrar
		urlBaseInfo.Registrant = whois.Registrant
		urlBaseInfo.Emali = whois.Emali
		urlBaseInfo.Phone = whois.Phone
		urlBaseInfo.CreateDate = whois.CreateDate
		urlBaseInfo.ExpireDate = whois.ExpireDate
		urlBaseInfo.DomainServer = whois.DomainServer
		urlBaseInfo.DNS = whois.DNS
		urlBaseInfo.Status = whois.Status
	}

	if isAlive {
		cms := CMSDetect(domain)
		urlBaseInfo.Cms = cms
		waf := WAFDetect(url)
		urlBaseInfo.Waf = waf
	}

	return urlBaseInfo
}
