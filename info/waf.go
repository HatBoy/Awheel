package info

import (
	"awheel/base"
	"strings"
	"regexp"
	"sync"
	"time"
	"net/url"
)

/*
WAF检测，参考下面两个开源程序
https://github.com/Ekultek/WhatWaf
https://github.com/EnableSecurity/wafw00f
*/

type detectFunc func(*base.Response) (string, bool)

func qihu360(response *base.Response) (string, bool) {
	signature1 := "X-Powered-By-360wzb"
	signature2 := "360wzws"
	headers := response.Headers
	_, ok := headers[signature1]
	if ok {
		return "360网站卫士", true
	} else {
		server, _ := headers["Server"]
		check := strings.Contains(server, signature2)
		if check {
			return "360网站卫士", true
		} else {
			return "", false
		}
	}
}

func airlock(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`^AL[_-]?(SESS|LB)=`)
	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	check := signature1.FindAllString(cookie, -1)
	if check != nil {
		return "InfoGuard Airlock", true
	} else {
		return "", false
	}
}

func akamaiGHost(response *base.Response) (string, bool) {
	signature1 := "AkamaiGHost"
	headers := response.Headers
	server, _ := headers["Server"]
	check := strings.Contains(server, signature1)
	if check {
		return "AkamaiGHost Website Protection", true
	} else {
		return "", false
	}
}

func anquanbao(response *base.Response) (string, bool) {
	signature1 := "X-Powered-By-Anquanbao"
	headers := response.Headers
	_, ok := headers[signature1]
	if ok {
		return "安全宝", true
	} else {
		return "", false
	}
}

func anyu(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)sorry(.)?.your.access.has.been.intercept(ed)?.by.anyu`)
	body := response.Body
	check := signature1.FindAllString(body, -1)
	if check != nil {
		return "AnYu WAF", true
	} else {
		return "", false
	}
}

func armor(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)blocked.by.website.protection.from.armour`)
	body := response.Body
	check := signature1.FindAllString(body, -1)
	if check != nil {
		return "Armor Protection", true
	} else {
		return "", false
	}
}

func asm(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)the.requested.url.was.rejected..please.consult.with.your.administrator`)
	body := response.Body
	check := signature1.FindAllString(body, -1)
	if check != nil {
		return "Application Security Manager (F5 Networks)", true
	} else {
		return "", false
	}
}

func aws(response *base.Response) (string, bool) {
	signature1 := "AWS"
	signature2 := "aws"
	headers := response.Headers
	poweredBy, _ := headers["X-Powered-By"]
	check1 := strings.Contains(poweredBy, signature1)
	check2 := strings.Contains(poweredBy, signature2)
	if check1 || check2 {
		return "AWS WAF", true
	} else {
		return "", false
	}
}

func baidu(response *base.Response) (string, bool) {
	signature1 := "yunjiasu"
	headers := response.Headers
	server, _ := headers["Server"]
	check := strings.Contains(server, signature1)
	if check {
		return "百度云加速", true
	} else {
		return "", false
	}
}

func barracuda(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`^barra_counter_session=`)
	signature2 := regexp.MustCompile(`^BNI__BARRACUDA_LB_COOKIE=`)
	signature3 := regexp.MustCompile(`^BNI_persistence=`)
	signature4 := regexp.MustCompile(`^BN[IE]S_.*?=`)
	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	check1 := signature1.FindAllString(cookie, -1)
	check2 := signature2.FindAllString(cookie, -1)
	check3 := signature3.FindAllString(cookie, -1)
	check4 := signature4.FindAllString(cookie, -1)
	if (check1 != nil) || (check2 != nil) || (check3 != nil) || (check4 != nil) {
		return "Barracuda WAF", true
	} else {
		return "", false
	}
}

func betterwpsecurity(response *base.Response) (string, bool) {
	signature1 := "https://api.w.org/"
	headers := response.Headers
	link, _ := headers["Link"]
	check := strings.Contains(link, signature1)
	if check {
		return "Better WP Security", true
	} else {
		return "", false
	}
}

func f5Bigip(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)biG]gip|bipserver`)
	signature2 := regexp.MustCompile(`(?i)^TS[a-zA-Z0-9]{3,8}=`)
	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	server, _ := headers["Server"]
	check1 := signature1.FindAllString(server, -1)
	check2 := signature2.FindAllString(cookie, -1)
	if (check1 != nil) || (check2 != nil) {
		return "BIG-IP ASM (F5 Networks)", true
	} else {
		return "", false
	}
}

