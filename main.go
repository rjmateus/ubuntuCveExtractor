package main
//https://ubuntu.com/security/cve?offset=34170
import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)
//var count = 0
func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("ubuntu.com"),
	)

	count := 0
	// On every a element which has href attribute call callback
	c.OnHTML("a.p-pagination__link--next[href]", func(e *colly.HTMLElement) {
		if count > 0 {
			return
		}
		count ++
		link := e.Attr("href")
		c.Visit(e.Request.AbsoluteURL(link))
		// Print link
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	result := make([]CveData, 0)
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
			for _, value := range cveLinks.Nodes{
				href := ""
				for _, att := range value.Attr{
					if strings.Compare(att.Key, "href") == 0{
						href = att.Val
					}
				}
				fullPath := response.Request.AbsoluteURL(href)
				cveID := value.FirstChild.Data
				result = append(result, processCVEDetailsPage(fullPath, cveID))
			}
		})
	})

	c.Visit("https://ubuntu.com/security/cve")

	b, err := json.Marshal(&result)
	if err != nil {
		fmt.Println("failed to serialize response:", err)
	}else {
		file, err := os.Create("ubuntu_cve.json")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		bufferWritter := bufio.NewWriter(file)
		bufferWritter.Write(b)
	}
}