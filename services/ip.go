package services

import (
	"github.com/luoxiaojun1992/http-dns/models"
	"github.com/luoxiaojun1992/http-dns/utils"
)

type ipService struct {
}

var IpService ipService

func init() {
	//todo DI

	IpService = ipService{}
}

func (s ipService) GetList(region, serviceName string) ([]models.IpList, error) {
	ips := make([]models.IpList, 0, 10)

	err := utils.Orm.Where("region = ? AND service_name = ?", region, serviceName).
		Limit(10).
		OrderBy("updated_at DESC").
		Select("ip, ttl").
		Find(&ips)
	if err == nil {
		return ips, nil
	} else {
		return []models.IpList{}, err
	}
}

func (s ipService) Add(region, serviceName, ip, ttl string) (int64, error) {
	return utils.Orm.Insert(models.IpList{
		Region:      region,
		ServiceName: serviceName,
		Ip:          ip,
		Ttl:         ttl,
	})
}

func (s ipService) Delete(region, serviceName string) (int64, error) {
	return utils.Orm.OrderBy("updated_at DESC").Limit(10).Delete(models.IpList{
		Region:      region,
		ServiceName: serviceName,
	})
}
