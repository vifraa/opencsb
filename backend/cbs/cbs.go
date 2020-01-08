package cbs

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var jar, _ = cookiejar.New(nil)

var client = &http.Client{
	Jar: jar,
}
var noRedirectClient = &http.Client{
	Jar: jar,
	CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

func LoginCbs(username, password string) error {
	u := "https://www.chalmersstudentbostader.se/wp-login.php"

	v := url.Values{}
	v.Add("log", username)
	v.Add("pwd", password)
	v.Add("redirect_to", "https://www.chalmersstudentbostader.se/min-bostad/")

	res, err := client.PostForm(u, v)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	correctLogin := cookieIsSet("Fast2User_ssoId", res.Request.URL)
	if correctLogin {
		return nil
	} else {
		return errors.New("invalid login")
	}
}

func LoginAptusPort() error {
	url, err := getAptusportLogin()
	if err != nil {
		return err
	}

	res, err := noRedirectClient.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	correctLogin := cookieIsSet(".ASPXAUTH", res.Request.URL)
	if correctLogin {
		return nil
	} else {
		return errors.New("invalid aptus login")
	}
}

func OpenDoor(id string) error {
	if id == "" {
		return errors.New("invalid door id")
	}
	url := "https://apt-www.chalmersstudentbostader.se/AptusPortal/Lock/UnlockEntryDoor/" + id
	res, err := noRedirectClient.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	b, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return errors.New("unable to open door" + string(b))
	}

	return nil
}

func cookieIsSet(name string, u *url.URL) bool {
	isSet := false
	for _, c := range client.Jar.Cookies(u) {
		if c.Name == name {
			isSet = true
		}

	}
	return isSet
}

func getAptusportLogin() (string, error) {
	widgetRes, err := fetchCsbWidget("aptuslogin@APTUSPORT")
	if err != nil {
		return "", err
	}

	url := widgetRes.Data.Aptuslogin_APTUSPORT.Objekt[0].AptusURL
	url = strings.ReplaceAll(url, " ", "+")

	return url, nil
}

func fetchCsbWidget(widgetName string) (CSBWidgetResponse, error) {
	var widgetRes CSBWidgetResponse

	u := "https://www.chalmersstudentbostader.se/widgets/"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return widgetRes, err
	}

	q := req.URL.Query()
	q.Add("callback", "jQuery")
	q.Add("widgets[]", widgetName)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return widgetRes, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return widgetRes, err
	}

	// Remove wrapping to be able to parse to json.
	s := string(data)
	s = s[7 : len(s)-2]

	err = json.Unmarshal([]byte(s), &widgetRes)
	if err != nil {
		return widgetRes, err
	}

	return widgetRes, nil
}

