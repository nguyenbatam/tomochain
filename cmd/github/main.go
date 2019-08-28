package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Coin struct {
	code   string
	name   string
	total  string
	source string
}
type Release struct {
	time        int64
	header      string
	description string
	c           Coin
}

func GetGithubInfo(c *Coin) []Release {
	// Request the HTML page.
	res, err := http.Get(c.source + "/releases")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Find the review items
	releases := []Release{}
	doc.Find(".release-main-section ").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		releaseHeader := s.Find(".release-header .text-normal").Text()
		releaseHeader = strings.TrimSpace(releaseHeader)
		if len(releaseHeader) == 0 {
			return
		}
		date, _ := s.Find("relative-time").Attr("datetime")
		date = strings.TrimSpace(date)
		description := s.Find(".markdown-body").Text()
		description = strings.TrimSpace(description)
		t, err := time.Parse(time.RFC3339, date)
		if err != nil {
			fmt.Println(err, date, c.source+"/releases")

		}
		releases = append(releases, Release{t.Unix(), releaseHeader, description, Coin{c.code, c.name, c.total, c.source}})
	})
	return releases
}

func GetListCoin() []Coin {
	// Request the HTML page.
	res, err := http.Get("https://coinmarketcap.com/vi/coins/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	urls := []string{}
	// Find the review items
	doc.Find("a[class*=currency-name-container][href*=currencies]").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		link, _ := s.Attr("href")
		urls = append(urls, "https://coinmarketcap.com"+link)
	})

	coins := []Coin{}
	for _, url := range urls {
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		doc, err = goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		sourceCode, exit := doc.Find("a[href*=github]").Attr("href")
		if !exit {
			sourceCode, exit = doc.Find("a[href*=gitlab]").Attr("href")
			if !exit {
				sourceCode, exit = doc.Find("a[href*=bitbucket]").Attr("href")
			}
			if !exit {
				fmt.Println("not exit", url)
				continue
			}
		}
		s := doc.Find(".details-panel-item--name ")
		code := s.Find("span").Text()
		name, _ := s.Find("img").Attr("alt")
		name = strings.TrimSpace(name)
		total := doc.Find(".details-panel-item--marketcap-stats span[data-currency-market-cap] span[data-currency-value]").Text()
		total = strings.TrimSpace(total)
		res, err = http.Get(sourceCode)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()
		doc, err = goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal("err", err)
		}
		if doc.Find(".octicon-code").Size() == 0 {
			if doc.Find(".pinned-item-list-item-content a").Size() > 0 {
				sourceCode, _ = doc.Find(".pinned-item-list-item-content a").First().Attr("href")
			} else {
				continue
			}
			sourceCode = "https://github.com" + sourceCode
		}
		coins = append(coins, Coin{name, code, total, sourceCode})
	}
	return coins
}

func main() {
	coins := GetListCoin()
	releases := []Release{}
	for _, coin := range coins {
		releases = append(releases, GetGithubInfo(&coin)...)
	}
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].time > releases[j].time
	})
	fmt.Println(len(releases))
	f, err := os.Create("info.txt")
	if os.IsExist(err){
		os.Remove("info.txt")
		f, err = os.Create("info.txt")
		if err != nil{
			fmt.Println(err)
		}
	}
	for i := 0; i < len(releases); i++ {
		unixTimeUTC := time.Unix(releases[i].time, 0) //gives unix time stamp in utc
		f.WriteString("=========>>" + unixTimeUTC.Format(time.RFC3339) + "\n")
		f.WriteString(strconv.Itoa(i) + "\t" + releases[i].c.code + "\t" + releases[i].c.name + "\t" + releases[i].c.total + "\t" + releases[i].c.source + "\t" + releases[i].header + "\n")
		f.WriteString(releases[i].description + "\n")
	}
	f.Sync()
	f.Close()
}
