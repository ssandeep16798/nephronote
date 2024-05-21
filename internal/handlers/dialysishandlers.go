package handlers

import (
	"encoding/json"
	"fmt"
	"nephronote/internal/db"
	"nephronote/internal/models"
	"net/http"
	"time"

	"github.com/labstack/gommon/log"
)

func PreDialysisHandler(w http.ResponseWriter, req *http.Request) {
	var session models.DialysisSession
	decoder := json.NewDecoder(req.Body)
	fmt.Println("inside Pre-Handler")
	if err := decoder.Decode(&session); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	log.Error("error before calculating weights")

	session.WeightGain = session.PreWeight - session.DryWeight
	log.Error("error calculating weights.")
	session.SessionDate = time.Now().Format("2006-01-02")

	err := db.SavePreDialysisData(session)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":      true,
		"weight_gain": session.WeightGain,
		"message":     "Pre-dialysis data saved successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func PostDialysisHandler(w http.ResponseWriter, req *http.Request) {
	var session models.DialysisSession
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&session); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	session.WeightLoss = session.PreWeight - session.PostWeight

	err := db.UpdatePostDialysisData(session)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	effective := session.PostWeight == session.DryWeight
	if effective {
		fmt.Println("You had an effective dialysis session.")
	} else {
		fmt.Println("Drink less water until you get rid off those excess water..")
	}
	response := map[string]interface{}{
		"status":      true,
		"effective":   effective,
		"weight_loss": session.WeightLoss,
		"message":     "Post-dialysis data saved successfully",
	}
	json.NewEncoder(w).Encode(response)
}
