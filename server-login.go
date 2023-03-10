package main

import (
	"bytes"
	"github.com/bmizerany/pat"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

const (
	username = "user@example.com"
	password = "Zapote1234567!" // hashed password
)

type Message struct {
	Email    string
	Password string
}

func main() {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(home))
	mux.Post("/", http.HandlerFunc(handleLogin))
	mux.Get("/success", http.HandlerFunc(confirmation))

	http.ListenAndServe(":8080", mux)
}

func (v *Message) Validate() bool {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(v.Password), bcrypt.DefaultCost)

	if err != nil {
		return false
	}

	errResult := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))

	return v.Email == username && errResult == nil
}

func createLogger() *log.Logger {
	var buf bytes.Buffer
	return log.New(&buf, "logger: ", log.Ldate)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {

	logger := createLogger()

	if r.Method == "POST" {
		msg := &Message{
			Email:    r.PostFormValue("email"),
			Password: r.PostFormValue("password"),
		}
		if msg.Validate() == false {
			render(w, "pages/login.html", msg)
			return
		}
		logger.Println("Login Successful")
		http.Redirect(w, r, "/success", http.StatusSeeOther)
		return
	}

}

func render(w http.ResponseWriter, templateName string, data interface{}) {
	tmpl, err := template.ParseFiles(templateName)
	logger := createLogger()
	if err != nil {
		logger.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		logger.Print(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}
func home(w http.ResponseWriter, r *http.Request) {
	render(w, "pages/login.html", nil)
}
func confirmation(w http.ResponseWriter, r *http.Request) {
	render(w, "pages/success.html", nil)

}
