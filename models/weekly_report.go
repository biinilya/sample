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
  EXTRACT(year FROM date) as year,
  EXTRACT(week FROM date) as week,
  AVG(distance) as distance,
  AVG(duration) as duration
FROM record
GROUP BY
  EXTRACT(year FROM date),
  EXTRACT(week FROM date)
ORDER BY (
  EXTRACT(year FROM date),
  EXTRACT(week FROM date)
) DESC;
`

func RecordsGetWeeklyReport(o orm.Ormer) (report []*WeeklyReport, opErr error) {
	var years = []int{}
	var weeks = []int{}
	var distance = []float64{}
	var duration = []float64{}
	_, dbErr := o.Raw(rQ).QueryRows(&years, &weeks, &distance, &duration)
	if dbErr != nil {
		return nil, dbErr
	}
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
