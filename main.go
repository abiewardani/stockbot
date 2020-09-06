package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var listStocks = []string{
	`ICBP`, `ASII`, `INDF`, `UNVR`, `HMSP`, `MYOR`, // consumer
	`KLBF`,                         // obat2an
	`CTRA`, `PWON`, `BSDE`, `SMGR`, // konstruksi
	`CPIN`, `INTP`, `TLKM`, `SRIL`, // lain2
	`ITMG`, `PGAS`, `UNTR`, // cyclic
	`BBCA`, `BBRI`, `BTPS`, `BRIS`, // bank
}

var growthStocks = []string{
	`INDY`, `ADRO`, `PTBA`, `MBAP`,
	`SMGR`, `ULTJ`, `EKAD`, `SIDO`,
	`WOOD`, `MYOH`, `BEEF`, `ERAA`,
	`ACES`,
}

func main() {
	// e := echo.New()
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Stockbot")
	// })
	// e.Logger.Fatal(e.Start(":8003"))

	fmt.Println(`LQ45`)
	for _, val := range listStocks {
		robot(val, `lq45`)
	}
	fmt.Println(`========================`)
	fmt.Println(`Second Linier`)
	for _, val := range growthStocks {
		robot(val, `growth`)
	}
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
	PBV             AccountHistory
	PER             AccountHistory
	ROA             AccountHistory
	EPS             AccountHistory
	BVPS            AccountHistory
	ROE             AccountHistory
	Equity          AccountHistory
	OperatingProfit AccountHistory
	DebEquityRatio  AccountHistory
}

var notAllowed = map[string]bool{
	`BALANCE SHEET`:    true,
	`INCOME STATEMENT`: true,
	`RATIO`:            true,
}

func robot(stock string, typeStock string) {
	c := colly.NewCollector(
		colly.UserAgent(`Robotorial Scrapping`),
		colly.AllowedDomains("www.indopremier.com", "indopremier.com"),
	)

	c.OnRequest(func(r *colly.Request) {
		// fmt.Println(``)
		// log.Println("visiting", r.URL.String())
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

				switch typeAccount {
				case `PER`:
					perArr := cleansingPercentage(tmpArr)
					account.PER = HistoryCalculation(perArr, isAnlz)
				case `PBV`:
					pbvArr := cleansingPercentage(tmpArr)
					account.PBV = HistoryCalculation(pbvArr, isAnlz)
				case `BVPS`:
					BVPSArr := cleansingPercentage(tmpArr)
					account.BVPS = HistoryCalculation(BVPSArr, isAnlz)
				case `EPS`:
					EPSArr := cleansingPercentage(tmpArr)
					account.EPS = HistoryCalculation(EPSArr, isAnlz)
				case `ROE`:
					roeArr := cleansingPercentage(tmpArr)
					account.ROE = HistoryCalculation(roeArr, isAnlz)
				case `Total Equity`:
					eqArr := cLeansingValue(tmpArr)
					account.Equity = HistoryCalculation(eqArr, isAnlz)
				case `Operating Profit`:
					opProfitArr := cLeansingValue(tmpArr)
					account.OperatingProfit = HistoryCalculation(opProfitArr, isAnlz)
				case `Debt/Equity`:
					opDebEquityArr := cleansingPercentage(tmpArr)
					account.DebEquityRatio = HistoryCalculation(opDebEquityArr, isAnlz)
				}
			}
		})

		// Score
		// If PBV more than intrinsik value give it 2
		// If PBV more than intrinsik and MOS more than 10 - 20 % give it 1
		// If PBV more than intrinsik and MOS more than 20 % give it 2

		// fairPricePBV := account.PBV.Average * account.BVPS.Current
		// fairPricePER := account.PER.Average * account.EPS.Current

		// fmt.Println(`Fair Price PBV : `, FormatFloat(fairPricePBV))
		// fmt.Println(`Margin PBV : `, FormatFloat(account.PBV.MOS))
		// fmt.Println(``)
		// fmt.Println(`Fair Price PER : `, FormatFloat(fairPricePER))
		// fmt.Println(`Margin PER : `, FormatFloat(account.PER.MOS))

		// fmt.Println(``)
		// fmt.Println(`Current ROE :`, FormatFloat(account.ROE.Current))
		// fmt.Println(`Previous ROE :`, FormatFloat(account.ROE.Previous))
		// fmt.Println(`Margin ROE : `, FormatFloat(account.ROE.MOS))

		switch typeStock {
		case `lq45`:
			scoring := ScoringLQ45(account)
			if scoring >= 13 {
				fmt.Println(`Stock : `, stock, ` | `, `Scoring : `, scoring)
			}
		case `growth`:
			scoring := scoringGrowthStock(account)
			if scoring >= 13 {
				fmt.Println(`Stock : `, stock, ` | `, `Scoring : `, scoring)
			}
		}

		//	TODO, play with data here

	})

	c.Visit(`https://www.indopremier.com/module/saham/include/fundamental.php?code=` + stock)
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

// FormatFloat ...
func FormatFloat(aFloat float32) string {
	s := fmt.Sprintf("%.2f", aFloat)
	return s
}

// ScoringLQ45 ...
func ScoringLQ45(account Account) int {
	scoring := 0

	// PBV
	if account.PBV.MOS < 0 {
		scoring -= 3
	}

	if account.PBV.MOS > 10 && account.PBV.MOS <= 20 {
		scoring += 5
	}

	if account.PBV.MOS > 20 {
		scoring += 6
	}

	// PER ...
	if account.PER.MOS < 0 {
		scoring -= 2
	}

	if account.PER.MOS > 10 && account.PER.MOS <= 20 {
		scoring += 4
	}

	if account.PER.MOS > 20 {
		scoring += 5
	}

	// ROE ...
	if (account.ROE.MOS * -1) > 0 {
		scoring++
	}

	if account.ROE.Current > account.ROE.Previous {
		scoring += 3
	}

	if (account.ROE.Average) > 10 {
		scoring++
	}

	// Equity
	if account.Equity.Current > account.Equity.Previous {
		scoring++
	}

	// Operating Profit
	if account.OperatingProfit.Current > account.OperatingProfit.Previous {
		scoring += 3
	}

	// If its good company, Max Score 19
	return scoring
}

func scoringGrowthStock(account Account) int {
	scoring := 0

	// No More Than 1.5 PBV
	if account.PBV.Current >= 1.5 {
		return -99
	}

	// Konsisten ROE ...
	if (account.ROE.Current) < 4 {
		return -99
	}

	// DER no more than 1 ...
	if (account.DebEquityRatio.Current) > 1 {
		return -99
	}

	// PBV
	if account.PBV.MOS < 0 {
		scoring -= 3
	}

	if account.PBV.Current < 1 {
		scoring += 3
	}

	if account.PBV.MOS > 5 && account.PBV.MOS <= 10 {
		scoring += 5
	}

	if account.PBV.MOS > 10 {
		scoring += 6
	}

	// PER ...
	if account.PER.MOS < 0 {
		scoring -= 2
	}

	if account.PER.MOS > 5 && account.PER.MOS <= 10 {
		scoring += 4
	}

	if account.PER.MOS > 10 {
		scoring += 5
	}

	// ROE ...
	if (account.ROE.MOS * -1) > 0 {
		scoring++
	}

	if account.ROE.Current > account.ROA.Previous {
		scoring += 3
	}

	// Equity
	if account.Equity.Current > account.Equity.Previous {
		scoring++
	}

	// Operating Profit
	if account.OperatingProfit.Current > account.OperatingProfit.Previous {
		scoring += 3
	}

	// If its good company, Max Score 19
	return scoring
}
