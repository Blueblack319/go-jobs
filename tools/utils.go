package tools

import (
	"log"
	"net/http"
	"strings"
)

func CleanString(toClean string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(toClean)), " ")
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatal("Status Code:", res.StatusCode)
	}
}
