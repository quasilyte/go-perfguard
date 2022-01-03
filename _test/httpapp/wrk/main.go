package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func main() {
	numClients := 10
	requestsPerClient := 75000

	var wg sync.WaitGroup
	wg.Add(numClients)

	methodsFreq := map[string]int{
		"stats":              1,
		"jsonInspect":        4,
		"validateIdentError": 4,
		"objectsSum":         5,
		"validateIdent":      10,
		"formatDate":         15,
	}

	var methods []string
	for name, n := range methodsFreq {
		for i := 0; i < n; i++ {
			methods = append(methods, name)
		}
	}
	rand.Seed(time.Now().Unix())
	rand.Shuffle(len(methods), func(i, j int) {
		methods[i], methods[j] = methods[j], methods[i]
	})

	start := time.Now()
	getRequest("http://localhost:8080/startProfiling")
	for i := 0; i < numClients; i++ {
		go func() {
			for i := 0; i < requestsPerClient; i++ {
				method := methods[rand.Intn(len(methods))]
				switch method {
				case "stats":
					doStats()
				case "formatDate":
					doFormatDate()
				case "jsonInspect":
					doJsonInspect()
				case "validateIdent":
					doValidateIdent(true)
				case "validateIdentError":
					doValidateIdent(false)
				case "objectsSum":
					doObjectsSum()
				}
				delayRand := rand.Intn(150)
				if delayRand >= 140 {
					delayRand = 300
				}
				delay := delayRand + 50
				time.Sleep(time.Microsecond * time.Duration(delay))
			}
			wg.Done()
		}()
	}
	wg.Wait()
	getRequest("http://localhost:8080/finishProfiling")
	end := time.Since(start)

	fmt.Printf("made %d requests\n", numClients*requestsPerClient)
	fmt.Printf("elapsed: %.2fs\n", end.Seconds())
	fmt.Printf("requests per second: %d\n", int(float64(numClients*requestsPerClient)/end.Seconds()))
}

func getRequest(urlFormat string, args ...interface{}) {
	resp, err := http.Get(fmt.Sprintf(urlFormat, args...))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
}

func doStats() {
	getRequest("http://localhost:8080/stats")
}

func doFormatDate() {
	now := time.Now()
	getRequest("http://localhost:8080/formatDate?arg=%d", now.Unix())
}

func doJsonInspect() {
	data := url.QueryEscape(`{"f": 1, "str": "v", "arr": [1]}`)
	getRequest("http://localhost:8080/jsonInspect?arg=%s", data)
}

func doValidateIdent(valid bool) {
	if valid {
		getRequest("http://localhost:8080/validateIdent?arg=RelativelyLongIdent183ValueButNotTooLong")
	} else {
		getRequest("http://localhost:8080/validateIdent?arg=This_Is_InvalidThis_Is_InvalidThis_Is_InvalidThis_Is_InvalidThis_Is_InvalidThis_Is_Invalid")
	}
}

func doObjectsSum() {
	getRequest("http://localhost:8080/objectsSum")
}
