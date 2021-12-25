package main

import (
	"strings"

	"github.com/crazybirdz/go-jobs/internals"
)

func main() {
	term := internals.GetTerm()
	term = strings.ToLower(term)
	jobs := internals.ScrapeJob(term)
	internals.WriteJobs(term, jobs)
}
