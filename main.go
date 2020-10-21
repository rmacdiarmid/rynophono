package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	gapikey       = os.Getenv("GAPIKEY")
	gfindplaceurl = "https://maps.googleapis.com/maps/api/place/findplacefromtext/json"
	client        = http.Client{}
)

type Borrower struct {
	Name   string
	Street string
	City   string
	State  string
	Zip    string
}

// ToTextQuery returns string formatted to google places api text query
func (b *Borrower) ToTextQuery() string {
	q := strings.Builder{}
	q.Grow(len(b.Street) + 25)
	q.WriteString(b.Street)
	q.WriteString(", ")
	q.WriteString(b.City)
	q.WriteString(" ")
	q.WriteString(b.State)
	q.WriteString(" ")
	q.WriteString(b.Zip)
	return q.String()
}

type PlaceResponse struct {
	Candidates []struct {
		PlaceID string `json:"place_id"`
	} `json:"candidates"`
	Status string `json:"status"`
}

func GetPlacesId(b *Borrower) (*PlaceResponse, error) {
	v := url.Values{}
	v.Set("key", gapikey)
	v.Set("input", b.ToTextQuery())
	v.Set("inputtype", "textquery")

	resp, err := client.Get(gfindplaceurl + "?" + v.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var pr PlaceResponse
	err = json.Unmarshal(body, &pr)
	if err != nil {
		return nil, err
	}
	//fmt.Println("response Body:", string(body))
	return &pr, nil
}

type PlaceDetailResponse struct {
	Result struct {
		FormattedPhoneNumber string `json:"formatted_phone_number"`
		Website              string `json:"website"`
	} `json:"result"`
	Status string `json:"status"`
}

func GetPlacesDetail(pr *PlaceResponse) (*PlaceDetailResponse, error) {
	v := url.Values{}
	v.Set("key", gapikey)
	v.Set("fields", "formatted_phone_number,website")
	v.Set("place_id", pr.Candidates[0].PlaceID)

	resp, err := client.Get("https://maps.googleapis.com/maps/api/place/details/json?" + v.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var pdr PlaceDetailResponse
	err = json.Unmarshal(body, &pdr)
	if err != nil {
		return nil, err
	}
	//fmt.Println("response Body:", string(body))
	return &pdr, nil
}

func AddPlaceToRecord(record []string, placeId, phone, website string) []string {
	record = append(record, placeId)
	record = append(record, phone)
	record = append(record, website)
	return record
}

func main() {
	return
}
