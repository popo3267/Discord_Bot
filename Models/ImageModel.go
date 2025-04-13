package models

import "time"

type WeatherData struct {
	Success string `json:"success"`
	Records struct {
		Station []struct {
			StationName string `json:"stationName"`
			GeoInfo     struct {
				CountyName string `json:"countyName"`
			}
			ObsTime struct {
				DateTime string `json:"datetime"`
			} `json:"obstime"`
			WeatherElement struct {
				AirTemperature float64 `json:"airtemperature"`
			} `json:"weatherElement"`
		} `json:"station"`
	} `json:"records"`
}

type ImageManage struct {
	ImageID    int64     `json:"imageid"`
	Keyword    string    `json:"keyword"`
	ImageURL   string    `json:"imageurl"`
	UserName   string    `json:"username"`
	CreateTime time.Time `json:"createtime"`
}
