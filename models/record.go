package models

import (
	"time"

	"github.com/astaxie/beego/orm"
)

type Record struct {
	Id        int64     `orm:"auto"`
	Date      time.Time `orm:"type(datetime)"`
	Distance  float64   `orm:"digits(6);decimals(3)"`
	Duration  float64   `orm:"digits(12);decimals(3)"`
	Latitude  float64   `orm:"type(real)"`
	Longitude float64   `orm:"type(real)"`
	User      *User     `orm:"rel(fk);index"`
}

func init() {
	orm.RegisterModel(new(Record))
}
