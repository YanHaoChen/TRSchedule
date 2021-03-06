# 臺鐵時刻表(CLI)

因臺鐵網頁版的時刻表操作上有點麻煩，所以就使用 GO，試著讓一切的操作簡潔一點。此專案寫的方式有點草率，但還算堪用，如有更好的寫法，或者是有什麼問題，再歡迎各位大大提出、指教。

## 使用前準備

* 請在電腦安裝 GO
  * [官方安裝說明](https://golang.org/doc/install)
  * 或者直接問 Google 大神，有沒有更好的安裝方式

## 編譯

```bash
$ git clone https://github.com/YanHaoChen/TRSchedule.git
$ cd TRSchedule
$ go build .
```

## 使用教學

> 記得將編譯後，得到的`TRSchedule`放置環境路目錄中哦。

```bash
$ TRSchedule                                                                                                                       

查詢車站代號：

        $ TRSchedule -l

使用範例如下：

1. 查詢當天從臺中(1319)到新烏日(1324)的車班。

	TRSchedule -from=1319 -to=1324

或者：

	TRSchedule -from=臺中 -to=新烏日

2. 查詢指定日期從臺中(1319)到新烏日(1324)的車班。

	TRSchedule -date=2019-01-01 -from=1319 -to=1324

3. 查詢指定日期及出發時間（下午兩點至四點）從臺中(1319)到新烏日(1324)的車班。

	TRSchedule -date=2019-01-01 -start=1400 -end=1600 -from=臺中 -to=新烏日

4. 查詢指定日期及到達時間（下午兩點至四點）從臺中(1319)到新烏日(1324)的車班。(type 預設為 1 ，也就是查詢出發時間。)

	TRSchedule -date=2019-01-01 -start=1400 -end=1600 -type=2 -from=1319 -to=1324
```

## 範例

```bash
$ TRSchedule -from=1319 -to=1238 -start=1600 -end=1800
|  車種|    車次|          啟站|          終站|  發車時間|  到達時間|     行駛時間|    票價|
--------------------------------------------------------------------------------------------
|  自強|     129|          基隆|          潮州|     16:12|     18:42|   02小時30分|     469|
	附註: 每日行駛。
	訂位連結: http://railway.hinet.net/Foreign/TW/etno1.html?from_station=146&to_station=185&getin_date=2018/12/23&train_no=129
--------------------------------------------------------------------------------------------
|區間車|    3257|          后里|          潮州|     16:57|     21:07|   04小時10分|     301|
	附註: 每日行駛。
	無法訂位
--------------------------------------------------------------------------------------------
...
```

