package qqwry

import (
	"errors"
	"github.com/xiaoqidun/qqwry"
	geoip "go-common/utils/ip_util"
	"go-common/utils/ip_util/qqwry/assets"
)

type QQWry struct {
	client *qqwry.Client
}

func New() (*QQWry, error) {
	client, err := qqwry.NewClientFromData(assets.QQWryIPDB)
	if err != nil {
		return nil, err
	}
	return &QQWry{client}, nil
}

func (q *QQWry) Close() error {
	return nil
}

func (q *QQWry) Query(rawIP string) (res geoip.Result, err error) {
	location, err := q.client.QueryIP(rawIP)
	if err != nil {
		return res, err
	}
	if location == nil {
		return res, errors.New("IP Not Found")
	}
	res.Country = location.Country
	res.City = location.City
	res.ISP = location.ISP
	res.Province = location.Province
	return res, nil
}
