package pkg

import (
	"fmt"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w,"Hello, my name is Aisha!Many people make APIs based on their favorite movies/TV series. So, i made an API about my favorite Harry Potter universe🧙")
}