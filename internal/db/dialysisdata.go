package db

import (
	"database/sql"
	"fmt"
	"nephronote/internal/models"
)

func SavePreDialysisData(session models.DialysisSession) (int, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return 0, fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	var sessionID int
	query := `INSERT INTO dialysis_sessions (user_id, pre_blood_pressure, pre_pulse_rate, pre_temperature, pre_weight, dry_weight, weight_gain, session_date)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING session_id`
	err = db.QueryRow(query, session.UserID, session.PreDialysisData.PreBloodPressure, session.PreDialysisData.PrePulseRate,
		session.PreDialysisData.PreTemperature, session.PreDialysisData.PreWeight, session.PreDialysisData.DryWeight, session.WeightGain).Scan(&sessionID)
	if err != nil {
		return 0, fmt.Errorf("error saving pre-dialysis data: %v", err)
	}

	fmt.Printf("Pre-dialysis data saved successfully with session ID: %d\n", sessionID)

	return sessionID, nil
}

func GetDialysisSession(sessionID int, userID int) (models.DialysisSession, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return models.DialysisSession{}, fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	var session models.DialysisSession
	query := `SELECT session_id, user_id, pre_blood_pressure, pre_pulse_rate, pre_temperature, pre_weight, dry_weight, weight_gain, session_date 
              FROM dialysis_sessions WHERE session_id = $1 AND user_id = $2`
	err = db.QueryRow(query, sessionID, userID).Scan(&session.Session_id, &session.UserID, &session.PreDialysisData.PreBloodPressure,
		&session.PreDialysisData.PrePulseRate, &session.PreDialysisData.PreTemperature, &session.PreDialysisData.PreWeight,
		&session.PreDialysisData.DryWeight, &session.WeightGain, &session.SessionDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return session, fmt.Errorf("session not found")
		}
		return session, fmt.Errorf("error fetching session data: %v", err)
	}

	fmt.Printf("Retrieved session from database: %+v\n", session)

	return session, nil
}

func UpdatePostDialysisData(session models.DialysisSession) error {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return fmt.Errorf("error opening database connection: %v", err)
	}
	defer db.Close()

	query := `UPDATE dialysis_sessions SET post_blood_pressure = $1, post_pulse_rate = $2, post_temperature = $3, post_weight = $4, weight_loss = $5 
              WHERE session_id = $6 AND user_id = $7`
	_, err = db.Exec(query, session.PostDialysisData.PostBloodPressure, session.PostDialysisData.PostPulseRate,
		session.PostDialysisData.PostTemperature, session.PostDialysisData.PostWeight, session.WeightLoss, session.Session_id, session.UserID)
	if err != nil {
		return fmt.Errorf("error updating post-dialysis data: %v", err)
	}

	fmt.Println("Post-dialysis data updated successfully")

	return nil
}