func binarySEC(response *base.Response) (string, bool) {
	signature1 := "X-BinarySEC-Via"
	signature2 := "X-BinarySEC-nocache"
	signature3 := regexp.MustCompile(`(?i).*BinarySec.*`)
	headers := response.Headers
	server, _ := headers["Server"]
	_, check1 := headers[signature1]
	_, check2 := headers[signature2]
	check3 := signature3.FindAllString(server, -1)
	if check1 || check2 || (check3 != nil) {
		return "BinarySEC WAF", true
	} else {
		return "", false
	}
}

func blockDoS(response *base.Response) (string, bool) {
	signature1 := "BlockDos"
	headers := response.Headers
	server, _ := headers["Server"]
	check := strings.Contains(server, signature1)
	if check {
		return "BlockDos", true
	} else {
		return "", false
	}
}

func chinaCache(response *base.Response) (string, bool) {
	signature1 := "Powered-By-ChinaCache"
	headers := response.Headers
	_, check := headers[signature1]
	if check {
		return "ChinaCache-CDN", true
	} else {
		return "", false
	}
}

func ciscoAceXml(response *base.Response) (string, bool) {
	signature1 := "ACE XML Gateway"
	headers := response.Headers
	server, _ := headers["Server"]
	check := strings.Contains(server, signature1)
	if check {
		return "Cisco ACE XML Gateway", true
	} else {
		return "", false
	}
}

func cloudFlare(response *base.Response) (string, bool) {
	signature1 := "cloudflare-nginx"
	signature2 := "__cfduid"
	headers := response.Headers
	server, _ := headers["Server"]
	cookie, _ := headers["Set-Cookie"]
	check1 := strings.Contains(server, signature1)
	check2 := strings.Contains(cookie, signature2)
	if check1 || check2 {
		return "CloudFlare WAF", true
	} else {
		return "", false
	}
}

func cloudFront(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)\d.\d.[a-zA-Z0-9]{32,60}.cloudfront.net`)
	signature2 := regexp.MustCompile(`(?i)cloudfront`)
	signature3 := regexp.MustCompile(`(?i)x.amz.cf.id`)
	headers := response.Headers
	for _, value := range headers {
		check1 := signature1.FindAllString(value, -1)
		check2 := signature2.FindAllString(value, -1)
		check3 := signature3.FindAllString(value, -1)
		if (check1 != nil) || (check2 != nil) || (check3 != nil) {
			return "CloudFront Firewall (Amazon)", true
		}
	}
	return "", false
}

func codeIgniter(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)the.uri.you.submitted.has.disallowed.characters`)
	body := response.Body
	check := signature1.FindAllString(body, -1)
	if check != nil {
		return "XSS/CSRF Filtering Protection (CodeIgniter)", true
	} else {
		return "", false
	}
}

func comodo(response *base.Response) (string, bool) {
	signature1 := "Protected by COMODO WAF"
	headers := response.Headers
	server, _ := headers["Server"]
	check := strings.Contains(server, signature1)
	if check {
		return "Comodo WAF", true
	} else {
		return "", false
	}
}

func ibmDataPower(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)^(OK|FAIL)`)
	headers := response.Headers
	xBacksideTransport, _ := headers["X-Backside-Transport"]
	check := signature1.FindAllString(xBacksideTransport, -1)
	if check != nil {
		return "IBM Websphere DataPower Firewall", true
	} else {
		return "", false
	}

}

func denyAll(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)\Acondition.intercepted`)
	signature2 := regexp.MustCompile(`(?i)\Asessioncookie=`)
	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	body := response.Body
	check1 := signature1.FindAllString(body, -1)
	check2 := signature2.FindAllString(cookie, -1)
	if (check1 != nil) || (check2 != nil) {
		return "Deny All WAF", true
	} else {
		return "", false
	}

}

func dod(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)dod.enterprise.level.protection.system`)
	body := response.Body
	check1 := signature1.FindAllString(body, -1)
	if check1 != nil {
		return "DoD Enterprise-Level Protection System (Department of Defense)", true
	} else {
		return "", false
	}

}

func dOSarrest(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i).*dosarrest.*`)
	headers := response.Headers
	server, _ := headers["Server"]
	check1 := signature1.FindAllString(server, -1)
	if check1 != nil {
		return "DOSarrest", true
	} else {
		return "", false
	}
}

