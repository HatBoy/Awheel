package info

import (
	"awheel/base"
	"regexp"
	"net/url"
	"bytes"
	"strings"
	"encoding/binary"
	"net"
	"fmt"
	"log"
)

func virustotalParse(body string) []string {
	var subDomains []string
	idsReg := regexp.MustCompile(`(?U)"id": "(.*)"`)
	ids := idsReg.FindAllString(body, -1)
	for _, sub := range ids {
		sub = sub[7 : len(sub)-1]
		subDomains = append(subDomains, sub)
	}
	return subDomains
}

func virustotal(domain string) []string {
	url1 := "https://www.virustotal.com/ui/domains/" + domain + "/subdomains"
	headers := map[string]string{
		"Accept":          "application/json",
		"Accept-Encoding": "gzip",
		"Accept-Language": "zh-CN,zh;q=0.9",
	}
	response, _ := base.Get(url1, headers, true)
	body := response.Body
	sub1 := virustotalParse(body)
	cursorReg := regexp.MustCompile(`(?U)"cursor": "(.*)"`)
	cursors := cursorReg.FindAllString(body, -1)
	var cursor string
	if len(cursors) == 1 {
		cursor = cursors[0][11 : len(cursors[0])-1]
	}
	url2 := "https://www.virustotal.com/ui/domains/" + domain + "/subdomains?cursor=" + url.QueryEscape(cursor)
	response2, _ := base.Get(url2, headers, true)
	sub2 := virustotalParse(response2.Body)
	var subDomains []string
	for _, s1 := range sub1 {
		subDomains = append(subDomains, s1)
	}

	for _, s2 := range sub2 {
		subDomains = append(subDomains, s2)
	}
	return subDomains
}

type dnsHeader struct {
	Id                                 uint16
	Bits                               uint16
	Qdcount, Ancount, Nscount, Arcount uint16
}

func (header *dnsHeader) SetFlag(QR uint16, OperationCode uint16, AuthoritativeAnswer uint16, Truncation uint16, RecursionDesired uint16, RecursionAvailable uint16, ResponseCode uint16) {
	header.Bits = QR<<15 + OperationCode<<11 + AuthoritativeAnswer<<10 + Truncation<<9 + RecursionDesired<<8 + RecursionAvailable<<7 + ResponseCode
}

type dnsQuery struct {
	QuestionType  uint16
	QuestionClass uint16
}

func ParseDomainName(domain string) []byte {
	var (
		buffer   bytes.Buffer
		segments []string = strings.Split(domain, ".")
	)
	for _, seg := range segments {
		binary.Write(&buffer, binary.BigEndian, byte(len(seg)))
		binary.Write(&buffer, binary.BigEndian, []byte(seg))
	}
	binary.Write(&buffer, binary.BigEndian, byte(0x00))

	return buffer.Bytes()
}
func Send(dnsServer, domain string) ([]byte, int) {
	requestHeader := dnsHeader{
		Id:      0x0010,
		Qdcount: 1,
		Ancount: 0,
		Nscount: 0,
		Arcount: 0,
	}
	requestHeader.SetFlag(0, 0, 0, 0, 1, 0, 0)

	requestQuery := dnsQuery{
		QuestionType:  1,
		QuestionClass: 1,
	}

	var (
		conn   net.Conn
		err    error
		buffer bytes.Buffer
	)

	if conn, err = net.Dial("udp", dnsServer); err != nil {
		fmt.Println(err.Error())
		return make([]byte, 0), 0
	}
	defer conn.Close()

	binary.Write(&buffer, binary.BigEndian, requestHeader)
	binary.Write(&buffer, binary.BigEndian, ParseDomainName(domain))
	binary.Write(&buffer, binary.BigEndian, requestQuery)

	buf := make([]byte, 1024)
	if _, err := conn.Write(buffer.Bytes()); err != nil {
		fmt.Println(err.Error())
		return make([]byte, 0), 0
	}
	length, err := conn.Read(buf)
	return buf, length
}

// 用于子域名爆破
func Nslookup(domain string) (bool, error) {
	resolver := base.New([]string{"223.5.5.5", "114.114.114.114", "119.29.29.29", "182.254.116.116", "180.76.76.76"})
	resolver.RetryTimes = 10

	ips, err := resolver.LookupHost(domain)
	if err != nil {
		log.Fatalln(err)
		return false, err
	}

	if len(ips) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

//判断一级域名是否是泛解析
func DomainIsPan(subdomain string) bool {
	randSub := base.GetRandomString() + "." + subdomain
	ok, _ := Nslookup(randSub)
	return ok
}

func SubDomain(domain string) []string {
	subs := virustotal(domain)
	return subs

}
