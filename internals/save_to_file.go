package internals

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/crazybirdz/go-jobs/tools"
)

func WriteJobs(term string, jobs []job) {
	var path string = "jobs/" + time.Now().Format("2006-January-02") + "_" + term + "_jobs.csv"
	file, err := os.Create(path)
	defer file.Close()
	defer fmt.Printf("Total %d jobs.", len(jobs))
	tools.CheckError(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"apply", "title", "location", "company", "summary"}
	writeErr := w.Write(headers)
	tools.CheckError(writeErr)

	for _, job := range jobs {
		jobCSV := []string{job.id, job.title, job.location, job.company, job.summary}
		writeErr := w.Write(jobCSV)
		tools.CheckError(writeErr)
	}
}
