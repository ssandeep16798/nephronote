package models

type DialysisSession struct {
	UserID            int     `json:"user_id"`
	PreBloodPressure  string  `json:"pre_blood_pressure"`
	PrePulseRate      int     `json:"pre_pulse_rate"`
	PreTemperature    float64 `json:"pre_temperature"`
	PreWeight         float64 `json:"pre_weight"`
	DryWeight         float64 `json:"dry_weight"`
	WeightGain        float64 `json:"weight_gain"`
	PostBloodPressure string  `json:"post_blood_pressure"`
	PostPulseRate     int     `json:"post_pulse_rate"`
	PostTemperature   float64 `json:"post_temperature"`
	PostWeight        float64 `json:"post_weight"`
	WeightLoss        float64 `json:"weight_loss"`
	SessionDate       string  `json:"session_date"`
}
