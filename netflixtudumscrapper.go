// Package netflixtudumscrapper scrappes data from the netflix tudum site
package netflixtudumscrapper

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// NetflixData stores the data we parse from the html
type NetflixData struct {
	Rank      string
	Title     string
	Weeks     string
	Hours     string
	Runtime   string
	Views     string
	DateRange string
}

// ScrapeNetflix fetches the HTML from the URL and parses the movie data based on the detected table structure
func ScrapeNetflix(url string) ([]NetflixData, error) {
	// Fetch the URL
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Parse the HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	var movies []NetflixData
	var dateRange string

	// Check for the first table structure
	if doc.Find("table.w-full.text-sm").Length() != 0 {
		dateRange = doc.Find("div.px-3.text-xs").Text()
		doc.Find("table.w-full.text-sm").Find("tbody tr").Each(func(i int, s *goquery.Selection) {
			movie := NetflixData{
				Rank:      s.Find("td.tbl-cell-rank").Text(),
				Title:     s.Find("td.tbl-cell-name").Text(),
				Weeks:     s.Find("td.tbl-cell-weeks .wk-number").Text(),
				Hours:     s.Find("td.tbl-cell-hours").Text(),
				Runtime:   s.Find("td.tbl-cell-runtime").Text(),
				Views:     s.Find("td.tbl-cell-vhor").Text(),
				DateRange: dateRange,
			}
			movies = append(movies, movie)
		})
	} else if doc.Find("table.w-full.mx-auto").Length() != 0 {
		// Check for the second table structure
		doc.Find("table.w-full.mx-auto").Find("tr").Each(func(i int, s *goquery.Selection) {
			if i != 0 { // Skip the header row
				var weeks string
				if s.Find("td.tbl-cell-weeks .wk-number").Length() != 0 {
					weeks = s.Find("td.tbl-cell-weeks .wk-number").Text()
				}
				movie := NetflixData{
					Rank:      s.Find("td:first-child").Text(),
					Title:     s.Find("td:nth-child(2)").Text(),
					Weeks:     weeks, // Extract weeks if present, else empty string
					Hours:     s.Find("td:nth-child(3)").Text(),
					Runtime:   s.Find("td:nth-child(4)").Text(),
					Views:     s.Find("td:nth-child(5)").Text(),
					DateRange: "",
				}
				movies = append(movies, movie)
			}
		})
	} else {
		return nil, fmt.Errorf("no recognized table format found")
	}

	return movies, nil
}