func dotDefender(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)dotdefender.blocked.your.request`)
	signature2 := "X-dotDefender-denied"
	headers := response.Headers
	check1 := signature1.FindAllString(response.Body, -1)
	_, check2 := headers[signature2]
	if check1 != nil || check2 {
		return "dotDefender", true
	} else {
		return "", false
	}
}

func dynamicWeb(response *base.Response) (string, bool) {
	signature1 := "X-403-status-by"
	headers := response.Headers
	_, check1 := headers[signature1]
	if check1 {
		return "DynamicWeb Injection Check", true
	} else {
		return "", false
	}
}

func edgecast(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)^ECD \(.*?\)$`)
	signature2 := regexp.MustCompile(`(?i)^ECS \(.*?\)$`)
	headers := response.Headers
	server, _ := headers["Server"]
	check1 := signature1.FindAllString(server, -1)
	check2 := signature2.FindAllString(server, -1)
	if (check1 != nil) || (check2 != nil) {
		return "EdgeCast WAF", true
	} else {
		return "", false
	}
}

func expressionEngine(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)<.+>error.-.expressionengine<.+.>`)
	signature2 := regexp.MustCompile(`(?i)<.+><.+>error<.+.>:.the.uri.you.submitted.has.disallowed.characters.<.+.>`)
	signature3 := regexp.MustCompile(`(?i)invalid.get.data`)
	check1 := signature1.FindAllString(response.Body, -1)
	check2 := signature2.FindAllString(response.Body, -1)
	check3 := signature3.FindAllString(response.Body, -1)
	if (check1 != nil) || (check2 != nil) || (check3 != nil) {
		return "ExpressionEngine (Ellislab WAF)", true
	} else {
		return "", false
	}
}

func f5BigIpAPM(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)^MRHSession`)
	signature2 := regexp.MustCompile(`(?i)^F5_fullWT`)
	signature3 := regexp.MustCompile(`(?i)^F5_ST`)
	signature4 := regexp.MustCompile(`(?i)^F5_HT_shrinked`)
	signature5 := regexp.MustCompile(`(?i)^MRHSequence`)
	signature6 := regexp.MustCompile(`(?i)^MRHSHint`)
	signature7 := regexp.MustCompile(`(?i)^LastMRH_Session`)
	signature8 := regexp.MustCompile(`(?i)BigIP|BIG-IP|BIGIP`)

	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	server, _ := headers["Server"]
	check1 := signature1.FindAllString(cookie, -1)
	check2 := signature2.FindAllString(cookie, -1)
	check3 := signature3.FindAllString(cookie, -1)
	check4 := signature4.FindAllString(cookie, -1)
	check5 := signature5.FindAllString(cookie, -1)
	check6 := signature6.FindAllString(cookie, -1)
	check7 := signature7.FindAllString(cookie, -1)
	check8 := signature8.FindAllString(server, -1)
	if (check1 != nil) || (check2 != nil) || (check3 != nil) || (check4 != nil) || (check5 != nil) || (check6 != nil) || (check7 != nil) || (check8 != nil) {
		return "F5 BIG-IP APM", true
	} else {
		return "", false
	}
}

func f5BigIpLTM(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)^BIGipServer`)
	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	check1 := signature1.FindAllString(cookie, -1)
	if check1 != nil {
		return "F5 BIG-IP LTM", true
	} else {
		return "", false
	}
}

func f5FirePass(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)^MRHSession`)
	signature2 := regexp.MustCompile(`(?i)^uRoamTestCookie`)
	signature3 := regexp.MustCompile(`(?i)^MRHCId`)
	signature4 := regexp.MustCompile(`(?i)^MRHIntranetSession`)

	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	check1 := signature1.FindAllString(cookie, -1)
	check2 := signature2.FindAllString(cookie, -1)
	check3 := signature3.FindAllString(cookie, -1)
	check4 := signature4.FindAllString(cookie, -1)

	if (check1 != nil) || (check2 != nil) || (check3 != nil) || (check4 != nil) {
		return "F5 FirePass", true
	} else {
		return "", false
	}
}

func f5Trafficshield(response *base.Response) (string, bool) {
	signature1 := "F5-TrafficShield"
	headers := response.Headers
	server, _ := headers["Server"]
	check := strings.Contains(server, signature1)
	if check {
		return "F5 Trafficshield", true
	} else {
		return "", false
	}
}

