package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		uri := r.URL.Path
		method := r.Method

		if method == "POST" {
			input := r.FormValue("input")
			if input == "" {
				http.Error(w, "Input is required", http.StatusInternalServerError)
				return
			}

			var response string
			switch uri {
			case "/sha256":
				hash := sha256.Sum256([]byte(input))
				response = fmt.Sprintf("%x", hash)
			case "/base64":
				response = base64.StdEncoding.EncodeToString([]byte(input))
			case "/urlencode":
				response = url.QueryEscape(input)
			default:
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			_, _ = w.Write([]byte(response))
		} else if method == "GET" {
			w.WriteHeader(201)
			w.Write([]byte("go"))
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	http.ListenAndServe(":8008", nil)
}
