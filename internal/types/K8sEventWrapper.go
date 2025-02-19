package types

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
)

type K8sEventWrapper struct {
	Interval string `json:"interval"`
	Object   string `json:"object"`
	v1.Event
}

func GetInterval(e v1.Event) string {
	var interval string
	firstTimestampSince := translateMicroTimestampSince(e.EventTime)
	if e.EventTime.IsZero() {
		firstTimestampSince = translateTimestampSince(e.FirstTimestamp)
	}
	if e.Series != nil {
		interval = fmt.Sprintf("%s (x%d over %s)", translateMicroTimestampSince(e.Series.LastObservedTime), e.Series.Count, firstTimestampSince)
	} else if e.Count > 1 {
		interval = fmt.Sprintf("%s (x%d over %s)", translateTimestampSince(e.LastTimestamp), e.Count, firstTimestampSince)
	} else {
		interval = firstTimestampSince
	}

	return interval
}

// translateMicroTimestampSince returns the elapsed time since timestamp in
// human-readable approximation.
func translateMicroTimestampSince(timestamp metav1.MicroTime) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}

// translateTimestampSince returns the elapsed time since timestamp in
// human-readable approximation.
func translateTimestampSince(timestamp metav1.Time) string {
	if timestamp.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Since(timestamp.Time))
}
