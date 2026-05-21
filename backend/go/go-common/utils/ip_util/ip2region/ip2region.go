package ip2region

import (
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	geoip "go-common/utils/ip_util"
	"go-common/utils/ip_util/ip2region/assets"
	"strings"
)

type Ip2Region struct {
	v4Search *xdb.Searcher
	v6Search *xdb.Searcher
}

func New() (*Ip2Region, error) {
	client := &Ip2Region{}
	var err error
	client.v4Search, err = xdb.NewWithBuffer(xdb.IPv4, assets.Ip2RegionV4)
	if err != nil {
		return nil, err
	}
	client.v6Search, err = xdb.NewWithBuffer(xdb.IPv6, assets.Ip2RegionV6)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Ip2Region) Query(rawIP string) (res geoip.Result, err error) {
	res.IP = rawIP
	ipBytes, err := xdb.ParseIP(rawIP)
	if err != nil {
		return res, err
	}
	var regionData string
	if l := len(ipBytes); l == 4 {
		regionData, err = c.v4Search.Search(ipBytes)
		if err != nil {
			return res, err
		}
	} else if l == 16 {
		regionData, err = c.v6Search.Search(ipBytes)
		if err != nil {
			return res, err
		}
	} else {
		return res, fmt.Errorf("invalid byte ip address with len=%d", l)
	}

	parts := strings.Split(regionData, "|")
	if len(parts) != 4 {
		return res, fmt.Errorf("invalid region data: %s", regionData)
	}
	res.Country = parts[0]
	res.Province = parts[1]
	res.City = parts[2]
	res.ISP = parts[3]
	return res, nil
}

func (c *Ip2Region) Close() error {
	if c.v4Search != nil {
		c.v4Search.Close()
	}
	if c.v6Search != nil {
		c.v6Search.Close()
	}
	return nil
}
