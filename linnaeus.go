package main

type LinnaeusJob struct {
	linnaeusPort int
}

// Goes through unclassified images and has Linnaeus classify them
func (lj *LinnaeusJob) ClassifyImages() {

}

func NewLinnaeusJob(port int) *LinnaeusJob {
	job := new(LinnaeusJob)
	job.linnaeusPort = port
	return job
}
