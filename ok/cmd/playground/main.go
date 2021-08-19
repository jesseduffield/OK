package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/jesseduffield/OK/ok/interpreter"
)

func withTimeout(f func()) error {
	c := make(chan struct{}, 1)
	go func() {
		f()
		c <- struct{}{}
	}()

	select {
	case <-c:
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("Timed out (program must complete within 5 seconds)")
	}
}

func main() {
	http.HandleFunc("/run", func(w http.ResponseWriter, r *http.Request) {
		// TODO: fix this when setting live
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		err := withTimeout(
			func() { interpreter.Interpret(r.Body, w) },
		)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
	})

	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
