package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GeoData struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func GetJobNum(html string) int {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Fatal("Error opening document: ", err)
	}

	var v int

	doc.Find(".jobsearch-JobCountAndSortPane-jobCount").Each(func(i int, s *goquery.Selection) {
		value := s.Find("span").Text()
		digits := make([]string, 0)
		for _, ch := range value {
			if ch > 47 && ch < 58 {
				digits = append(digits, string(ch))
			}
		}
		filtered := strings.Join(digits, "")

		ret, err := strconv.Atoi(filtered)
		if err != nil {
			fmt.Println("Atoi failed for: ", filtered, err.Error())
		} else {
			v = ret
		}
	})

	return v
}

func GetJobUrlList(html string) ([]string, error) {
	var ret []string

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	doc.Find(".jcs-JobTitle").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		ret = append(ret, link)
	})

	return ret, nil
}

func GetOneJob(html string, link string) (Listing, error) {
	var l Listing

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return Listing{}, err
	}

	l.Link = link
	l.Title = doc.Find(".jobsearch-JobInfoHeader-title").Find("span").Text()
	c := doc.Find(".css-1ioi40n").AttrOr("aria-label", "")
	c = strings.Replace(c, " (opens in a new tab)", "", -1)
	l.Company = c
	l.CompanyURL = doc.Find(".css-1ioi40n").AttrOr("href", "")
	loc := doc.Find(".css-1ojh0uo").Text()
	loc = trimRepeat(loc)
	l.Location = loc

	gL, err := getGeolocation(loc)
	if err != nil {
		l.Lat = ""
		l.Lon = ""
	} else {
		l.Lat = gL.Lat
		l.Lon = gL.Lon
	}

	l.Description = doc.Find("#jobDescriptionText").Text()

	return l, nil
}

func trimRepeat(s string) string {
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[i+len(s)/2] {
			return s
		}
	}
	return s[0 : len(s)/2]
}

func getGeolocation(addr string) (GeoData, error) {
	url := "https://geocode.maps.co/search?q=%s&api_key=6660b6109539d568840582itg6230e9"
	addr = strings.ReplaceAll(addr, " ", "+")
	url = fmt.Sprintf(url, addr)

	res, err := http.Get(url)
	if err != nil {
		return GeoData{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return GeoData{}, err
	}

	var gData []GeoData
	err = json.Unmarshal(body, &gData)
	if err != nil {
		return GeoData{}, err
	}

	if len(gData) > 0 {
		return gData[0], nil
	}

	return GeoData{}, errors.New("")
}
