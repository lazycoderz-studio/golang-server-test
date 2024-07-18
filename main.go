package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type EmailRequest struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
	HTML    string `json:"html"`
}

func sendEmailHandler(w http.ResponseWriter, r *http.Request) {
	var emailReq EmailRequest
	err := json.NewDecoder(r.Body).Decode(&emailReq)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	from := mail.NewEmail("Example Sender", emailReq.From)
	to := mail.NewEmail("Example Receiver", emailReq.To)
	message := mail.NewSingleEmail(from, emailReq.Subject, to, emailReq.Text, emailReq.HTML)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.WriteHeader(response.StatusCode)
	w.Write([]byte(response.Body))
}

func receiveMail(w http.ResponseWriter, r *http.Request) {
	req := map[string]interface{}{}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Println(req)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]string{
			"host":   r.Host,
			"status": "success",
		}

		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON response
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("error encoding response: %v", err)
			http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
			return
		}
	})

	r.HandleFunc("/send-email", sendEmailHandler).Methods("POST")
	r.HandleFunc("/receive-mail", receiveMail).Methods("POST")
	// Start the HTTP server
	err := http.ListenAndServe(":8090", r)
	if err != nil {
		log.Println("Error while running server:", err)
	}
}
