package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"nephronote/internal/db"
	"nephronote/internal/models"
	"net/http"
	"time"
)

func PreDialysisHandler(w http.ResponseWriter, req *http.Request) {
	var preData models.PreDialysisData
	decoder := json.NewDecoder(req.Body)
	fmt.Println("inside Pre-Handler")
	if err := decoder.Decode(&preData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Check if UserID exists in the context
	userID, ok := req.Context().Value("userID").(int)
	if !ok {
		fmt.Println("Error: UserID not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Println("PreDialysisHandler: UserID found in context:", userID)

	fmt.Println("UserID:", userID)

	fmt.Println("before session")
	session := models.DialysisSession{
		UserID:          userID,
		PreDialysisData: preData,
		WeightGain:      preData.PreWeight - preData.DryWeight,
		SessionDate:     time.Now(),
	}
	fmt.Println("after session.")

	err := db.SavePreDialysisData(session)
	if err != nil {
		fmt.Println("Error saving pre-dialysis data:", err)
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
	var postData models.PostDialysisData
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&postData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	sessionID := req.Context().Value("sessionID").(int) // Assuming sessionID is passed in context after pre-dialysis data entry
	session, err := db.GetDialysisSession(sessionID)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session.PostDialysisData = postData
	session.WeightLoss = session.PreDialysisData.PreWeight - postData.PostWeight

	err = db.UpdatePostDialysisData(session)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	effective := session.PostDialysisData.PostWeight == session.PreDialysisData.DryWeight
	responseMessage := "You had an effective dialysis session."
	if !effective {
		responseMessage = "Drink less water until you get rid of the excess water."
	}

	response := map[string]interface{}{
		"status":      true,
		"effective":   effective,
		"weight_loss": session.WeightLoss,
		"message":     responseMessage,
	}
	json.NewEncoder(w).Encode(response)
}
