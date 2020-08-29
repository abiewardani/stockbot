package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gocolly/colly"
)

func main() {
	robot()
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
		e.ForEach("table.table-fundamental tr", func(_ int, el *colly.HTMLElement) {
			for i := 3; i < 20; i++ {
				iStr := strconv.Itoa(i)
				txt := el.ChildText("th:nth-child(" + iStr + ")")
				if txt == "" {
					break
				}
				arrHeader = append(arrHeader, txt)
			}
		})

		arrAccount := map[string][]string{}
		e.ForEach("table.table-fundamental tr", func(_ int, el *colly.HTMLElement) {
			typeAccount := el.ChildText("td:first-child")
			if !notAllowed[typeAccount] {
				tmpArr := []string{}
				for i := 0; i < len(arrHeader); i++ {
					iStr := strconv.Itoa(i + 2)
					tmpArr = append(tmpArr, el.ChildText("td:nth-child("+iStr+")"))
				}
				arrAccount[el.ChildText("td:first-child")] = tmpArr
			}
		})

		fmt.Println(arrAccount)
		//	TODO, play with data here

	})

	c.Visit(`https://www.indopremier.com/module/saham/include/fundamental.php?code=pgas&quarter=1`)
}
