package db

import (
	"database/sql"
	"fmt"
	"nephronote/internal/models"
)

func SavePreDialysisData(session models.DialysisSession) error {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	query := `INSERT INTO dialysis_sessions (user_id, pre_blood_pressure, pre_pulse_rate, pre_temperature, pre_weight, dry_weight, weight_gain, session_date) 
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = db.Exec(query, session.UserID, session.PreBloodPressure, session.PrePulseRate, session.PreTemperature, session.PreWeight, session.DryWeight, session.WeightGain, session.SessionDate)
	if err != nil {
		return fmt.Errorf("error saving pre-dialysis data: %v", err)
	}

	return nil
}

func UpdatePostDialysisData(session models.DialysisSession) error {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	query := `UPDATE dialysis_sessions 
              SET post_blood_pressure=$1, post_pulse_rate=$2, post_temperature=$3, post_weight=$4, weight_loss=$5 
              WHERE user_id=$6 AND session_date=$7`
	_, err = db.Exec(query, session.PostBloodPressure, session.PostPulseRate, session.PostTemperature, session.PostWeight, session.WeightLoss, session.UserID, session.SessionDate)
	if err != nil {
		return fmt.Errorf("error updating post-dialysis data: %v", err)
	}

	return nil
}
