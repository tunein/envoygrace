package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	var (
		envoyTarget string
		evictionPeriodMS int
		envoyRequestTimeoutMS int
		preSleepMS int
	)
	flag.StringVar(&envoyTarget, "envoy-url", "localhost:8001", "the base url where envoy lives")
	flag.IntVar(&evictionPeriodMS, "eviction-period", 5000, "amount of milliseconds to sleep after calling envoy graceful commands")
	flag.IntVar(&envoyRequestTimeoutMS, "envoy-timeout", 1000, "envoy request timeout in milliseconds")
	flag.IntVar(&preSleepMS, "pre-sleep", 5000, "how long to sleep before sending requests to envoy in milliseconds")
	flag.Parse()

	<-time.After(time.Duration(preSleepMS) * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(envoyRequestTimeoutMS)*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("http://%s/healthcheck/fail", envoyTarget), nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("response: %s", string(b))
	fmt.Printf("status code: %d", res.StatusCode)

	<-time.After(time.Duration(evictionPeriodMS) * time.Millisecond)
}
