package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const TR_MAIN_URL = "http://twtraffic.tra.gov.tw/twrail/"

type Train struct {
	Train_Code           string
	Class_Code           string
	Begin_Code           string
	Begin_Name           string
	Begin_EName          string
	End_Code             string
	End_Name             string
	End_EName            string
	Over_Night           string
	Direction            string
	MainViaRoad          string
	Handicapped          string
	Package              string
	Dining               string
	TrainType            string
	From_Departure_Time  string
	To_Arrival_Time      string
	Fare                 int
	Comment              string
	Discount_Price_Adult string
	Discount_Begin_Date  string
	Discount_End_Date    string
	From_Ticket_Code     string
	To_Ticket_Code       string
	Everyday             string
	TicketLink           string
}

type Station struct {
	Station_Code    string
	City_Code       string
	Station_Name    string
	Station_EName   string
	Station_Order   int
	STN_TICKET_CODE string
	Station_Name_JP string
	Station_Name_KR string
	Station_Name_SC string
	IDCode          string
	TextValue       string
	CityCode        string
}

func getStationCode() (map[string]string, map[string][]string) {
	stationMap := make(map[string]string)
	stationList := make(map[string][]string)
	values := url.Values{
		"datatype": {"station"},
		"language": {"tw"}}
	resp, err := http.PostForm(TR_MAIN_URL+"Services/BaseDataServ.ashx", values)
	if nil != err {
		fmt.Println("errorination happened getting the response", err)
		return nil, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		fmt.Println("errorination happened reading the body", err)
		return nil, nil
	}
	stringBody := string(body)
	res := []Station{}
	_ = json.Unmarshal([]byte(stringBody), &res)
	for _, stationIns := range res {
		stationMap[stationIns.Station_Name] = stationIns.Station_Code
		stationMap[stationIns.Station_Code] = stationIns.Station_Name
		stationList[stationIns.CityCode] = append(stationList[stationIns.CityCode], stationIns.Station_Name+": "+stationIns.Station_Code)
	}
	return stationMap, stationList
}

func findMatchingTrains(values url.Values) []Train {
	resp, err := http.PostForm("http://twtraffic.tra.gov.tw/twrail/TW_SearchResult.aspx", values)

	if nil != err {
		fmt.Println("errorination happened getting the response", err)
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		fmt.Println("errorination happened reading the body", err)
		return nil
	}
	stringBody := string(body)
	stringBody = strings.Replace(stringBody, "\r\n", " ", -1)
	stringBody = strings.Replace(stringBody, " ", "", -1)

	r, _ := regexp.Compile("JSONData=(.*)\\]")
	stringBody = r.FindString(stringBody)
	stringBody = strings.Replace(stringBody, "JSONData=", "", -1)

	res := []Train{}
	_ = json.Unmarshal([]byte(stringBody), &res)

	return res
}
