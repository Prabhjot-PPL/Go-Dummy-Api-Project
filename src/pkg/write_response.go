package pkg

import (
	"encoding/json"
	"go-project/src/internal/core/coreinterfaces"
	"log"
	"net/http"
)

func WriteResponse(w http.ResponseWriter, resp coreinterfaces.UserAPIHandler) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "Failed to encode Categories", http.StatusInternalServerError)
	} else {
		log.Println("All Categories fetched successfully!!!")
	}
}