func fortiWeb(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)<.+>powered.by.fortinet<.+.>`),
		regexp.MustCompile(`(?i)<.+>fortigate.ips.sensor<.+.>`),
		regexp.MustCompile(`(?i)fortigate`),
		regexp.MustCompile(`(?i).fgd_icon`),
		regexp.MustCompile(`(?i)FORTIWAFSID=`),
		regexp.MustCompile(`(?i)application.blocked.`),
		regexp.MustCompile(`(?i).fortiGate.application.control`),
		regexp.MustCompile(`(?i)(http(s)?)?://\w+.fortinet(.\w+:)?`),
		regexp.MustCompile(`(?i)fortigate.hostname`)}

	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(cookie, -1)
		check2 := s.FindAllString(body, -1)
		if (check1 != nil) || (check2 != nil) {
			return "FortiWeb WAF", true
		}
	}
	return "", false
}

func gladius(response *base.Response) (string, bool) {
	signature1 := "gladius_blockchain_driven_cyber_protection_network_session"
	headers := response.Headers
	_, check := headers[signature1]
	if check {
		return "Gladius network WAF", true
	} else {
		return "", false
	}
}

func hyperGuard(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)^WODSESSION=`)
	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	check1 := signature1.FindAllString(cookie, -1)
	if check1 != nil {
		return "Art of Defence HyperGuard", true
	} else {
		return "", false
	}
}

func incapsula(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)incap_ses|visid_incap|incapsula`),
		regexp.MustCompile(`(?i)incapsula.incident.id`),
		regexp.MustCompile(`(?i)^visid.*=`)}

	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(cookie, -1)
		check2 := s.FindAllString(body, -1)
		if (check1 != nil) || (check2 != nil) {
			return "Incapsula WAF", true
		}
	}
	return "", false
}

func mission(response *base.Response) (string, bool) {
	signature1 := "Mission Control Application Shield"
	headers := response.Headers
	server, _ := headers["Server"]
	check := strings.Contains(server, signature1)
	if check {
		return "Mission Control Application Shield", true
	} else {
		return "", false
	}
}

func opModsecurity(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)ModSecurity|NYOB`),
		regexp.MustCompile(`(?i)mod.security`),
		regexp.MustCompile(`(?i)This.error.was.generated.by.mod.security`),
		regexp.MustCompile(`(?i)web.server at`),
		regexp.MustCompile(`(?i)page.you.are.(accessing|trying)?.(to|is)?.(access)?.(is|to)?.(restricted)?`)}

	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		if check1 != nil {
			return "Open Source WAF (Modsecurity)", true
		}
	}
	return "", false
}

func trustwave(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)(mod_security|Mod_Security|NOYB)`)
	headers := response.Headers
	server, _ := headers["Server"]
	check1 := signature1.FindAllString(server, -1)
	if check1 != nil {
		return "Trustwave ModSecurity", true
	} else {
		return "", false
	}
}

func owaspModsecurity(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)not.acceptable`),
		regexp.MustCompile(`(?i)additionally\S.a.406.not.acceptable`),
	}

	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		if check1 != nil {
			return "Mod Security (OWASP CSR)", true
		}
	}
	return "", false
}

func knownsec(response *base.Response) (string, bool) {
	signature1 := "当前访问疑似黑客攻击，已被网站管理员设置为拦截"
	signature2 := "请登录知道创宇云安全"
	body := response.Body
	check1 := strings.Contains(body, signature1)
	check2 := strings.Contains(body, signature2)
	if check1 && check2 {
		return "知道创宇云安全", true
	} else {
		return "", false
	}
}

func naxsi(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)^naxsi`)
	headers := response.Headers
	server, _ := headers["X-Data-Origin"]
	check1 := signature1.FindAllString(server, -1)
	if check1 != nil {
		return "Naxsi", true
	} else {
		return "", false
	}
}

func netContinuum(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)^NCI__SessionId=`)
	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	check1 := signature1.FindAllString(cookie, -1)
	if check1 != nil {
		return "NetContinuum", true
	} else {
		return "", false
	}
}

func netScaler(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)^(ns_af=|citrix_ns_id|NSC_)`)
	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	check1 := signature1.FindAllString(cookie, -1)
	if check1 != nil {
		return "Citrix NetScaler", true
	} else {
		return "", false
	}
}

