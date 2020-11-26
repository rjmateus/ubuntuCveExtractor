package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
	"log"
	"net/http"
	"strings"
)

type PackagesChannel struct {
	channel string
	packages []string
}

type OsData struct {
	Name, Url string
	PackagesChannelList []PackagesChannel
}
type PackageCve struct {
	PackageName     string
	OsVersionStatus map[string]string
	OsData          map[string]OsData
}


type CveData struct {
	CveID  string
	PublishData string
	PkgCve map[string]PackageCve
}

func processCVEDetailsPage(url string, cveID string) CveData{
	c := colly.NewCollector()

	cveData := CveData{CveID: cveID, PkgCve: make(map [string]PackageCve, 0)}

	c.OnResponse(func(response *colly.Response) {
		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(response.Body))
		if err != nil {
			log.Fatal(err)
		}

		cveData.PublishData = doc.Find(".p-strip h1 ~ p strong").Text()
		// Find the review items
		doc.Find(".cve-table > tbody").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			currentPackage := ""
			s.Find("tr").Each(func(i int, selection *goquery.Selection) {
				tds := selection.Find("td")
				if len(tds.Nodes) == 3 {
					name := selection.Find("td:nth-child(1) > a").Text()
					currentPackage = name
					cveData.PkgCve[currentPackage] = PackageCve{currentPackage, make(map[string]string), make(map[string]OsData, 0)}
					fmt.Println(name)
					selection.Find("td:nth-child(1) > small > a").Each(func(i int, s *goquery.Selection) {
						val,_ := s.Attr("href")
						fmt.Printf("system: %s --> %s\n", s.Text(), val)
						packagesChannel := make([]PackagesChannel, 0)
						if strings.Compare(s.Text(), "Ubuntu") == 0{
							packagesChannel = getPackageList(val)
						}
						cveData.PkgCve[currentPackage].OsData[s.Text()] = OsData{s.Text(), val, packagesChannel}
					})
					release := strings.TrimSpace(strings.Replace(selection.Find("td:nth-child(2)").Text(), "\n", " ", -1))
					status := strings.TrimSpace(strings.Replace(selection.Find("td:nth-child(3)").Text(), "\n", " ", -1))
					fmt.Printf("\t %s --> %s\n", release, status)

					cveData.PkgCve[currentPackage].OsVersionStatus[release] = status
				} else if len(tds.Nodes) == 2 {
					release := strings.TrimSpace(strings.Replace(selection.Find("td:nth-child(1)").Text(), "\n", " ", -1))
					status := strings.TrimSpace(strings.Replace(selection.Find("td:nth-child(2)").Text(), "\n", " ", -1))
					fmt.Printf("\t %s --> %s\n", release, status)
					cveData.PkgCve[currentPackage].OsVersionStatus[release] = status
				}
			})
		})
	})
	c.Visit(url)
	//fmt.Println(cveData)
	b, err := json.Marshal(&cveData)
	if err != nil {
		fmt.Println("failed to serialize response:", err)
	}else {
		fmt.Println(string(b))
	}
	return cveData
}

func getPackageList(val string) []PackagesChannel{
	result :=  make([]PackagesChannel, 0)
	res, err := http.Get(val)
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
	doc.Find("#psearchres li").Each(func(i int, s *goquery.Selection) {

		channel := s.Find("a.resultlink").Text()
		packages := make([]string, 0)
		s.Find("span.binaries a").Each(func(i int, s *goquery.Selection) {
			packages = append(packages, s.Text())
		})
		fmt.Printf("%s - %s\n", channel, packages)
		result = append(result, PackagesChannel{channel, packages})
	})
	return result
}