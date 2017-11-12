package main

import "net/http"
import "encoding/json"
import "io/ioutil"
import "fmt"

type LinnaeusJob struct {
	ClassifierPort int
	Conf           *Configuration
	Pg             *PostgresClient
}

// Format of linnaeus response
type Class struct {
	ClassId       string  // "n00000"
	HumanReadable string  // "corgi, pembrook welsh"
	Probability   float64 // 0.76
}

// Goes through unclassified images and has Linnaeus classify them
func (lj *LinnaeusJob) ClassifyImages() {
	fmt.Println("Classifying images...")
	unclassified := Pg.GetUnclassified()

	for _, fname := range unclassified {
		url := fmt.Sprintf("http://localhost:%s?filename=%s", lj.ClassifierPort, fname)
		resp, err := http.Get(url)
		if err != nil {
			// TODO: probably *don't* want to panic here
			panic(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		classes := make([]Class, 0)
		json.Unmarshal(body, &classes)

		// TODO: Do stuff with classes
		fmt.Println(classes)
	}
}

func NewLinnaeusJob(conf Configuration) *LinnaeusJob {
	job := new(LinnaeusJob)
	job.ClassifierPort = conf.ClassifierPort
	job.pg = NewPostgresClient(job.Conf.PGHost, job.Conf.PGPort,
		job.Conf.PGUser, job.Conf.PGPassword, job.Conf.PGDbname)
	return job
}
