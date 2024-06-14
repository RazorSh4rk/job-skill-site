package cmd

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"

	"razorsh4rk.github.io/indeedscrape/browser"
	"razorsh4rk.github.io/indeedscrape/scraper"
)

type Config struct {
	Term      string
	Location  string
	Pages     int
	SkipPages int
	Browser   string
}

type IConfig interface {
	Process()
}

func (cf *Config) Process() {
	cf.Term = strings.Replace(cf.Term, " ", "+", -1)
	cf.Location = strings.Replace(cf.Location, " ", "+", -1)
	if !slices.Contains([]string{"firefox", "chrome"}, cf.Browser) {
		cf.Browser = "firefox"
	}
}

func Process(c *Config, callback func(scraper.Listing)) {
	c.Process()

	url :=
		fmt.Sprintf("https://www.indeed.com/jobs?q=%s&l=%s", c.Term, c.Location)
	urlPager := "&start=%d"

	b := createBrowser(c)
	defer b.Destroy()

	numJobs := scraper.GetJobNum(b.GetPage(fmt.Sprintf(url, 0)))
	if numJobs == 0 {
		log.Printf("Found 0 listed jobs, please check your search on indeed dot com\n")
		os.Exit(0)
	}
	log.Printf("Found %d jobs\n", numJobs)

	jobLinks := make([]string, 0)

	for idx := range c.Pages {
		start := idx*10 + (c.SkipPages * 10)
		suffix := fmt.Sprintf(urlPager, start)
		pageUrl := url + suffix
		fmt.Printf("idx: %d, start: %d, skipped: %d, pageurl: %s\n", idx, start, c.SkipPages, pageUrl)

		page := b.GetPage(pageUrl)
		j, err := scraper.GetJobUrlList(page)
		if err == nil {
			jobLinks = append(jobLinks, j...)
		} else {
			fmt.Println("error: ", err.Error())
		}
	}

	fmt.Printf("Found %d jobs\n", len(jobLinks))

	jobs := make([]scraper.Listing, len(jobLinks))

	for idx, link := range jobLinks {
		fmt.Println("Job #", idx)

		j := addOneJob(c, jobs, idx, "https://indeed.com"+link)
		callback(j)
	}
}

func addOneJob(c *Config, l []scraper.Listing, idx int, link string) scraper.Listing {
	b := createBrowser(c)

	defer b.Destroy()

	page := b.GetPage(link)
	job, err := scraper.GetOneJob(page, link)
	if err == nil && !job.IsEmpty() {
		l[idx] = job
	}
	return job
}

func createBrowser(c *Config) browser.Browser {
	var b browser.Browser

	if c.Browser == "chrome" {
		b.FromChrome()
	} else {
		b.FromFirefox()
	}

	return b
}
