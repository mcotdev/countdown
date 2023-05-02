package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

const indexHTML = `
<!DOCTYPE html>
<html>
<head>
	<title>Redirecting in {{ .Countdown }} seconds...</title>
	<meta charset="UTF-8">
	<meta http-equiv="refresh" content="3;url={{ .RedirectURL }}">
	<style>
		body {
			font-family: sans-serif;
			text-align: center;
		}
		#countdown {
			font-size: 4em;
			font-weight: bold;
			fill: white;
		}
		#circle {
			fill: #eee;
			stroke: #555;
			stroke-width: 8;
		}
	</style>
</head>
<body>
	<svg width="200" height="200">
		<circle id="circle" cx="100" cy="100" r="90"></circle>
		<text id="countdown" x="50%" y="50%" dominant-baseline="middle" text-anchor="middle">{{ .Countdown }}</text>
	</svg>
	<p>Redirecting in {{ .Countdown }} seconds...</p>
</body>
</html>
`

func main() {
	// Define the redirect URL and the rate limiter.
	redirectURL := "https://www.google.com"
	rateLimiter := rate.NewLimiter(rate.Every(time.Second), 1)

	// Define the HTTP handler function that displays the countdown page.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the rate limiter allows the request.
		if !rateLimiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		// Set the response header to indicate that the response contains HTML.
		w.Header().Set("Content-Type", "text/html")

		// Parse the HTML template and execute it with the countdown and redirect URL values.
		tmpl, err := template.New("index").Parse(indexHTML)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		data := struct {
			Countdown   int
			RedirectURL string
		}{
			Countdown:   3,
			RedirectURL: redirectURL,
		}
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	})

	// Start the HTTP server and listen for incoming requests.
	fmt.Println("Starting server on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}





