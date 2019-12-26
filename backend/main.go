package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
)

var jar, _ = cookiejar.New(nil)

var client = &http.Client{
	Jar: jar,
}

func main() {
	err := loginCbs("9802089251", "k3EfVSamW&W8F")
	if err != nil {
		log.Fatal(err)
	}
	fetchCsbWidget("aptuslogin@APTUSPORT")
}

func loginCbs(username, password string) error {
	url := "https://www.chalmersstudentbostader.se/wp-login.php"
	rb := struct {
		Log        string `json:"log"`
		Pwd        string `json:"pwd"`
		RedirectTo string `json:"redirect_to"`
	}{
		Log:        username,
		Pwd:        password,
		RedirectTo: "https://www.chalmersstudentbostader.se/min-bostad/",
	}

	j, err := json.Marshal(rb)
	if err != nil {
		return err
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func fetchCsbWidget(widgetName string) error {
	url := "https://www.chalmersstudentbostader.se/widgets/"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	q := req.URL.Query()
	q.Add("callback", "jQuery")
	q.Add("widgets[]", widgetName)
	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println(resp.Body)
	return nil
}
