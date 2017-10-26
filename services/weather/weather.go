package weather

import (
	"fmt"
	"time"

	"net/http"

	"errors"
	"io/ioutil"

	"encoding/json"

	"github.com/astaxie/beego"
)

var apiKey string

func Init() {
	apiKey = beego.AppConfig.String("WeatherKey")
}

const forecastUrlTemplate = "https://api.darksky.net/forecast/%s/%f,%f,%s?exclude=currently,flags,hourly"

func GetWeather(lat, lng float64, when time.Time) (json.RawMessage, error) {
	var url = fmt.Sprintf(forecastUrlTemplate, apiKey, lat, lng, when.Format(time.RFC3339))
	var resp, respErr = http.Get(url)
	if respErr != nil {
		return nil, respErr
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	var data, readErr = ioutil.ReadAll(resp.Body)
	if readErr != nil {
		return nil, readErr
	}
	return data, nil
}
