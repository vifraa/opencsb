package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var jar, _ = cookiejar.New(nil)

var client = &http.Client{
	Jar: jar,
}

func main() {
	err := loginCbs("9802089251", "k3EfVSamW&W8F^")
	if err != nil {
		log.Fatal(err)
	}
	fetchCsbWidget("aptuslogin@APTUSPORT")
}

func loginCbs(username, password string) error {
	u := "https://www.chalmersstudentbostader.se/wp-login.php"

	v := url.Values{}
	v.Add("log", username)
	v.Add("pwd", password)
	v.Add("redirect_to", "https://www.chalmersstudentbostader.se/min-bostad/")

	resp, err := client.PostForm(u, v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	correctLogin := loginCookieSet(resp.Request.URL)
	if correctLogin {
		return nil
	} else {
		return errors.New("invalid login")
	}
}

func loginCookieSet(u *url.URL) bool {
	isSet := false
	for _, c := range client.Jar.Cookies(u) {
		if c.Name == "Fast2User_ssoId" {
			isSet = true
		}

	}
	return isSet
}

func fetchCsbWidget(widgetName string) error {
	u := "https://www.chalmersstudentbostader.se/widgets/"
	req, err := http.NewRequest("GET", u, nil)
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

	b, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(b))

	return nil
}
