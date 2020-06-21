package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type DatapointProvider struct{}

type param struct {
	Name	string		`json:"name"`
	Units	string		`json:"units"`
	Text	string		`json:"$"`
}

type wx struct {
	Param	[]param		`json:"Param"`
}

type rep struct {
	D		string		`json:"D"`
	Gn		string		`json:"Gn"`
	Hn		string		`json:"Hn"`
	PPd		string		`json:"PPd"`
	S		string		`json:"S"`
	V		string		`json:"V"`
	Dm		string		`json:"Dm"`
	FDm		string		`json:"FDm"`
	W		string		`json:"W"`
	U		string		`json:"U"`
	Text	string		`json:"$"`
}

type period struct {
	Type	string		`json:"type"`
	Value	string		`json:"value"`
	Rep		[]rep		`json:"Rep"`
}

type location struct {
	I			string		`json:"i"`
	Lat			string		`json:"lat"`
	Lon			string		`json:"lon"`
	Name		string		`json:"name"`
	Country		string		`json:"country"`
	Continent	string		`json:"continent"`
	Elevation	string		`json:"elevation"`
	Period		[]period	`json:"Period"`
}

type dv struct {
	DataDate	string		`json:"dataDate"`
	Type		string		`json:"type"`
	Location	location	`json:"Location"`
}

type siteRep struct {
	Wx	wx		`json:"Wx"`
	DV	dv		`json:"DV"`
}

type siteRepMessage struct {
	SiteRep		siteRep		`json:"SiteRep"`
}


type paragraph struct {
	Title	string		`json:"title"`
	Text	string		`json:"$"`
}

type paragraphList struct {
	Paragraph	[]paragraph
}

// Some funky custom JSON unmarshalling
// Because a regPeriod can contain either a single paragraph or
// and array of them.
func (pList *paragraphList) UnmarshalJSON(data []byte) error {
	var p paragraph
	if err := json.Unmarshal(data, &p); err == nil {
		// It's a single paragraph
		pList.Paragraph = make([]paragraph, 1)
		pList.Paragraph[0] = p
	} else {
		// It's an array of the buggers
		err := json.Unmarshal(data, &(pList.Paragraph))
		if err != nil {
			return err
		}
	}
	return nil
}

type regPeriod struct {
	Id			string		`json:"id"`
	Paragraphs	paragraphList	`json:"Paragraph"`
}

type fcstPeriod struct {
	Period		[]regPeriod	`json:"Period"`
}

type regionalFcst struct {
	CreatedOn	string		`json:"createdOn"`
	IssuedAt	string		`json:"issuedAt"`
	RegionId	string		`json:"regionId"`
	FcstPeriods	fcstPeriod	`json:"FcstPeriods"`
}

type regionalFcstMessage struct {
	RegionalFcst	regionalFcst	`json:"RegionalFcst"`
}

// Obviously this is not a great idea.
// But the API requires this key in the clear. And this is a commandline
// app that runs on a 3rd party computer. What to do?
const apiKey = "ca899d51-c28b-49d2-a664-53264953d263"
const baseUrl = "http://datapoint.metoffice.gov.uk/public/data"
const locationEwell = 351409
const regionSE = 514

const (
	ResolutionDaily = "daily"
	Resolution3Hourly = "3hourly"
)

func init() {
	AddProvider("datapoint", DatapointProvider{})
}

func (d DatapointProvider) OneDayForecast(dayDelta int) OneDayForecast {

	_, err := dailyForecast(locationEwell, Resolution3Hourly)
	if err != nil {
		log.Fatalln(err)
	}

	regional, err := regionalForecast(regionSE)
	if err != nil {
		log.Fatalln(err)
	}

	var period0 = regional.RegionalFcst.FcstPeriods.Period[0]
	var description = fmt.Sprintf("%s\n\n%s\n%s\n\n%s\n%s",
		period0.Paragraphs.Paragraph[0].Text,
		period0.Paragraphs.Paragraph[1].Title,
		period0.Paragraphs.Paragraph[1].Text,
		period0.Paragraphs.Paragraph[2].Title,
		period0.Paragraphs.Paragraph[2].Text,
	)

	var f = OneDayForecast{
		Description: description,
	}
	return f
}

func dailyForecast(locationId int, resolution string) (siteRepMessage, error) {
	url := fmt.Sprintf("%s/%s/%s/%d?res=%s&key=%s",
		baseUrl, "val/wxfcs/all", "json", locationId, resolution, apiKey)

	body, err := httpGet(url)
	if err != nil {
		return siteRepMessage{}, err
	}

	var m siteRepMessage
	err = json.Unmarshal(body, &m)
	if err != nil {
		return siteRepMessage{}, err
	}

	return m, nil
}

func regionalForecast(regionId int) (regionalFcstMessage, error) {
	url := fmt.Sprintf("%s/%s/%s/%d?key=%s",
		baseUrl, "txt/wxfcs/regionalforecast", "json", regionId, apiKey)

	body, err := httpGet(url)
	if err != nil {
		return regionalFcstMessage{}, err
	}

	var m regionalFcstMessage
	err = json.Unmarshal(body, &m)
	if err != nil {
		return regionalFcstMessage{}, err
	}

	return m, nil
}

func httpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("datapoint: unexpected status code %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}
