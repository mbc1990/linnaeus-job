package main

import "net/http"
import "encoding/json"
import "io/ioutil"
import "fmt"
import "strconv"
import "log"

type LinnaeusJob struct {
	ClassifierPort int
	Conf           *Configuration
	Pg             *PostgresClient
}

// Format of linnaeus response
type Class struct {
	ClassId     string  // "n00000"
	ClassName   string  // "corgi, pembrook welsh"
	Probability float64 // 0.76
}

// Goes through unclassified images and has Linnaeus classify them
func (lj *LinnaeusJob) ClassifyImages() {
	unclassified := lj.Pg.GetUnclassified()

	port := strconv.Itoa(lj.ClassifierPort)
	for _, img := range unclassified {
		fmt.Println(*img)
		url := fmt.Sprintf("http://localhost:%s/classify?filename=%s", port, img.Filename)
		resp, err := http.Get(url)
		if err != nil {
			log.Print(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		classes := make([]Class, 0)
		json.Unmarshal(body, &classes)

		// TODO: Do stuff with classes
		fmt.Println(classes)
	}
}

func NewLinnaeusJob(conf *Configuration) *LinnaeusJob {
	job := new(LinnaeusJob)
	job.Conf = conf
	job.ClassifierPort = conf.LinnaeusPort
	job.Pg = NewPostgresClient(job.Conf.PGHost, job.Conf.PGPort,
		job.Conf.PGUser, job.Conf.PGPassword, job.Conf.PGDbname)
	return job
}
