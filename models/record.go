package models

import (
	"time"

	"encoding/json"

	"github.com/astaxie/beego/orm"
)

func init() {
	orm.RegisterModel(new(Record))
}

type Record struct {
	Id          int64           `orm:"auto"`
	Date        time.Time       `orm:"type(datetime)"`
	Distance    float64         `orm:"digits(6);decimals(3)"`
	Duration    float64         `orm:"digits(12);decimals(3)"`
	Latitude    float64         `orm:"type(real)"`
	Longitude   float64         `orm:"type(real)"`
	User        *User           `orm:"rel(fk);index"`
	WeatherData string          `orm:"type(text)"`
	Weather     json.RawMessage `orm:"-"`
}

func (r *Record) New(o orm.Ormer, linkTo uint64) (opErr error) {
	r.User = &User{Id: int64(linkTo)}
	if r.Id, opErr = o.Insert(r); opErr != nil {
		return
	}
	r.Weather = json.RawMessage(r.WeatherData)
	return
}

func (r *Record) LoadById(o orm.Ormer, id uint64) (opErr error) {
	if opErr = o.QueryTable(r).Filter("Id", id).One(r); opErr != nil {
		return
	}
	r.Weather = json.RawMessage(r.WeatherData)
	return
}

func (r *Record) Save(o orm.Ormer) (opErr error) {
	r.Weather = []byte(r.Weather)
	if _, opErr = o.Update(r); opErr != nil {
		return
	}
	return
}

func (r *Record) Delete(o orm.Ormer) (opErr error) {
	if _, opErr = o.QueryTable(r).Filter("Id", r.Id).Delete(); opErr != nil {
		return
	}
	return
}

func RecordsGetAll(o orm.Ormer, cond *orm.Condition) (records []*Record, opErr error) {
	records = []*Record{}
	if _, opErr = o.QueryTable(Record{}).SetCond(cond).All(&records); opErr != nil {
		return
	}
	for _, r := range records {
		r.Weather = json.RawMessage(r.WeatherData)
	}
	return
}

type RecordView struct {
	Id        int64           `json:"id"`
	Date      time.Time       `json:"date"`
	Distance  float64         `json:"distance"`
	Duration  float64         `json:"duration"`
	Latitude  float64         `json:"latitude"`
	Longitude float64         `json:"longitude"`
	Weather   json.RawMessage `json:"weather"`
}

type RecordData struct {
	Date      time.Time `json:"date"`
	Distance  float64   `json:"distance"`
	Duration  float64   `json:"duration"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}
