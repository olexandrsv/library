package main

import (
	"fmt"
	"library/error_handler"
	"library/log"
	"library/request"
	"library/trace"
	"net/http"
)

func main() {
	testRequest()
}

func testRequest() {
	log.Init()

	http.HandleFunc("/logs", log.Endpoint)
	http.HandleFunc("/traces", log.TraceEndpoint)
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	trace.Save(w, r)
	function1(1, 2)
	req := request.New(r)
	n := req.Form.GetInt("n").Do()
	m := req.Form.GetInt("m").Do()
	s := req.Form.GetString("s").Do()

	err := req.Err()
	if err != nil {
		fmt.Printf("type: %T", err)
		error_handler.HandleError(err, w)
		return
	}

	fmt.Println(n, m, s)
}

func function1(n, m int){
	trace.Save(n, m)
}