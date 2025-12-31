package models

type WeatherResponse struct {
	CurrentConditions struct {
		Temp       float64 `json:"temp"`
		Conditions string  `json:"conditions"`
	} `json:"currentConditions"`
}
