package services

import (
	"github.com/luoxiaojun1992/http-dns/models"
	"github.com/go-xorm/xorm"
)

type ipService struct {

}

var IpService ipService

func init()  {
	//todo DI

	IpService = ipService{}
}

func (s ipService) GetList (region, serviceName string, orm *xorm.Engine) ([]models.IpList, error) {
	ips := make([]models.IpList, 0, 10)

	err := orm.Where("region = ? AND service_name = ?", region, serviceName).
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
