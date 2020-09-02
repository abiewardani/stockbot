package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func main() {
	// e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Stockbot")
	// })
	// e.Logger.Fatal(e.Start(":8003"))
	robot()
}

// AccountHistory ....
type AccountHistory struct {
	Current  float32
	Previous float32
	Average  float32
	Min      float32
	Max      float32
	MOS      float32
}

// Account ..
type Account struct {
	PBV  AccountHistory
	PER  AccountHistory
	ROA  AccountHistory
	EPS  AccountHistory
	BVPS AccountHistory
}

var notAllowed = map[string]bool{
	`BALANCE SHEET`:    true,
	`INCOME STATEMENT`: true,
	`RATIO`:            true,
}

func robot() {
	c := colly.NewCollector(
		colly.UserAgent(`Robotorial Scrapping`),
		colly.AllowedDomains("www.indopremier.com", "indopremier.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	c.OnHTML(`.table-responsive`, func(e *colly.HTMLElement) {
		arrHeader := []string{}
		isAnlz := false
		e.ForEach("table.table-fundamental tr", func(_ int, el *colly.HTMLElement) {
			for i := 3; i < 20; i++ {
				iStr := strconv.Itoa(i)
				txt := el.ChildText("th:nth-child(" + iStr + ")")
				if txt == "" {
					break
				}
				if strings.Contains(txt, `Anlz`) {
					isAnlz = true
				}
				arrHeader = append(arrHeader, txt)
			}
		})

		account := Account{}
		e.ForEach("table.table-fundamental tr", func(_ int, el *colly.HTMLElement) {
			typeAccount := el.ChildText("td:first-child")
			if !notAllowed[typeAccount] {
				tmpArr := []string{}
				for i := 0; i < len(arrHeader); i++ {
					iStr := strconv.Itoa(i + 2)
					tmpArr = append(tmpArr, el.ChildText("td:nth-child("+iStr+")"))
				}

				if typeAccount == `PER` {
					perArr := cleansingPercentage(tmpArr)
					account.PER = HistoryCalculation(perArr, isAnlz)
					fmt.Println(account.PER)
				}

				if typeAccount == `PBV` {
					pbvArr := cleansingPercentage(tmpArr)
					account.PBV = HistoryCalculation(pbvArr, isAnlz)
					fmt.Println(account.PBV)
				}

				if typeAccount == `BVPS` {
					BVPSArr := cleansingPercentage(tmpArr)
					account.BVPS = HistoryCalculation(BVPSArr, isAnlz)
					fmt.Println(account.BVPS)
				}
			}
		})

		//	fmt.Println(arrAccount)
		//	TODO, play with data here

	})

	c.Visit(`https://www.indopremier.com/module/saham/include/fundamental.php?code=ICBP`)
}

// cLeansingValue ...
func cLeansingValue(params []string) []float32 {
	var finalAccount float32
	var res []float32
	for _, val := range params {
		val = strings.Replace(val, ",", "", -1)
		if strings.Contains(val, "T") {
			val = strings.Replace(val, " T", "", -1)
			f, err := strconv.ParseFloat(val, 32)
			if err != nil {
				continue
			}
			finalAccount = float32(f) * 1000
			res = append(res, finalAccount)
		}

		if strings.Contains(val, "B") {
			val = strings.Replace(val, "B", "", -1)
			f, err := strconv.ParseFloat(val, 32)
			if err != nil {
				continue
			}
			finalAccount = float32(f) * 1
			res = append(res, finalAccount)
		}
	}

	return res
}

func cleansingPercentage(params []string) []float32 {
	var finalAccount float32
	var res []float32
	for _, val := range params {
		val = strings.Replace(val, ",", "", -1)
		val = strings.Replace(val, " x", "", -1)
		val = strings.Replace(val, " %", "", -1)
		f, err := strconv.ParseFloat(val, 32)
		if err != nil {
			continue
		}
		finalAccount = float32(f) * 1
		res = append(res, finalAccount)
	}

	return res
}

// HistoryCalculation ...
func HistoryCalculation(arrValueHistory []float32, isAnlz bool) AccountHistory {
	var accountHistory AccountHistory
	if len(arrValueHistory) < 2 {
		return accountHistory
	}

	if isAnlz { // quartal 1,2,3
		arrValueHistory = RemoveIndex(arrValueHistory, 0)
	} else { // quartal 4
		arrValueHistory = RemoveIndex(arrValueHistory, 1)
	}

	var totalValueHistory float32
	max := arrValueHistory[0]
	min := arrValueHistory[0]

	for _, val := range arrValueHistory {
		totalValueHistory += val
		// get min & max value
		if max < val {
			max = val
		}
		if min > val {
			min = val
		}
	}

	current := arrValueHistory[0]
	previous := arrValueHistory[1]
	average := totalValueHistory / float32(len(arrValueHistory))
	MOS := ((average - current) / average) * 100

	return AccountHistory{Current: current, Previous: previous, Average: average, Min: min, Max: max, MOS: MOS}
}

// RemoveIndex ...
func RemoveIndex(s []float32, index int) []float32 {
	return append(s[:index], s[index+1:]...)
}
