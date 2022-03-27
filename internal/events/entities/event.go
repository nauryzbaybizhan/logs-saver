package eventEntities

import (
	"net"
	"time"
)

type RawEvent struct {
	UserId      string  `json:"id" form:"id"`
	Ip          string  `json:"i" form:"i"`
	ApiKey      string  `json:"k" form:"k"`
	Url         string  `json:"u" form:"u"`
	UserAgent   string  `json:"a" form:"a"`
	RequestTime float64 `json:"t" form:"t"`
}

type Event struct {
	Url         string    `json:"url"`
	UserId      string    `json:"user_id"`
	Ip          net.IP    `json:"ip"`
	ApiKey      string    `json:"api_key"`
	UserAgent   string    `json:"user_agent"`
	IpInfo      *IpInfo   `json:"ip_info"`
	RequestTime time.Time `json:"request_time"`
}
type IpInfo struct {
	Id          int32  `json:"id"`
	Bot         bool   `json:"bot"`
	Datacenter  bool   `json:"datacenter"`
	Tor         bool   `json:"tor"`
	Proxy       bool   `json:"proxy"`
	Vpn         bool   `json:"vpn"`
	Country     string `json:"country"`
	DomainCount string `json:"domaincount"`
	DomainList  string `json:"domain_list"`
}
