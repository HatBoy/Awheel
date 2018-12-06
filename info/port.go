package info

import (
	"strings"
	"awheel/base"
	"regexp"
	"bytes"
	"encoding/json"
	"awheel/data"
)

/*
调用ZoomEye等搜索引擎获取IP的端口服务信息
*/

func zoomeyeLogin(username, password string) (string, error) {
	url := "https://api.zoomeye.org/user/login"
	data := "{\"username\":\"USER\",\"password\":\"PWD\"}"
	data = strings.Replace(data, "USER", username, -1)
	data = strings.Replace(data, "PWD", password, -1)
	response, err := base.Post(url, nil, data, true)
	if err != nil {
		return "", nil
	}
	body := response.Body
	token := body[18 : len(body)-2]
	return token, nil
}

type ZoomPort struct {
	Extrainfo string `json:extrainfo`
	Service   string `json:service`
	Hostname  string `json:hostname`
	Version   string `json:version`
	Device    string `json:device`
	Os        string `json:os`
	Port      int    `json:port`
	App       string `json:app`
	Banner    string `banner`
}

// 调用ZoomEye接口查询IP端口服务情况
func zoomeyeSearch(ip string) ([]ZoomPort, error) {
	token, err := zoomeyeLogin("hatboy-dj@qq.com", "7r^3&bfswX^K6g85")
	if err != nil {
		return nil, err
	}
	url := "https://api.zoomeye.org/host/search?query=ip:" + ip
	headers := map[string]string{"Authorization": "JWT " + token}

	response, err := base.Get(url, headers, true)
	if err != nil {
		return nil, err
	}

	body := response.Body

	portinfoReg := regexp.MustCompile(`(?U)"portinfo": {(.*)},`)
	portinfos := portinfoReg.FindAllString(body, -1)

	var portString bytes.Buffer
	portString.WriteString("[")

	for _, pi := range portinfos {
		pistring := pi[11 : len(pi)-1]
		portString.WriteString(pistring)
		portString.WriteString(",")
	}
	portString.WriteString("]")

	var portInfos []ZoomPort
	s := portString.String()
	err = json.Unmarshal([]byte(s[0:len(s)-2]+"]"), &portInfos)
	if err != nil {
		return nil, err
	}
	return portInfos, nil
}

// IP端口查询
func IpPort(ip string) ([]data.PortInfo, error) {
	zoomPorts, err := zoomeyeSearch(ip)
	if err != nil {
		return nil, err
	}

	var ipPorts []data.PortInfo

	for _, zp := range zoomPorts {
		pi := new(data.PortInfo)
		pi.Ip = ip
		pi.Extrainfo = zp.Extrainfo
		pi.Service = zp.Service
		pi.Hostname = zp.Hostname
		pi.Version = zp.Version
		pi.Device = zp.Device
		pi.Os = zp.Os
		pi.Port = zp.Port
		pi.App = zp.App
		ipPorts = append(ipPorts, *pi)
	}
	return ipPorts, nil
}

// 调用https://fofa.so接口搜索，调用失败，先不写了
func FofaSearch(ip string) {

}
