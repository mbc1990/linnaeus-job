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
	ImageDir       string
}

// Format of linnaeus response
type Class struct {
	ClassId     string // "n00000"
	ClassName   string // "corgi, pembrook welsh"
	Probability string // "0.76"
}

// Goes through unclassified images and has Linnaeus classify them
func (lj *LinnaeusJob) ClassifyImages() {
	unclassified := lj.Pg.GetUnclassified()

	port := strconv.Itoa(lj.ClassifierPort)
	for _, img := range unclassified {
		fmt.Println(*img)
		fname := lj.ImageDir + img.Filename
		url := fmt.Sprintf("http://localhost:%s/classify?filename=%s", port, fname)
		resp, err := http.Get(url)
		if err != nil {
			log.Print(err)
		}
		body, err := ioutil.ReadAll(resp.Body)

		// Close immediately rather than defer because we're in a big loop
		resp.Body.Close()
		classes := make([]Class, 0)
		json.Unmarshal(body, &classes)

		// Until we're returning probability from classifier, just use top class
		if len(classes) >= 1 {
			class := classes[0]
			prob, _ := strconv.ParseFloat(class.Probability, 32)
			lj.Pg.SaveClassification(img.ImageId, class.ClassId,
				class.ClassName, prob)
		}
	}
}

func NewLinnaeusJob(conf *Configuration) *LinnaeusJob {
	job := new(LinnaeusJob)
	job.Conf = conf
	job.ClassifierPort = conf.LinnaeusPort
	job.ImageDir = conf.ImageDirAbsPath
	job.Pg = NewPostgresClient(job.Conf.PGHost, job.Conf.PGPort,
		job.Conf.PGUser, job.Conf.PGPassword, job.Conf.PGDbname)
	return job
}
