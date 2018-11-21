package models

import "time"

type IpList struct {
	Id int64
	Region string `xorm:"varchar(255) notnull default '' 'region'"`
	ServiceName string `xorm:"varchar(255) notnull default '' 'service name'"`
	Ip string `xorm:"varchar(20) notnull default '' 'ip'"`
	Ttl string `xorm:""`
	CreatedAt time.Time `xorm:"TIMESTAMP notnull 'created at'"`
	UpdatedAt time.Time `xorm:"TIMESTAMP notnull 'updated at'"`
}
