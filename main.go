package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// O
type job struct {
	id       string
	title    string
	location string
	company  string
	summary  string
}

// O
func main() {
	term := getTerm()
	handleScrape(term)
}

// O
func getTerm() string {
	var term string
	fmt.Println("Input for searching...")
	fmt.Scanln(&term)
	return term
}

// O
func handleScrape(t string) {
	term := strings.ToLower(t)
	jobs := scrapeJob(term)
	writeJobs(jobs)
}

// handleScrape
// O
func scrapeJob(term string) []job {
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

// scrapeJob
// O
func getPages(url string) int {
	pages := 0
	res, err := http.Get(url)
	checkError(err)
	checkStatusCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		total := s.Find("a").Length()
		pages = total
	})

	return pages
}

// scrapeJob
// O
func getPage(number int, mainChannel chan []job, url string) {
	var jobs []job
	pageURL := url + "&start=" + strconv.Itoa(number*50)
	fmt.Println("Scrapping Indeed: Page", number)
	res, err := http.Get(pageURL)
	checkError(err)
	checkStatusCode(res)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkError(err)

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

// getPage
// O
func extractJob(s *goquery.Selection, c chan job) {
	// id, title, location, company, summary
	id, _ := s.Attr("data-jk")
	title, _ := s.Find(".jobTitle>span").Attr("title")
	title = cleanString(title)
	location := s.Find(".companyLocation").Text()
	location = cleanString(location)
	company := cleanString(s.Find(".companyName").Text())
	summary := cleanString(s.Find(".job-snippet").Text())
	c <- job{id: "https://www.indeed.com/viewjob?jk=" + id, title: title, location: location, company: company, summary: summary}
}

// handleScrape
// O
func writeJobs(jobs []job) {
	file, err := os.Create("jobs.csv")
	defer file.Close()
	checkError(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"apply", "title", "location", "company", "summary"}
	writeErr := w.Write(headers)
	checkError(writeErr)

	for _, job := range jobs {
		jobCSV := []string{job.id, job.title, job.location, job.company, job.summary}
		writeErr := w.Write(jobCSV)
		checkError(writeErr)
	}
}

// extractJob
// O
func cleanString(toClean string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(toClean)), " ")
}

// getPages
// O
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// O
func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatal("Status Code:", res.StatusCode)
	}
}