type CSBWidgetResponse struct {
	CSRFtoken string `json:"CSRFtoken"`
	Data      struct {
		Aptuslogin struct {
			Objekt []struct {
				AptusURL string `json:"aptusUrl"`
				ObjektNr string `json:"objektNr"`
			} `json:"objekt"`
		} `json:"aptuslogin"`
		Aptuslogin_APTUSPORT struct {
			Objekt []struct {
				AptusURL string `json:"aptusUrl"`
				ObjektNr string `json:"objektNr"`
			} `json:"objekt"`
		} `json:"aptuslogin@APTUSPORT"`
		Avisummering struct {
			AntalAviserade int `json:"antalAviserade"`
			AntalBetalda   int `json:"antalBetalda"`
			AntalObetalda  int `json:"antalObetalda"`
		} `json:"avisummering"`
		Avtalsummering []struct {
			Antal       int    `json:"antal"`
			Beskrivning string `json:"beskrivning"`
			Kategori    string `json:"kategori"`
		} `json:"avtalsummering"`
		Felanmalansummering []struct {
			Antal    int    `json:"antal"`
			CSSClass string `json:"cssClass"`
			Label    string `json:"label"`
		} `json:"felanmalansummering"`
		Intresseerbjudandesummering struct {
			AntalErbjudanden           string      `json:"antalErbjudanden"`
			AntalIntressen             string      `json:"antalIntressen"`
			AntalObesvaradeErbjudanden string      `json:"antalObesvaradeErbjudanden"`
			AntalPrelBokadeErbjudanden interface{} `json:"antalPrelBokadeErbjudanden"`
		} `json:"intresseerbjudandesummering"`
		Koerochprenumerationer_STD struct {
			Action             interface{} `json:"action"`
			Aktiv              bool        `json:"aktiv"`
			AktivPrenumeration bool        `json:"aktivPrenumeration"`
			Kodagar            int         `json:"kodagar"`
			Namn               string      `json:"namn"`
			RegistreradDatum   string      `json:"registreradDatum"`
			Rendered           bool        `json:"rendered"`
			Status             string      `json:"status"`
			StatusDetalj       interface{} `json:"statusDetalj"`
			TypKod             string      `json:"typKod"`
		} `json:"koerochprenumerationer@STD"`
		Kontaktuppgifter struct {
			Adress struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        string      `json:"value"`
			} `json:"adress"`
			BekraftaEpost struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    interface{} `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         interface{} `json:"type"`
				Value        interface{} `json:"value"`
			} `json:"bekraftaEpost"`
			Co struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        string      `json:"value"`
			} `json:"co"`
			Enamn struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        string      `json:"value"`
			} `json:"enamn"`
			Epost struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        string      `json:"value"`
			} `json:"epost"`
			ErbjudandeSms struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        interface{} `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    interface{} `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         interface{} `json:"type"`
				Value        string      `json:"value"`
			} `json:"erbjudandeSms"`
			Fnamn struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        string      `json:"value"`
			} `json:"fnamn"`
			Hyresreferenser struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        interface{} `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    interface{} `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         interface{} `json:"type"`
				Value        string      `json:"value"`
			} `json:"hyresreferenser"`
			KundItems               []interface{} `json:"kundItems"`
			KundnrPlaceholderNonSwe string        `json:"kundnrPlaceholderNonSwe"`
			KundnrPlaceholderSwe    string        `json:"kundnrPlaceholderSwe"`
			Kundnummer              struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        string      `json:"value"`
			} `json:"kundnummer"`
			Land struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        interface{} `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    interface{} `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         interface{} `json:"type"`
				Value        string      `json:"value"`
			} `json:"land"`
			Media struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        string      `json:"value"`
			} `json:"media"`
			Mobilnr struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        string      `json:"value"`
			} `json:"mobilnr"`
			Ort struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        string      `json:"value"`
			} `json:"ort"`
			Postnr struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        string      `json:"value"`
			} `json:"postnr"`
			Sprak struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        interface{} `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    interface{} `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         interface{} `json:"type"`
				Value        string      `json:"value"`
			} `json:"sprak"`
			SvPersonnr struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    interface{} `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         interface{} `json:"type"`
				Value        string      `json:"value"`
			} `json:"svPersonnr"`
			TelefonArbete struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        interface{} `json:"value"`
			} `json:"telefonArbete"`
			TelefonBostad struct {
				Helptext     interface{} `json:"helptext"`
				ID           interface{} `json:"id"`
				InputType    interface{} `json:"inputType"`
				Label        string      `json:"label"`
				Legend       interface{} `json:"legend"`
				Mandatory    string      `json:"mandatory"`
				Maxlength    interface{} `json:"maxlength"`
				Name         interface{} `json:"name"`
				Placeholder  string      `json:"placeholder"`
				Readonly     bool        `json:"readonly"`
				Rendered     bool        `json:"rendered"`
				Required     bool        `json:"required"`
				SelectValues struct{}    `json:"selectValues"`
				Type         string      `json:"type"`
				Value        interface{} `json:"value"`
			} `json:"telefonBostad"`
		} `json:"kontaktuppgifter"`
	} `json:"data"`
	Events     []interface{} `json:"events"`
	FileResult interface{}   `json:"fileResult"`
	HTML       struct {
		Alert                  string `json:"alert"`
		Avisummering           string `json:"avisummering"`
		Avtalsummering         string `json:"avtalsummering"`
		Felanmalansummering    string `json:"felanmalansummering"`
		Koerochprenumerationer string `json:"koerochprenumerationer"`
		Kontaktuppgifter       string `json:"kontaktuppgifter"`
	} `json:"html"`
	Javascripts []interface{} `json:"javascripts"`
	Messages    []interface{} `json:"messages"`
	OpenWindow  interface{}   `json:"openWindow"`
	RedirectURL interface{}   `json:"redirectUrl"`
	ReplaceURL  interface{}   `json:"replaceUrl"`
}
