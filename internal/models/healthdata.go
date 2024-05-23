package models

import "time"

type PreDialysisData struct {
	PreBloodPressure string  `json:"pre_blood_pressure"`
	PrePulseRate     int     `json:"pre_pulse_rate"`
	PreTemperature   float64 `json:"pre_temperature"`
	PreWeight        float64 `json:"pre_weight"`
	DryWeight        float64 `json:"dry_weight"`
}

type PostDialysisData struct {
	PostBloodPressure string  `json:"post_blood_pressure"`
	PostPulseRate     int     `json:"post_pulse_rate"`
	PostTemperature   float64 `json:"post_temperature"`
	PostWeight        float64 `json:"post_weight"`
}

type DialysisSession struct {
	Session_id       int `json:"session_id"`
	UserID           int `json:"user_id"`
	PreDialysisData  PreDialysisData
	PostDialysisData PostDialysisData
	WeightGain       float64   `json:"weight_gain"`
	WeightLoss       float64   `json:"weight_loss"`
	SessionDate      time.Time `json:"session_date"`
}
