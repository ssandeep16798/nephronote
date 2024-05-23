package handlers

import (
	"encoding/json"
	"fmt"
	"nephronote/internal/db"
	"nephronote/internal/models"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func PreDialysisHandler(w http.ResponseWriter, req *http.Request) {
	var preData models.PreDialysisData
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&preData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userIDInterface := req.Context().Value("userID")
	if userIDInterface == nil {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		http.Error(w, "User ID not of type int", http.StatusInternalServerError)
		return
	}

	session := models.DialysisSession{
		UserID:          userID,
		PreDialysisData: preData,
		WeightGain:      preData.PreWeight - preData.DryWeight,
	}

	sessionID, err := db.SavePreDialysisData(session)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println("Pre-dialysis data saved successfully with session ID:", sessionID)

	response := map[string]interface{}{
		"status":      true,
		"message":     "Pre-dialysis data saved successfully",
		"weight_gain": session.WeightGain,
		"session_id":  sessionID, // Add session_id to the response
	}
	json.NewEncoder(w).Encode(response)
}
func PostDialysisHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	sessionIDStr, ok := vars["sessionID"]
	if !ok {
		http.Error(w, "Session ID is missing in the URL", http.StatusBadRequest)
		return
	}

	sessionID, err := strconv.Atoi(sessionIDStr)
	if err != nil {
		http.Error(w, "Invalid session ID format", http.StatusBadRequest)
		return
	}

	var postData models.PostDialysisData
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&postData); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	userIDInterface := req.Context().Value("userID")
	if userIDInterface == nil {
		http.Error(w, "User ID not found in context", http.StatusInternalServerError)
		return
	}

	userID, ok := userIDInterface.(int)
	if !ok {
		http.Error(w, "User ID not of type int", http.StatusInternalServerError)
		return
	}

	session, err := db.GetDialysisSession(sessionID, userID)
	if err != nil {
		http.Error(w, "Error fetching dialysis session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if session.UserID != userID {
		http.Error(w, "Invalid session ID", http.StatusUnauthorized)
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
