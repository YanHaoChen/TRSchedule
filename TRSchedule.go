package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

func main() {

	trainClass := map[string]string{
		"1100": "  自強",
		"1101": "  自強",
		"1103": "  自強",
		"1108": "  自強",
		"1109": "  自強",
		"110A": "  自強",
		"110B": "  自強",
		"110C": "  自強",
		"110D": "  自強",
		"110E": "  自強",
		"110F": "  自強",
		"1102": "太魯閣",
		"1107": "普悠瑪",
		"1110": "  莒光",
		"1111": "  莒光",
		"1114": "  莒光",
		"1115": "  莒光",
		"1120": "  復興",
		"1130": "  電車",
		"1131": "區間車",
		"1132": "區間快",
		"1140": "普快車",
		"1141": "柴快車",
		"1151": "普通車"}

	const RFC3339FullDate = "2006-01-02"
	const ReservationDate = "2006/01/02"

	// Get codes of stations
	stationCodeMap, stationList := getStationCode()

	if len(os.Args) < 2 {
		showExample()
		return
	} else if len(os.Args) == 2 {
		isHelpOrList := os.Args[1]
		if isHelpOrList == "-h" || isHelpOrList == "--help" {
			showExample()
			return
		} else if isHelpOrList == "-l" {
			showStationList(stationList)
			return
		} else {
			showExample()
			return
		}
	}

	t := time.Now()
	searchDate := flag.String("date", t.Format(RFC3339FullDate), "Date")
	fromStation := flag.String("from", "None", "From station")
	toStation := flag.String("to", "None", "To station")
	start := flag.String("start", "0000", "Start of time")
	end := flag.String("end", "2359", "End of time")
	timeType := flag.String("type", "1", "Start(1) or arriving(2)")
	flag.Parse()

	if stationCodeMap == nil {
		fmt.Println("無法取得車站代號。")
		return
	}

	if *fromStation == "None" || *toStation == "None" {
		fmt.Println("請輸入起迄站。")
		return
	}

	codeOfFrom, foundFrom := stationCodeMap[*fromStation]
	codeOfTo, foundTo := stationCodeMap[*toStation]

	if !foundFrom {
		println("起站名稱有誤。")
		return
	} else {
		_, err := strconv.Atoi(codeOfFrom)
		if err != nil {
			codeOfFrom = *fromStation
		}
	}
	if !foundTo {
		println("迄站名稱有誤。")
		return
	} else {
		_, err := strconv.Atoi(codeOfTo)
		if err != nil {
			codeOfTo = *toStation
		}
	}

	startInt, err := strconv.Atoi(*start)
	if startInt > 2359 || startInt < 0 || err != nil {
		fmt.Println("起始時間有誤。")
		return
	}

	endInt, err := strconv.Atoi(*end)
	if endInt > 2359 || endInt < 0 || err != nil {
		fmt.Println("結束時間有誤。")
		return
	}

	reservationDate, err := time.Parse(RFC3339FullDate, *searchDate)
	if err != nil {
		fmt.Println("日期格式有誤。")
		return
	}

	values := url.Values{
		"FromStation":     {codeOfFrom},
		"FromStationName": {"0"},
		"ToStation":       {codeOfTo},
		"ToStationName":   {"0"},
		"TrainClass":      {"2"},
		"searchdate":      {*searchDate},
		"FromTimeSelect":  {*start},
		"ToTimeSelect":    {*end},
		"Timetype":        {*timeType}}

	res := findMatchingTrains(values)

	if res != nil {
		if len(res) < 1 {
			fmt.Println("沒有對應的火車。")
			return
		}
	} else {
		return
	}

	fmt.Printf("|%s|%6s|%12s|%12s|%6s|%6s|%9s|%6s|\n", "  車種", "車次", "啟站", "終站", "發車時間", "到達時間", "行駛時間", "票價")
	fmt.Printf("--------------------------------------------------------------------------------------------\n")
	for _, element := range res {
		closeCode := element.Class_Code
		trans_DT, _ := strconv.Atoi(element.From_Departure_Time)
		trans_AT, _ := strconv.Atoi(element.To_Arrival_Time)
		print_DT := element.From_Departure_Time[:2] + ":" + element.From_Departure_Time[2:]
		print_AT := element.To_Arrival_Time[:2] + ":" + element.To_Arrival_Time[2:]
		duration := 0

		if element.Over_Night == "1" {
			duration = (23-trans_DT/100)*60 + (60 - trans_DT%100) + (trans_AT/100)*60 + trans_AT%100
		} else {
			duration = ((trans_AT/100)*60 + trans_AT%100) - ((trans_DT/100)*60 + trans_DT%100)
		}
		transDuration := fmt.Sprintf("%02d小時%02d分", duration/60, duration%60)
		startStation := element.Begin_Name
		endStation := element.End_Name
		resultInterface := ""
		if len(startStation) == 6 && len(endStation) == 6 {
			resultInterface = "|%s|%8s|%12s|%12s|%10s|%10s|%10s|%8d|\n"
		} else if len(startStation) == 6 && len(endStation) == 9 {
			resultInterface = "|%s|%8s|%12s|%11s|%10s|%10s|%10s|%8d|\n"
		} else if len(startStation) == 6 && len(endStation) == 12 {
			resultInterface = "|%s|%8s|%12s|%10s|%10s|%10s|%10s|%8d|\n"
		} else if len(startStation) == 9 && len(endStation) == 6 {
			resultInterface = "|%s|%8s|%11s|%12s|%10s|%10s|%10s|%8d|\n"
		} else if len(startStation) == 9 && len(endStation) == 9 {
			resultInterface = "|%s|%8s|%11s|%11s|%10s|%10s|%10s|%8d|\n"
		} else if len(startStation) == 9 && len(endStation) == 12 {
			resultInterface = "|%s|%8s|%11s|%10s|%10s|%10s|%10s|%8d|\n"
		} else if len(startStation) == 12 && len(endStation) == 6 {
			resultInterface = "|%s|%8s|%10s|%12s|%10s|%10s|%10s|%8d|\n"
		} else if len(startStation) == 12 && len(endStation) == 9 {
			resultInterface = "|%s|%8s|%10s|%11s|%10s|%10s|%10s|%8d|\n"
		} else if len(startStation) == 12 && len(endStation) == 12 {
			resultInterface = "|%s|%8s|%10s|%10s|%10s|%10s|%10s|%8d|\n"
		}

		fmt.Printf(resultInterface,
			trainClass[element.Class_Code],
			element.Train_Code,
			startStation,
			endStation,
			print_DT,
			print_AT,
			transDuration,
			element.Fare)
		fmt.Printf("\t附註: %s\n", element.Comment)
		if element.TicketLink == "Y" {
			nowTime := t.Hour()*100 + t.Minute()
			if nowTime < trans_DT {
				if (closeCode == "1130") || (closeCode == "1131") || (closeCode == "1132") || (closeCode == "1140") || (closeCode == "1141") {
					fmt.Printf("\t無法訂位。\n")
				} else {
					fmt.Printf("\t訂位連結: http://railway.hinet.net/Foreign/TW/etno1.html?from_station=%s&to_station=%s&getin_date=%d/%02d/%02d&train_no=%s\n",
						element.From_Ticket_Code,
						element.To_Ticket_Code,
						reservationDate.Year(),
						reservationDate.Month(),
						reservationDate.Day(),
						element.Train_Code)
				}
			} else {
				fmt.Printf("\t無法訂位\n")
			}
		}
		fmt.Printf("--------------------------------------------------------------------------------------------\n")
	}

}

