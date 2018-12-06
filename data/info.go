package data

/*
存放info包中的struct结构
*/

// baseinfo数据结构
type UrlBaseInfo struct {
	Domain       string   //域名
	Ip           string   // IP
	IpList       []string // IP列表
	IsCDN        bool     //是否是CDN
	Country      string   //国家
	Region       string   //省
	City         string   //市
	Isp          string   //运营商
	Registrar    string   //注册商
	Registrant   string   //注册人
	Emali        string   //联系邮箱
	Phone        string   //联系电话
	CreateDate   string   //创建时间
	ExpireDate   string   //过期时间
	DomainServer string   //域名服务器
	DNS          []string //DNS列表
	Status       []string //状态
	Cms          string   //CMS类型
	Waf          string   //WAF类型
	Url          string
	Title        string
	PoweredBy    string
	Server       string
	Via          string
	Robots       string
	PortInfos    []PortInfo
}

// 域名IP信息
type IpInfo struct {
	Domain  string   //域名
	Ip      string   // IP
	IpList  []string // IP列表
	IsCDN   bool     //是否是CDN
	Country string   //国家
	Region  string   //省
	City    string   //市
	Isp     string   //运营商
}

//端口信息
type PortInfo struct {
	Ip        string
	Extrainfo string
	Service   string
	Hostname  string
	Version   string
	Device    string
	Os        string
	Port      int
	App       string
}

//Webinfo
type WebInfoData struct {
	Url       string
	Title     string
	PoweredBy string
	Server    string
	Via       string
	Robots    string
}

type WhoisInfo struct {
	Domain       string   //域名
	Registrar    string   //注册商
	Registrant   string   //注册人
	Emali        string   //联系邮箱
	Phone        string   //联系电话
	CreateDate   string   //创建时间
	ExpireDate   string   //过期时间
	DomainServer string   //域名服务器
	DNS          []string //DNS列表
	Status       []string //状态
}

// 任务数据结构
type Target struct {
	TargetId         string //任务ID，每一个任务共用唯一一个ID
	Proto            string //协议
	Domain           string //完整的域名
	Master           bool   //是否是主任务，只有最开始的任务是主任务，其他的都是从任务
	DomainId         string //域名ID，同一个二级域名共用一个ID
	Types            string //任务类型，包括baseinfo、subdomain、subdir
	IsSubdomian      bool   //是否查询子域名
	IsDomainPan      bool   //域名是否是泛解析
	IsBruteSubdomain bool   //默认只进行接口查询，通过设置可选是否进行爆破子域名
	FirstDomain      string //需要查询子域名时一级域名
	IsSubdir         bool   //是否查询子目录
	IsAlive          bool   //域名是否可以访问，不可访问则不进行子目录爆破
	Subdir           string //需要目录爆破的时候爆破的目录
	DirType          string //子目录类别，这里指的是主任务的子目录,php,asp,aspx,jsp,dir,mdb,auto自动判断
	SubdirType       string //子目录类别，这里指的是其他子域名的子目录,php,asp,aspx,jsp,dir,mdb,auto自动判断，auto只能单独存在
	IsSubdirRec      bool   //子域名是否递归查询子目录
}

// baseinfo结果数据结构
type BaseInfoResult struct {
	Task     Target
	Baseinfo UrlBaseInfo
}

type DirResult struct {
	TargetId    string //任务ID，每一个任务共用唯一一个ID
	FirstDomain string //需要查询子域名时一级域名
	Proto       string //协议
	Domain      string //完整的域名
	DomainId    string //域名ID，同一个二级域名共用一个ID
	Subdir      string //需要目录爆破的时候爆破的目录
	StatusCode  int    //状态码
	Master      bool   //是否是主任务，只有最开始的任务是主任务，其他的都是从任务
	DirType     string //子目录类别，这里指的是主任务的子目录,php,asp,aspx,jsp,dir,mdb,auto自动判断
	SubdirType  string //子目录类别，这里指的是其他子域名的子目录,php,asp,aspx,jsp,dir,mdb,auto自动判断，auto只能单独存在
	IsSubdirRec bool   //子域名是否递归查询子目录
}