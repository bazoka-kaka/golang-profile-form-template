package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

type User struct {
	Username string `json:"username"`
	Online   bool   `json:"online"`
}

func main() {
	http.HandleFunc("/", handleText)
	http.HandleFunc("/index", handleHtml)
	http.HandleFunc("/users", handleJson)

	// upload form
	http.HandleFunc("/form", handleForm)
	http.HandleFunc("/submit", handleSubmit)
	http.HandleFunc("/upload", handleUpload)

	// server responses
	http.HandleFunc("/route", handleRequest)

	fmt.Println("server running on port 3000")
	http.ListenAndServe(":3000", nil)
}

func handleText(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Success!"))
}

func handleHtml(w http.ResponseWriter, r *http.Request) {
	tmpl := `
		<h1>Benzion Yehezkel</h1>
		<p>My hobbies:</p>
		<ul>
			<li>Programming</li>
			<li>Gaming</li>
			<li>Anime</li>
		</ul>
	`
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

func handleJson(w http.ResponseWriter, r *http.Request) {
	users := []User{
		{"Benzion", false},
		{"Benjamin", true},
		{"Deborah", false},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	err := json.NewEncoder(w).Encode(&users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		tmpl, err := template.ParseFiles(filepath.Join("views", "index.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		firstName := r.FormValue("firstname")
		lastName := r.FormValue("lastname")
		onlineStr := r.FormValue("online")

		data := map[string]interface{}{
			"first_name": firstName,
			"last_name":  lastName,
			"online":     onlineStr,
		}

		tmpl, err := template.ParseFiles(filepath.Join("views", "success.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	http.Error(w, "", http.StatusBadRequest)
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(1024); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	alias := r.FormValue("alias")
	uploadedFile, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filename := handler.Filename
	if alias != "" {
		filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
	}

	targetFile, err := os.OpenFile(filepath.Join(dir, "uploads", filename), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte("File Uploaded!"))
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// method, url, param, header values, header get, body
	// method
	method := r.Method

	// url, param query
	url := r.URL
	urlParam := r.URL.Query().Get("param")

	// header values and get
	headerValues := r.Header.Values("foo")
	headerGet := r.Header.Get("foo")

	// body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writer := fmt.Sprintf("Method:\t\t %v \nUrl:\t\t %v \nUrl param:\t %v \nHeaderValues:\t %v \nHeaderGet:\t %v \nBody:\t\t %v", method, url, urlParam, headerValues, headerGet, string(body))
	w.Write([]byte(writer))
}
