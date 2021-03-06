package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

var (
	answer         = ""
	webRoot        = "web-root"
	updateFlagFile = "tmp/update"
)

func main() {
	loadAnswer()
	router := httprouter.New()
	router.POST("/answer", answerPost)
	router.NotFound = http.HandlerFunc(serveStaticFilesOr404)
	log.Fatal(http.ListenAndServe(":8089", router))
}

func loadAnswer() {
	loadedAnswer, error := ioutil.ReadFile("private/answer.txt")
	panicOnError(error)
	answer = string(loadedAnswer)
}

func answerPost(responseWriter http.ResponseWriter, request *http.Request, requestParameters httprouter.Params) {
	if answer == "" || answer != request.PostFormValue("answer") {
		fmt.Fprint(responseWriter, "false")
		return
	}
	_, error := os.Create(updateFlagFile)
	if error != nil {
		fmt.Fprint(responseWriter, "false")
		return
	}
	fmt.Fprint(responseWriter, "true")
}
