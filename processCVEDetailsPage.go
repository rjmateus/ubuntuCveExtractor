package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)



func processCVEDetailsPage(url string, cveID string) CveData{
	fmt.Println("processCVEDetailsPage--> %s", url)
	cveData := CveData{CveID: cveID, Packages: make(map [string]PackageInfo, 0)}

	// Load the HTML document
	doc := GetPageForSearch(url)

	cveData.UbuntuCvePublishDate = doc.Find(".p-strip h1 ~ p strong").Text()
	// Find the review items
	doc.Find(".cve-table > tbody").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		currentPackage := ""
		s.Find("tr").Each(func(i int, selection *goquery.Selection) {
			tds := selection.Find("td")
			if len(tds.Nodes) == 3 {
				name := selection.Find("td:nth-child(1) > a").Text()
				currentPackage = name
				cveData.Packages[currentPackage] = PackageInfo{currentPackage, make(map[string]string), make(map[string]OsData, 0)}
				selection.Find("td:nth-child(1) > small > a").Each(func(i int, s *goquery.Selection) {
					val,_ := s.Attr("href")
					packagesChannel := make([]PackagesChannel, 0)
					if strings.Compare(s.Text(), "Ubuntu") == 0{
						packagesChannel = getPackageList(val)
					}
					cveData.Packages[currentPackage].OsPackages[s.Text()] = OsData{s.Text(), val, packagesChannel}
				})
				release := strings.TrimSpace(strings.Replace(selection.Find("td:nth-child(2)").Text(), "\n", " ", -1))
				status := strings.TrimSpace(strings.Replace(selection.Find("td:nth-child(3)").Text(), "\n", " ", -1))

				cveData.Packages[currentPackage].ReleaseStatus[release] = status
			} else if len(tds.Nodes) == 2 {
				release := strings.TrimSpace(strings.Replace(selection.Find("td:nth-child(1)").Text(), "\n", " ", -1))
				status := strings.TrimSpace(strings.Replace(selection.Find("td:nth-child(2)").Text(), "\n", " ", -1))
				cveData.Packages[currentPackage].ReleaseStatus[release] = status
			}
		})
	})
	return cveData
}

func getPackageList(val string) []PackagesChannel{
	fmt.Println("getPackageList--> %s", val)
	result :=  make([]PackagesChannel, 0)
	doc := GetPageForSearch(val)

	// Find the review items
	doc.Find("#psearchres").Each(func(i int, s *goquery.Selection) {
		h2 := s.Find("h2")
		if h2.Nodes == nil {
			return
		}
		searchResultsTile := h2.Nodes[0].FirstChild.Data
		if strings.Compare(strings.ToLower(searchResultsTile), "exact hits") == 0 {
			uls := s.Find("ul")
			if len(uls.Nodes) > 0 {
				ulSelect := goquery.NewDocumentFromNode(uls.Nodes[0])
				ulSelect.Find("li").Each(func(i int, s *goquery.Selection) {
					channel := s.Find("a.resultlink").Text()
					packages := make([]string, 0)
					s.Find("span.binaries a").Each(func(i int, s *goquery.Selection) {
						packages = append(packages, s.Text())
					})
					s.Find("ul:first li")
					result = append(result, PackagesChannel{channel, packages})
				})
			}
		}
	})
	return result
}

func GetPageForSearch(url string) *goquery.Document {
	res, err := http.Get(url)
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
	return doc
}