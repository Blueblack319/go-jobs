package internals

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/crazybirdz/go-jobs/tools"
)

type job struct {
	id       string
	title    string
	location string
	company  string
	summary  string
}

func ScrapeJob(term string) []job {
	baseURL := "https://kr.indeed.com/jobs?q=" + term + "&limit=50"
	var jobs []job
	c := make(chan []job)

	totalPages := getPages(baseURL)
	fmt.Println("Extracted", totalPages, "pages")

	for i := 0; i < totalPages; i++ {
		go getPage(i, c, baseURL)
	}

	for i := 0; i < totalPages; i++ {
		pageJobs := <-c
		jobs = append(jobs, pageJobs...)
	}
	return jobs
}

func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
	tools.CheckError(err)
	tools.CheckStatusCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	tools.CheckError(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		total := s.Find("a").Length()
		pages = total
	})

	return pages
}

func getPage(number int, mainChannel chan []job, url string) {
	var jobs []job
	pageURL := url + "&start=" + strconv.Itoa(number*50)
	fmt.Println("Scrapping Indeed: Page", number)
	res, err := http.Get(pageURL)
	tools.CheckError(err)
	tools.CheckStatusCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	tools.CheckError(err)

	innerChannel := make(chan job)

	searchCards := doc.Find(".result")

	searchCards.Each(func(index int, s *goquery.Selection) {
		go extractJob(s, innerChannel)
	})

	for i := 0; i < searchCards.Length(); i++ {
		extracted := <-innerChannel
		jobs = append(jobs, extracted)
	}
	mainChannel <- jobs
}

func extractJob(s *goquery.Selection, c chan job) {
	// id, title, location, company, summary
	id, _ := s.Attr("data-jk")
	title, _ := s.Find(".jobTitle>span").Attr("title")
	title = tools.CleanString(title)
	location := s.Find(".companyLocation").Text()
	location = tools.CleanString(location)
	company := tools.CleanString(s.Find(".companyName").Text())
	summary := tools.CleanString(s.Find(".job-snippet").Text())
	c <- job{id: "https://www.indeed.com/viewjob?jk=" + id, title: title, location: location, company: company, summary: summary}
}
