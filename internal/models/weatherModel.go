package models

type WeatherResponse struct {
	CurrentConditions struct {
		Temp       float64 `json:"temp"`
		Conditions string  `json:"conditions"`
	} `json:"currentConditions"`
}

type UserQuery struct {
	Location string `json:"location"`
}
