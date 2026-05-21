package geolite

import (
	"errors"
	"github.com/oschwald/geoip2-golang/v2"
	geoip "go-common/utils/ip_util"
	"go-common/utils/ip_util/geolite/assets"
	"net/netip"
)

type GeoLite struct {
	client *geoip2.Reader
}

func New() (*GeoLite, error) {
	db, err := geoip2.OpenBytes(assets.GeoLite2CityData)
	if err != nil {
		return nil, err
	}
	return &GeoLite{client: db}, nil
}

func (g *GeoLite) Close() error {
	if g.client == nil {
		return nil
	}
	return g.client.Close()
}

func (g *GeoLite) Query(rawIP string) (res geoip.Result, err error) {
	ip, err := netip.ParseAddr(rawIP)
	if err != nil {
		return res, errors.New("invalid ip address")
	}

	record, err := g.client.City(ip)
	if err != nil {
		return res, err
	}
	if !record.HasData() {
		return res, errors.New("no data found for this ip")
	}

	res.Country = record.Country.Names.SimplifiedChinese
	if len(record.Subdivisions) > 0 {
		res.Province = record.Subdivisions[0].Names.SimplifiedChinese
	}
	res.City = record.City.Names.SimplifiedChinese
	return
}