func nevisProxy(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)^Navajo.*?$`)
	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	check1 := signature1.FindAllString(cookie, -1)
	if check1 != nil {
		return "AdNovum nevisProxy", true
	} else {
		return "", false
	}
}

func nsFocus(response *base.Response) (string, bool) {
	signature1 := "NSFocus"
	headers := response.Headers
	server, _ := headers["Server"]
	check := strings.Contains(server, signature1)
	if check {
		return "NSFocus", true
	} else {
		return "", false
	}
}

func paloAlto(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)\bhas.been.blocked.in.accordance.with.company.policy\b`),
		regexp.MustCompile(`(?i)<.+>Virus.Spyware.Download.Blocked<.+.>`),
	}

	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		if check1 != nil {
			return "Palo Alto Firewall", true
		}
	}
	return "", false
}

func perimeterX(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)access.to.this.page.has.been.denied.because.we.believe.you.are.using.automation.tool`),
		regexp.MustCompile(`(?i)http(s)?://(www.)?perimeterx.\w+.whywasiblocked`),
		regexp.MustCompile(`(?i)(..)?client.perimeterx.*/[a-zA-Z]{8,15}/*.*.js`),
	}

	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		if check1 != nil {
			return "Anti Bot Protection (PerimeterX)", true
		}
	}
	return "", false
}

func pkSecurityModule(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)<.+>pkSecurityModule\W..\WSecurity.Alert<.+.>`),
		regexp.MustCompile(`(?i)<.+http(s)?.//([w]{3})?.kitnetwork.\w+.+>`),
		regexp.MustCompile(`(?i)<.+>A.safety.critical.request.was.discovered.and.blocked.<.+.>`),
	}

	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		if check1 != nil {
			return "pkSecurityModule IDS", true
		}
	}
	return "", false
}

func powerCDN(response *base.Response) (string, bool) {
	signature1 := "PowerCDN"
	headers := response.Headers
	_, check := headers[signature1]
	if check {
		return "NSFocus", true
	} else {
		return "", false
	}
}

func powerful(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)Powerful Firewall`),
		regexp.MustCompile(`(?i)http(s)?...tiny.cc.powerful.firewall`),
	}

	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		if check1 != nil {
			return "Powerful Firewall (MyBB plugin)", true
		}
	}
	return "", false
}

func profense(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)profense`)
	headers := response.Headers
	server, _ := headers["Server"]
	check := signature1.FindAllString(server, -1)
	if check != nil {
		return "Profense", true
	} else {
		return "", false
	}
}

func radware(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i).\bcloudwebsec.radware.com\b.`),
		regexp.MustCompile(`(?i)<.+>unauthorized.activity.has.been.detected<.+.>`),
		regexp.MustCompile(`(?i)with.the.following.case.number.in.its.subject:.\d+.`),
	}
	signature2 := "X-SL-CompState"
	headers := response.Headers
	_, check2 := headers[signature2]
	if check2 {
		return "Radware (AppWall WAF)", true
	}
	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		if check1 != nil {
			return "Radware (AppWall WAF)", true
		}
	}
	return "", false
}

func sabre(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)dxsupport@sabre.com`)
	check := signature1.FindAllString(response.Body, -1)
	if check != nil {
		return "Sabre Firewall (WAF)", true
	} else {
		return "", false
	}
}

func safeDog(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)(http(s)?)?(://)?(www|404|bbs|\w+)?.safedog.\w+`),
		regexp.MustCompile(`(?i)waf(.?\d+(.)?\d+)`),
		regexp.MustCompile(`(?i)^Safedog`),
	}
	headers := response.Headers
	server, _ := headers["Server"]
	xPowerBy, _ := headers["X-Powered-By"]
	for _, s := range signatures {
		check1 := s.FindAllString(server, -1)
		check2 := s.FindAllString(xPowerBy, -1)
		if check1 != nil || check2 != nil {
			return "Safedog", true
		}
	}
	return "", false
}

