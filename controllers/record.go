package controllers

import (
	"encoding/json"
	"toptal/lib"
	"toptal/models"
	"toptal/presenter"
	"toptal/services/weather"
)

//  RecordController operations for Record
type RecordController struct {
	ApiController
}

var recordView = presenter.New().AddHook(
	models.Record{}, presenter.StructView().TranslateCase().Exclude("User", "WeatherData"),
)

var recordDataView = presenter.New().AddHook(
	models.Record{}, presenter.StructView().TranslateCase().Exclude("Id", "User", "Weather", "WeatherData"),
)

func (c *ApiController) LoadRecordFromBody(record *models.Record) {
	if dataInErr := json.Unmarshal(c.Ctx.Input.RequestBody, recordDataView.AsJson(&record)); dataInErr != nil {
		c.AbortWith(400, dataInErr)
	}

	switch {
	case record.Distance >= 1000 || record.Distance <= 0:
		c.AbortWith(400, "'distance' should be a positive number lesser then 1000")
	case record.Duration >= 1000000 || record.Duration <= 0:
		c.AbortWith(400, "'duration' should be a positive number lesser then 10000000")
	}
	return
}

// @Title Post
// @Description create new Record, returns new record data
// @Param   X-Access-Token header  string              true	 "Access Token"
// @Param   uid            path    uint64              true  "User ID"
// @Param   body           body    models.RecordData   true  "Record"
// @Success 201 {object} models.RecordView
// @Failure 400 Bad request
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @Failure 503 weather service is down
// @router / [post]
func (c *RecordController) Post() {
	var u = c.RequireOwnerOrPerm(models.PERM_ADMIN)
	var record models.Record
	c.LoadRecordFromBody(&record)

	if weatherData, err := weather.GetWeather(record.Latitude, record.Longitude, record.Date); err != nil {
		c.AbortWith(503, err)
	} else {
		record.WeatherData = string(weatherData)
	}

	if dbErr := record.New(lib.GetDB(), uint64(u.Id)); dbErr != nil {
		c.AbortWith(500, dbErr)
	}

	c.Ctx.Output.SetStatus(201)
	c.Data["json"] = recordView.AsJson(&record)
	c.ServeJSON()
}

// @Title Put
// @Description update existing record, returns updated record data
// @Param   X-Access-Token header  string              true	 "Access Token"
// @Param   uid            path    uint64              true  "User ID"
// @Param   record_id      path    uint64              true  "Record ID"
// @Param   body           body    models.RecordData   true  "Record"
// @Success 200 {object} models.RecordView
// @Failure 400 Bad request
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router /:record_id [put]
func (c *RecordController) Put() {
	var recordId, recordIdErr = c.GetUint64(":record_id")
	if recordIdErr != nil {
		c.AbortWith(400, "Invalid :record_id")
	}
	c.RequireOwnerOrPerm(models.PERM_ADMIN)
	var record models.Record
	if dbErr := record.LoadById(lib.GetDB(), recordId); dbErr != nil {
		c.AbortWith(403, "Access to record is forbidden")
	}
	c.LoadRecordFromBody(&record)

	if dbErr := record.Save(lib.GetDB()); dbErr != nil {
		c.AbortWith(500, dbErr)
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = recordView.AsJson(&record)
	c.ServeJSON()
}

// @Title Delete
// @Description update existing record, returns updated record data
// @Param   X-Access-Token header  string  true  "Access Token"
// @Param   uid            path    uint64  true  "User ID"
// @Param   record_id      path    uint64  true  "Record ID"
// @Success 200 {object} models.RecordView
// @Failure 400 Bad request
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router /:record_id [delete]
func (c *RecordController) Delete() {
	var recordId, recordIdErr = c.GetUint64(":record_id")
	if recordIdErr != nil {
		c.AbortWith(400, "Invalid :record_id")
	}
	c.RequireOwnerOrPerm(models.PERM_ADMIN)
	var record models.Record
	if dbErr := record.LoadById(lib.GetDB(), recordId); dbErr != nil {
		c.AbortWith(403, "Access to record is forbidden")
	}

	if dbErr := record.Delete(lib.GetDB()); dbErr != nil {
		c.AbortWith(500, dbErr)
	}

	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = recordView.AsJson(&record)
	c.ServeJSON()
}

// @Title GetAll
// @Description Get all user-related records
// @Param   X-Access-Token header  string  true  "Access Token"
// @Param   uid            path    uint64  true  "User ID"
// @Param   filter         query   string  false "Filter records e.x. (date = '2017-01-01')"
// @Success 200 {object} []models.RecordView
// @Failure 400 Bad request
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router / [get]
func (c *RecordController) GetAll() {
	var u = c.RequireOwnerOrPerm(models.PERM_ADMIN)
	var f = c.LoadFilter("date", "distance", "duration", "latitude", "longitude").And("user__id", u.Id)

	var records, dbErr = models.RecordsGetAll(lib.GetDB(), f)
	if dbErr != nil {
		c.AbortWith(500, dbErr)
	}

	c.Data["json"] = recordView.AsJson(&records)
	c.ServeJSON()
}

// @Title WeeklyReport
// @Description Get average distance and duration per week
// @Param   X-Access-Token header  string  true  "Access Token"
// @Param   uid            path    uint64  true  "User ID"
// @Success 200 {object} []models.WeeklyReport
// @Failure 400 Bad request
// @Failure 401 unauthorized
// @Failure 403 forbidden
// @router /report/weekly [get]
func (c *RecordController) WeeklyReport() {
	c.RequireOwnerOrPerm(models.PERM_ADMIN)

	var report, dbErr = models.RecordsGetWeeklyReport(lib.GetDB())
	if dbErr != nil {
		c.AbortWith(500, dbErr)
	}

	c.Data["json"] = report
	c.ServeJSON()
}
