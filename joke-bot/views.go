package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type dadJokeResponse struct {
	ID     string `json:"id,omitempty"`
	Joke   string `json:"joke,omitempty"`
	Status int    `json:"status,omitempty"`
}

type response struct {
	Joke *string `json:"joke"`
}

func requestJSON(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.Print("response:\n")
	log.Print(resp)

	if resp.StatusCode != 200 {
		err := errors.New("Request to Dad Joke API failed")
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getDadJokeView(w http.ResponseWriter, r *http.Request) {
	url := "https://icanhazdadjoke.com/"

	body, err := requestJSON(url)
	if err != nil {
		se := serverError{}
		se.applyServerError(err, "Whoops, no jokes right now", w)
		return
	}

	jk := dadJokeResponse{}
	err = json.Unmarshal(body, &jk)
	if err != nil {
		log.Fatal(err)
	}
	data := response{
		Joke: &jk.Joke,
	}

	applyResponseJSON(200, data, w)
}

type healthCheck struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func healthCheckView(w http.ResponseWriter, r *http.Request) {
	data := healthCheck{
		Status:  "OK",
		Message: "Healthy",
	}
	applyResponseJSON(200, data, w)
}

func applyResponseJSON(code int, data interface{}, w http.ResponseWriter) {
	if code != 200 {
		w.WriteHeader(code)
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		se := serverError{}
		se.applyServerError(err, "Server error", w)
		return
	}
	w.Write(jsonData)
	w.Header().Set("Content-Type", "application/json")
}

type serverError struct {
	Err          error  `json:"error"`
	ErrorMessage string `json:"message"`
}

func (serverErr *serverError) applyServerError(err error, errorMessage string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "application/json")
	se := &serverError{
		Err:          err,
		ErrorMessage: errorMessage,
	}
	jsonData, err := json.Marshal(se)
	w.Write(jsonData)
}
