package models

import "github.com/astaxie/beego/orm"

type WeeklyReport struct {
	Year     int     `json:"year"`
	Week     int     `json:"week"`
	Distance float64 `json:"avg_distance"`
	Duration float64 `json:"avg_duration"`
}

var rQ = `
SELECT
  EXTRACT(year FROM r.date) as year,
  EXTRACT(week FROM r.date) as week,
  AVG(r.distance) as distance,
  AVG(r.duration) as duration
FROM record r INNER JOIN "user" u on u.id = r.user_id
WHERE
  u.id = ?
GROUP BY
  EXTRACT(year FROM r.date),
  EXTRACT(week FROM r.date)
ORDER BY (
  EXTRACT(year FROM r.date),
  EXTRACT(week FROM r.date)
) DESC;
`

func RecordsGetWeeklyReport(o orm.Ormer, uid uint64) (report []*WeeklyReport, opErr error) {
	var years = []int{}
	var weeks = []int{}
	var distance = []float64{}
	var duration = []float64{}
	_, dbErr := o.Raw(rQ, uid).QueryRows(&years, &weeks, &distance, &duration)
	if dbErr != nil {
		return nil, dbErr
	}
	report = []*WeeklyReport{}
	for idx := range years {
		report = append(report, &WeeklyReport{
			Year:     years[idx],
			Week:     weeks[idx],
			Distance: distance[idx],
			Duration: duration[idx],
		})
	}
	return
}
