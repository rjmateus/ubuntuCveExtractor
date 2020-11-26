package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"

	"github.com/gocolly/colly"
)
var count = 0
func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("ubuntu.com"),
	)

	//count := 0
	// On every a element which has href attribute call callback
	c.OnHTML("a.p-pagination__link--next[href]", func(e *colly.HTMLElement) {
		if count > 1 {
			return
		}
		count ++
		link := e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
		// Print link
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnResponse(func(response *colly.Response) {
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body))
		if err != nil {
			log.Fatal(err)
		}

		// Find the review items
		doc.Find(".cve-table").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			cveLinks := s.Find("a[href]")
			for i, value := range cveLinks.Nodes{
				href := ""
				for _, att := range value.Attr{
					if strings.Compare(att.Key, "href") == 0{
						href = att.Val
					}
				}
				fullPath := response.Request.AbsoluteURL(href)
				fmt.Printf("Review %d: %s --> %s\n", i, href, fullPath)
				cveID := value.FirstChild.Data
				processCVEDetailsPage(fullPath, cveID)
				return
			}
		})
		fmt.Println()
	})


	// Start scraping on https://hackerspaces.org
	c.Visit("https://ubuntu.com/security/cve")
}