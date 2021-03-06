package goalpost

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"
)

//JobStatus is a enumerated int representing the processing status of a Job
type JobStatus int

// JobStatus types
const (
	Ack JobStatus = iota + 1
	Uack
	Nack
	Failed
)

type AuditMessage struct {
	Message string
	Time    time.Time
}

//Job wraps arbitrary data for processing
type Job struct {
	Status JobStatus
	//Unique identifier for a Job
	ID uint64

	Label string

	//Data contains the bytes that were pushed using Queue.PushBytes()
	Data []byte

	//RetryCount is the number of times the job has been retried
	//If your work can have a temporary failure state, it is recommended
	//that you check retry count and return a fatal error after a certain
	//number of retries
	RetryCount int

	//Message is primarily used for debugging. It conatains status info
	//about what was last done with the job.
	AuditLog []*AuditMessage
}

//DecodeJob decodes a gob encoded byte array into a Job struct and returns a pointer to it
func DecodeJob(b []byte) *Job {
	//TODO: this should return an error in the event decoder.Decode returns an err
	buffer := bytes.NewReader(b)
	decoder := gob.NewDecoder(buffer)
	job := Job{}
	decoder.Decode(&job)
	return &job
}

//Bytes returns a gob encoded byte array representation of *j
func (j *Job) Bytes() []byte {
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	encoder.Encode(j)
	return buffer.Bytes()
}

//RecoverableWorkerError defines an error that a worker DoWork func
//can return that indicates the message should be retried
type RecoverableWorkerError struct {
	message string
}

func (e RecoverableWorkerError) Error() string {
	return fmt.Sprintf("Worker encountered a temporary error: %s", e.message)
}

//NewRecoverableWorkerError creates a new RecoverableWorkerError
func NewRecoverableWorkerError(message string) RecoverableWorkerError {
	return RecoverableWorkerError{message}
}
