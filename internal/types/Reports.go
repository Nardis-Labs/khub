package types

import "time"

// Report represents the dumpfile s3 object metadata
type Report struct {
	Name         string    `json:"name"`
	LastModified time.Time `json:"lastModified"`
	Created      string    `json:"created"`
	Expires      string    `json:"expires"`
	Bucket       string    `json:"bucket"`
	Size         int64     `json:"size"`
	SizeUnits    string    `json:"sizeUnits"`
	User         string    `json:"user"`
	Reason       string    `json:"reason"`
	Type         string    `json:"type"`
}
