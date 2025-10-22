package logger

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func Middleware(next http.Handler, pathLog string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dir := filepath.Dir(pathLog)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create log directory: %v\n", err)
		}
		f, err := os.OpenFile(pathLog, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0755)
		if err != nil {
			log.Printf("Error opening log file: %v\n", err)
			return
		}
		defer func() {
			err = f.Close()
			if err != nil {
				log.Printf("Error closing log file: %v\n", err)
			}
		}()

		_, err = f.Write([]byte(r.Method + " " + r.URL.Path + " " + time.Now().String() + "\n"))
		if err != nil {
			log.Printf("Error writing to log file: %v\n", err)
		}

		next.ServeHTTP(w, r)
	})
}
