package services

import (
	"context"
	"go-common/utils/ip_util"
	"go-common/utils/ip_util/ip2region"
)

type Geo struct {
	client ip_util.GeoIP
}

func NewGeo() *Geo {
	return &Geo{}
}

func (g *Geo) Start(ctx context.Context) error {
	region, err := ip2region.New()
	if err != nil {
		return err
	}
	ip_util.Client = region
	g.client = region
	return nil
}

func (g *Geo) String() string {
	return "geo"
}

func (g *Geo) State(ctx context.Context) (string, error) {
	return "HEALTHY", nil
}

func (g *Geo) Terminate(ctx context.Context) error {
	return g.client.Close()
}
