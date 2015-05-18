package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/influxdb/influxdb/client"
)

const (
	MyHost        = "localhost"
	MyPort        = 8086
	MyDB          = "sonar"
	MyMeasurement = "loadTime"
)

func main() {
	rand.Seed(int64(time.Now().Second()))
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", MyHost, MyPort))
	if err != nil {
		log.Fatal(err)
	}

	conf := client.Config{
		URL:      *u,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
	}

	con, err := client.NewClient(conf)
	if err != nil {
		log.Fatal(err)
	}
	dur, ver, err := con.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Happy as a Hippo! %v, %s", dur, ver)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)
	tick := time.Tick(500 * time.Millisecond)
outer:
	for {
		select {
		case <-tick:
			writePoints(con)
		case <-sigs:
			break outer
		}
	}
	printMeasures()
}

func writePoints(con *client.Client) {
	var (
		countries  = []string{"Iran", "China", "Syria", "Philippines"}
		versions   = []string{"1.5.17", "2.0.0-beta4", "2.0.0-beta5"}
		sites      = []string{"Twitter", "Facebook", "Google", "Baidu"}
		sampleSize = 10000
		pts        = make([]client.Point, sampleSize)
	)

	for i := 0; i < sampleSize; i++ {
		pts[i] = client.Point{
			Name: "sonar",
			Tags: map[string]string{
				"site":    sites[rand.Intn(len(sites))],
				"version": versions[rand.Intn(len(versions))],
				"country": countries[rand.Intn(len(countries))],
			},
			Fields: map[string]interface{}{
				"loadTime": time.Duration((10000 + rand.Intn(1000))) * time.Millisecond,
			},
			Time:      time.Now(),
			Precision: "s",
		}
	}
	bps := client.BatchPoints{
		Points:          pts,
		Database:        MyDB,
		RetentionPolicy: "default",
	}
	start := time.Now()
	_, err := con.Write(bps)
	if err != nil {
		log.Fatalf("Write error: %s", err)
	}
	log.Printf("Write %d data points, takes %v", sampleSize, time.Now().Sub(start))
}

func printMeasures() {
}
