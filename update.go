package update

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type Update struct {
	Host    string `json:"Host"`
	User    string `json:"User"`
	Subject string `json:"Subject"`
	Body    string `json:"Body"`
}

/* Update Receivers */
func (update Update) Notify() {
	subject := fmt.Sprintf("%s (%s@%s)", update.Subject, update.User, update.Host)
	body := "\n" + update.Body
	cmd := exec.Command("notify-send", subject, body)
	cmd.Run()
}

func (update Update) SendRequest(server string) {
	url := "http://" + server

	// set host if not specified
	if update.Host == "" {
		update.Host, _ = os.Hostname()
	}

	// set user if not specified
	if update.User == "" {
		update.User = os.Getenv("USER")
	}

	// package up update into json
	body, err := json.Marshal(&update)
	if err != nil {
		log.Fatal(err)
	}

	// create request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return
		//log.Fatal(err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	// send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
		//log.Fatal(err.Error())
	}
	defer resp.Body.Close()

	// read body
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
		//log.Fatal(err.Error())
	}
}