func siteGuard(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)>Powered.by.SiteGuard.Lite<`),
		regexp.MustCompile(`(?i)refuse.to.browse`),
	}
	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		if check1 != nil {
			return "Website Security SiteGuard (Lite)", true
		}
	}
	return "", false
}

func sonicWALL(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)This.request.is.blocked.by.the.SonicWALL`),
		regexp.MustCompile(`(?i)Dell.SonicWALL`),
		regexp.MustCompile(`(?i)\bDell\b`),
		regexp.MustCompile(`(?i)Web.Site.Blocked.+\bnsa.banner`),
		regexp.MustCompile(`(?i)SonicWALL`),
		regexp.MustCompile(`(?i)<.+>policy.this.site.is.blocked<.+.>`),
	}
	headers := response.Headers
	server, _ := headers["Server"]
	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		check2 := s.FindAllString(server, -1)
		if check1 != nil || check2 != nil {
			return "SonicWALL Firewall (Dell)", true
		}
	}
	return "", false
}

func squid(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)squid`),
		regexp.MustCompile(`(?i)Access control configuration prevents`),
		regexp.MustCompile(`(?i)X.Squid.Error`),
	}
	headers := response.Headers
	server, _ := headers["Server"]
	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		check2 := s.FindAllString(server, -1)
		if check1 != nil || check2 != nil {
			return "Squid Proxy IDS", true
		}
	}
	return "", false
}

func stingray(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)\AX-Mapping-`)
	headers := response.Headers
	cookie, _ := headers["Set-Cookie"]
	check := signature1.FindAllString(cookie, -1)
	if check != nil {
		return "Stingray Application Firewall (Riverbed / Brocade)", true
	} else {
		return "", false
	}
}

func sucuri(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)Access Denied - Sucuri Website Firewall`),
		regexp.MustCompile(`(?i)Sucuri WebSite Firewall - CloudProxy - Access Denied`),
		regexp.MustCompile(`(?i)Questions\?.+cloudproxy@sucuri\.net`),
	}
	signature2 := "X-Sucuri-ID"
	headers := response.Headers
	_, check2 := headers[signature2]
	if check2 {
		return "Sucuri WAF", true
	}
	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		if check1 != nil {
			return "Radware (AppWall WAF)", true
		}
	}
	return "", false
}

func teros(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)st8(id|.wa|.wf)?(.)?(\d+|\w+)?`)
	headers := response.Headers

	for _, value := range headers {
		check := signature1.FindAllString(value, -1)
		if check != nil {
			return "Teros WAF", true
		}
	}
	return "", false

}

func urlScan(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)rejected.by.url.scan`),
		regexp.MustCompile(`(?i)/rejected.by.url.scan`),
	}
	headers := response.Headers
	location, _ := headers["Location"]
	for _, s := range signatures {
		check1 := s.FindAllString(location, -1)
		if check1 != nil {
			return "UrlScan (Microsoft)", true
		}
	}
	return "", false
}

func uspses(response *base.Response) (string, bool) {
	signature1 := "Secure Entry Server"
	headers := response.Headers
	server, _ := headers["Server"]
	check := strings.Contains(server, signature1)
	if check {
		return "USP Secure Entry Server", true
	} else {
		return "", false
	}
}

func varnish(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)<.+>(.)?security.by.cachewall(.)?<.+.>`)
	signatures := []string{"X-Varnish", "X-Cachewall-Action", "X-Cachewall-Reason"}
	headers := response.Headers
	for _, s := range signatures {
		_, check1 := headers[s]
		if check1 {
			return "Varnish/CacheWall WAF", true
		}
	}

	check := signature1.FindAllString(response.Body, -1)
	if check != nil {
		return "Varnish/CacheWall WAF", true
	}
	return "", false
}

func wallarm(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)nginix.wallarm`)
	headers := response.Headers
	server, _ := headers["Server"]
	check := signature1.FindAllString(server, -1)
	if check != nil {
		return "Wallarm WAF", true
	}
	return "", false
}

func webknight(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)WebKnight`)
	headers := response.Headers
	server, _ := headers["Server"]
	check := signature1.FindAllString(server, -1)
	if check != nil {
		return "WebKnight Application Firewall (AQTRONIX)", true
	}
	return "", false
}

func webSEAL(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)WebSEAL`)
	check := signature1.FindAllString(response.Body, -1)
	if check != nil {
		return "IBM Security Access Manager (WebSEAL)", true
	}
	return "", false
}

func west236(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)(.)?wt(\d+)?cdn(.)?`)
	headers := response.Headers
	xCache, _ := headers["X-Cache"]
	check := signature1.FindAllString(xCache, -1)
	if check != nil {
		return "West263 CDN", true
	}
	return "", false
}

