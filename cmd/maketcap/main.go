package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

func main() {
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
		github, _ := doc.Find("a[href*=github]").Attr("href")
		s := doc.Find(".details-panel-item--name ")
		code := s.Find("span").Text()
		coin, _ := s.Find("img").Attr("alt")
		coin = strings.TrimSpace(coin)
		fmt.Println(coin, "a", code,  github)
	}
}
