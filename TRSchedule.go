package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
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

	if len(os.Args) < 2 {
		showExample()
		return
	} else if len(os.Args) == 2 {
		isHelpOrList := os.Args[1]
		if isHelpOrList == "-h" || isHelpOrList == "--help" {
			showExample()
			return
		} else if isHelpOrList == "-l" {
			showStationList()
			return
		} else {
			showExample()
			return
		}
	}

	const RFC3339FullDate = "2006-01-02"
	const ReservationDate = "2006/01/02"

	t := time.Now()
	searchDate := flag.String("date", t.Format(RFC3339FullDate), "Date")
	fromStation := flag.String("from", "None", "From station")
	toStation := flag.String("to", "None", "To station")
	start := flag.String("start", "0000", "Start of time")
	end := flag.String("end", "2359", "End of time")
	timeType := flag.String("type", "1", "Start(1) or arriving(2)")
	flag.Parse()

	// Get codes of stations
	stationCodeMap := getStationCode()

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

	resp, err := http.PostForm("http://twtraffic.tra.gov.tw/twrail/TW_SearchResult.aspx", values)

	if nil != err {
		fmt.Println("errorination happened getting the response", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if nil != err {
		fmt.Println("errorination happened reading the body", err)
		return
	}
	stringBody := string(body)
	stringBody = strings.Replace(stringBody, "\r\n", " ", -1)
	stringBody = strings.Replace(stringBody, " ", "", -1)

	r, _ := regexp.Compile("JSONData=(.*)\\]")
	stringBody = r.FindString(stringBody)
	stringBody = strings.Replace(stringBody, "JSONData=", "", -1)

	res := []Train{}
	_ = json.Unmarshal([]byte(stringBody), &res)

	if len(res) < 1 {
		fmt.Println("沒有對應的火車。")
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
					fmt.Printf("\t無法訂位\n")
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

func showStationList() {
	locationList :=
		"\n臺北/基隆地區(0):" +
			"\n\t 福隆:    1810, 貢寮:    1809, 雙溪:    1808" +
			"\n\t 牡丹:    1807, 三貂嶺:  1806, 猴硐:    1805" +
			"\n\t 瑞芳:    1804, 四腳亭:  1803, 暖暖:    1802" +
			"\n\t 基隆:    1801, 三坑:    1029, 八堵:    1002" +
			"\n\t 七堵:    1003, 百福:    1030, 五堵:    1004" +
			"\n\t 汐止:    1005, 汐科:    1031, 南港:    1006" +
			"\n\t 松山:    1007, 臺北:    1008, 萬華:    1009" +
			"\n\t 板橋:    1011, 浮洲:    1032, 樹林:    1012" +
			"\n\t 南樹林:  1034, 山佳:    1013, 鶯歌:    1014" +
			"\n\n桃園地區(1):" +
			"\n\t 桃園:    1015, 內壢:    1016, 中壢:    1017" +
			"\n\t 埔心:    1018, 楊梅:    1019, 富岡:    1020" +
			"\n\t 新富:    1036" +
			"\n\n新竹地區(2):" +
			"\n\t 北湖:    1033, 湖口:    1021, 新豐:    1022" +
			"\n\t 竹北:    1023, 北新竹:  1024, 新竹:    1025" +
			"\n\t 三姓橋:  1035, 香山:    1026" +
			"\n\n苗栗地區(3):" +
			"\n\t 崎頂:    1027, 竹南:    1028, 談文:    1102" +
			"\n\t 大山:    1104, 後龍:    1105, 龍港:    1106" +
			"\n\t 白沙屯:  1107, 新埔:    1108, 通霄:    1109" +
			"\n\t 苑裡:    1110, 造橋:    1302, 豐富:    1304" +
			"\n\t 苗栗:    1305, 南勢:    1307, 銅鑼:    1308" +
			"\n\t 三義:    1310" +
			"\n\n臺中地區(4):" +
			"\n\t 日南:    1111, 大甲:    1112, 臺中港:  1113" +
			"\n\t 清水:    1114, 沙鹿:    1115, 龍井:    1116" +
			"\n\t 大肚:    1117, 追分:    1118, 泰安:    1314" +
			"\n\t 后里:    1315, 豐原:    1317, 栗林:    1325" +
			"\n\t 潭子:    1318, 頭家厝:  1326, 松竹:    1327" +
			"\n\t 太原:    1323, 精武:    1328, 臺中:    1319" +
			"\n\t 五權:    1329, 大慶:    1322, 烏日:    1320" +
			"\n\t 新烏日:  1324, 成功:    1321" +
			"\n\n彰化地區(5):" +
			"\n\t 彰化:    1120, 花壇:    1202, 大村:    1240" +
			"\n\t 員林:    1203, 永靖:    1204, 社頭:    1205" +
			"\n\t 田中:    1206, 二水:    1207" +
			"\n\n南投地區(6):" +
			"\n\t 源泉:    2702, 濁水:    2703, 龍泉:    2704" +
			"\n\t 集集:    2705, 水里:    2706, 車埕:    2707" +
			"\n\n雲林地區(7):" +
			"\n\t 林內:    1208, 石榴:    1209, 斗六:    1210" +
			"\n\t 斗南:    1211, 石龜:    1212" +
			"\n\n嘉義地區(8):" +
			"\n\t 大林:    1213, 民雄:    1214, 嘉北:    1241" +
			"\n\t 嘉義:    1215, 水上:    1217, 南靖:    1218" +
			"\n\n臺南地區(9):" +
			"\n\t 後壁:    1219, 新營:    1220, 柳營:    1221" +
			"\n\t 林鳳營:  1222, 隆田:    1223, 拔林:    1224" +
			"\n\t 善化:    1225, 南科:    1244, 新市:    1226" +
			"\n\t 永康:    1227, 大橋:    1239, 臺南:    1228" +
			"\n\t 保安:    1229, 仁德:    1243, 中洲:    1230" +
			"\n\t 長榮大學:5101, 沙崙:    5102" +
			"\n\n高雄地區(10):" +
			"\n\t 大湖:    1231, 路竹:    1232, 岡山:    1233" +
			"\n\t 橋頭:    1234, 楠梓:    1235, 新左營:  1242" +
			"\n\t 左營:    1236, 內惟:    1245, 美術館:  1246" +
			"\n\t 鼓山:    1237, 三塊厝:  1247, 高雄:    1238" +
			"\n\t 民族:    1419, 科工館:  1420, 正義:    1421" +
			"\n\t 鳳山:    1402, 後庄:    1403, 九曲堂:  1404" +
			"\n\n屏東地區(11):" +
			"\n\t 六塊厝:  1405, 屏東:    1406, 歸來:    1407" +
			"\n\t 麟洛:    1408, 西勢:    1409, 竹田:    1410" +
			"\n\t 潮州:    1411, 崁頂:    1412, 南州:    1413" +
			"\n\t 鎮安:    1414, 林邊:    1415, 佳冬:    1416" +
			"\n\t 東海:    1417, 枋寮:    1418, 加祿:    1502" +
			"\n\t 內獅:    1503, 枋山:    1504" +
			"\n\n臺東地區(12):" +
			"\n\t 大武:    1508, 瀧溪:    1510, 金崙:    1512" +
			"\n\t 太麻里:  1514, 知本:    1516, 康樂:    1517" +
			"\n\t 臺東:    1632, 山里:    1631, 鹿野:    1630" +
			"\n\t 瑞源:    1629, 瑞和:    1628, 關山:    1626" +
			"\n\t 海端:    1625, 池上:    1624" +
			"\n\n花蓮地區(13):" +
			"\n\t 富里:    1623, 東竹:    1622, 東里:    1621" +
			"\n\t 玉里:    1619, 三民:    1617, 瑞穗:    1616" +
			"\n\t 富源:    1614, 大富:    1613, 光復:    1612" +
			"\n\t 萬榮:    1611, 鳳林:    1610,  南平:   1609" +
			"\n\t 林榮新光:1608, 豐田:    1607, 壽豐:    1606" +
			"\n\t 平和:    1605, 志學:    1604, 吉安:    1602" +
			"\n\t 花蓮:    1715, 北埔:    1714, 景美:    1713" +
			"\n\t 新城:    1712, 崇德:    1711, 和仁:    1710" +
			"\n\t 和平:    1709" +
			"\n\n宜蘭地區(14):" +
			"\n\t 漢本:    1708, 武塔:    1706, 南澳:    1705" +
			"\n\t 東澳:    1704, 永樂:    1703, 蘇澳:    1827" +
			"\n\t 蘇澳新:  1826, 新馬:    1825, 冬山:    1824" +
			"\n\t 羅東:    1823, 中里:    1822, 二結:    1821" +
			"\n\t 宜蘭:    1820, 四城:    1819, 礁溪:    1818" +
			"\n\t 頂埔:    1817, 頭城:    1816, 外澳:    1815" +
			"\n\t 龜山:    1814, 大溪:    1813, 大里:    1812" +
			"\n\t 石城:    1811" +
			"\n\n平溪/深澳線(15):" +
			"\n\t 瑞芳:    1804, 猴硐:    1805, 三貂嶺:  1806" +
			"\n\t 菁桐:    1908, 平溪:    1907, 嶺腳:    1906" +
			"\n\t 望古:    1905, 十分:    1904, 大華:    1903" +
			"\n\t 海科館:  6103, 八斗子:   2003" +
			"\n\n內灣/六家線(16):" +
			"\n\t 新竹:    1025, 北新竹:  1024, 千甲:    2212" +
			"\n\t 新莊:    2213, 竹中:    2203, 六家:    2214" +
			"\n\t 上員:    2204, 榮華:    2211, 竹東:    2205" +
			"\n\t 橫山:    2206, 九讚頭:  2207, 合興:    2208" +
			"\n\t 富貴:    2209, 內灣:    2210" +
			"\n\n集集線(17):" +
			"\n\t 二水:    1207, 源泉:    2702, 濁水:    2703" +
			"\n\t 龍泉:    2704, 集集:    2705, 水里:    2706" +
			"\n\t 車埕:    2707" +
			"\n\n沙崙線(18):" +
			"\n\t 中洲:    1230, 長榮大學:5101, 沙崙:    5102"

	fmt.Println(locationList)

}
