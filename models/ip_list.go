package models

import "time"

type IpList struct {
	Id int64 `json:"-"`
	Region string `xorm:"varchar(255) notnull default '' 'region'" json:"-"`
	ServiceName string `xorm:"varchar(255) notnull default '' 'service_name'" json:"-"`
	Ip string `xorm:"varchar(20) notnull default '' 'ip'" json:"ip"`
	Ttl string `xorm:"varchar(255) not null default '0' 'ttl'" json:"ttl"`
	CreatedAt time.Time `xorm:"TIMESTAMP notnull created default CURRENT_TIMESTAMP 'created_at'" json:"-"`
	UpdatedAt time.Time `xorm:"TIMESTAMP notnull updated default CURRENT_TIMESTAMP 'updated_at'" json:"-"`
}
