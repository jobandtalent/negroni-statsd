package negronistatsd

import "time"

type nopClient struct{}

func (nopClient) Counter(sampleRate float32, bucket string, n ...int)          {}
func (nopClient) Timing(sampleRate float32, bucket string, d ...time.Duration) {}
func (nopClient) Gauge(sampleRate float32, bucket string, value ...string)     {}