func showExample() {
	example := "\n查詢車站代號：\n" +
		"\n\tTRSchedule -l\n" +
		"\n使用範例如下：\n" +
		"\n1. 查詢當天從臺中(1319)到新烏日(1324)的車班。\n" +
		"\n\tTRSchedule -from=1319 -to=1324\n" +
		"\n或者：\n" +
		"\n\tTRSchedule -from=臺中 -to=新烏日\n" +
		"\n2. 查詢指定日期從臺中(1319)到新烏日(1324)的車班。\n" +
		"\n\tTRSchedule -date=2019-01-01 -from=1319 -to=1324\n" +
		"\n3. 查詢指定日期及出發時間（下午兩點至四點）從臺中(1319)到新烏日(1324)的車班。\n" +
		"\n\tTRSchedule -date=2019-01-01 -start=1400 -end=1600 -from=1319 -to=1324\n" +
		"\n4. 查詢指定日期及到達時間（下午兩點至四點）從臺中(1319)到新烏日(1324)的車班。(type 預設為 1 ，也就是查詢出發時間。)\n" +
		"\n\tTRSchedule -date=2019-01-01 -start=1400 -end=1600 -type=2 -from=1319 -to=1324\n"

	fmt.Println(example)
}

func showStationList(stationList map[string][]string) {
	areaList := []string{"臺北/基隆地區", "桃園地區", "新竹地區", "苗栗地區",
		"臺中地區", "彰化地區", "南投地區", "雲林地區", "嘉義地區",
		"臺南地區", "高雄地區", "屏東地區", "臺東地區", "花蓮地區",
		"宜蘭地區", "平溪/深澳線", "內灣/六家線", "集集線", "沙崙線"}
	for index, cityList := range stationList {
		cvIndex, _ := strconv.Atoi(index)
		println(areaList[cvIndex] + ":")
		for _, item := range cityList {
			println("\t" + item)
		}
		println("")
	}

}
