package region

import (
	"github.com/ihuanglei/authenticator/pkg/logger"
	"github.com/lionsoul2014/ip2region/binding/golang/ip2region"
	"github.com/simplexwork/common"
)

var (
	ip2Region *ip2region.Ip2Region
	err       error
)

// Region 地区信息
type Region struct {
	ip2region.IpInfo
}

func init() {
	if ip2Region, err = ip2region.New("lib/ip2region.db"); err != nil {
		logger.Error(err)
	}
}

// IP2Region ip转地区信息
func IP2Region(ip string) (*Region, error) {
	if ip2Region == nil {
		return nil, err
	}
	ipInfo, err := ip2Region.MemorySearch(ip)
	if err != nil {
		logger.Warn(err)
	}

	if ipInfo.Country == "0" || common.IsEmpty(ipInfo.Country) {
		ipInfo.Country = "-"
	}

	if ipInfo.Province == "0" || common.IsEmpty(ipInfo.Province) {
		ipInfo.Province = "-"
	}

	if ipInfo.City == "0" || common.IsEmpty(ipInfo.City) {
		ipInfo.City = "-"
	}

	if ipInfo.Region == "0" || common.IsEmpty(ipInfo.Region) {
		ipInfo.Region = "-"
	}

	region := Region{IpInfo: ipInfo}

	return &region, nil

}