func wordfence(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)generated.by.wordfence`),
		regexp.MustCompile(`(?i)your.access.to.this.site.has.been.limited`),
		regexp.MustCompile(`(?i)<.+>wordfence<.+.>`),
	}
	body := response.Body
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		if check1 != nil {
			return "Wordfence (Feedjit)", true
		}
	}
	return "", false
}

func wts(response *base.Response) (string, bool) {
	signature1 := regexp.MustCompile(`(?i)wts.waf`)
	check := signature1.FindAllString(response.Body, -1)
	if check != nil {
		return "WTS WAF ", true
	}
	return "", false
}

func yundun(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)YUNDUN`),
		regexp.MustCompile(`(?i)^yd.cookie=(\w+)?`),
		regexp.MustCompile(`(?i)http(s)?.//(www)?.(\w+)?(.)?yundun(.com)?`),
	}
	headers := response.Headers
	xCache, _ := headers["X-Cache"]
	server, _ := headers["Server"]
	cookie, _ := headers["Set-Cookie"]
	for _, s := range signatures {
		check1 := s.FindAllString(xCache, -1)
		check2 := s.FindAllString(server, -1)
		check3 := s.FindAllString(cookie, -1)
		if (check1 != nil) || (check2 != nil) || (check3 != nil) {
			return "Yundun WAF", true
		}
	}
	return "", false
}

func yunsuo(response *base.Response) (string, bool) {
	signatures := []*regexp.Regexp{regexp.MustCompile(`(?i)<img.class=.yunsuologo.`),
		regexp.MustCompile(`(?i)yunsuo.session`),
	}
	headers := response.Headers
	body := response.Body
	cookie, _ := headers["Set-Cookie"]
	for _, s := range signatures {
		check1 := s.FindAllString(body, -1)
		check2 := s.FindAllString(cookie, -1)
		if (check1 != nil) || (check2 != nil) {
			return "Yunsuo WAF", true
		}
	}
	return "", false
}

var wafWG sync.WaitGroup

// 发送请求页面，然后逐一匹配
func checkWAF(url string, ch chan string) {
	defer wafWG.Done()
	headers := map[string]string{
		"Accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	response, err := base.Get(url, headers, true)
	if err != nil {
		return
	}
	// 所有的检测函数数组，按顺序进行检测
	detectFuncs := []detectFunc{qihu360, airlock, akamaiGHost, anquanbao, anyu, armor, asm, aws, baidu, barracuda,
		betterwpsecurity, f5Bigip, f5BigIpAPM, binarySEC, blockDoS, chinaCache, ciscoAceXml, cloudFlare, cloudFront,
		codeIgniter, comodo, ibmDataPower, denyAll, dod, dOSarrest, dotDefender, dynamicWeb, edgecast, expressionEngine,
		f5BigIpLTM, f5FirePass, f5Trafficshield, fortiWeb, gladius, hyperGuard, incapsula, mission, opModsecurity,
		trustwave, owaspModsecurity, knownsec, naxsi, netContinuum, netScaler, nevisProxy, nsFocus, paloAlto, perimeterX,
		pkSecurityModule, powerCDN, powerful, profense, radware, sabre, safeDog, siteGuard, sonicWALL, squid, stingray,
		sucuri, teros, urlScan, uspses, varnish, wallarm, webknight, webSEAL, west236, wordfence, wts, yundun, yunsuo}
	for _, f := range detectFuncs {
		waf, ok := f(response)
		if ok {
			ch <- waf
		}
	}
}

func WAFDetect(u string) string {
	// 正式运行时替换路径，下面的路径仅供测试使用
	dir := base.GetCurrentDirectory()
	path := dir + "/dicts/payloads.txt"
	//fmt.Println(path)

	//path := "D:/Code/GoCode/src/awheel/dicts/payloads.txt"
	lines := base.ReadLines(path)
	num := len(lines)

	wafWG.Add(num)
	ch := make(chan string, num)

	for _, line := range lines {
		attack := u + "?id=1" + url.QueryEscape(line)
		go checkWAF(attack, ch)
	}

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Second * 30) // 等待30秒钟
		timeout <- true
	}()
	select {
	case waf := <-ch:
		return waf
	case <-timeout:
		return ""
	}

	wafWG.Wait()
	close(ch)
	return ""
}
