package models

import "time"

type IpList struct {
	Id int64
	Region string `xorm:"varchar(255) notnull default '' 'region'"`
	ServiceName string `xorm:"varchar(255) notnull default '' 'service name'"`
	Ip string `xorm:"varchar(20) notnull default '' 'ip'"`
	Ttl string `xorm:"varchar(255) not null default '0' 'ttl'"`
	CreatedAt time.Time `xorm:"TIMESTAMP notnull created default CURRENT_TIMESTAMP 'created at'"`
	UpdatedAt time.Time `xorm:"TIMESTAMP notnull updated default CURRENT_TIMESTAMP 'updated at'"`
}
