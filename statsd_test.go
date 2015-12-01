package negronistatsd

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/codegangsta/negroni"
	"github.com/stretchr/testify/assert"
)

type mockClient struct {
	wg                     *sync.WaitGroup
	counterName, timerName string
}

func (c *mockClient) Counter(sampleRate float32, bucket string, n ...int) {
	if c.wg != nil {
		defer c.wg.Done()
	}
	c.counterName = bucket
}
func (c *mockClient) Timing(sampleRate float32, bucket string, d ...time.Duration) {
	if c.wg != nil {
		defer c.wg.Done()
	}
	c.timerName = bucket
}
func (mockClient) Gauge(sampleRate float32, bucket string, value ...string) {}

func TestCounters(t *testing.T) {
	c := mockClient{}
	m := Middleware{client: &c, prefix: "prefix"}

	recorder := negroni.NewResponseWriter(httptest.NewRecorder())
	recorder.WriteHeader(200)
	m.countResponse(recorder)
	assert.Equal(t, "prefix.request.200", c.counterName)
}

func TestTimers(t *testing.T) {
	c := mockClient{}
	m := Middleware{client: &c, prefix: "prefix"}

	m.timeRequest(time.Now(), &http.Request{RequestURI: "/a/b/c"})
	assert.Equal(t, c.timerName, "prefix.a.b.c")
}

func TestMiddleware(t *testing.T) {
	var wg sync.WaitGroup

	m := Middleware{client: &mockClient{wg: &wg}}

	wg.Add(2)
	m.ServeHTTP(httptest.NewRecorder(), &http.Request{}, func(rw http.ResponseWriter, r *http.Request) {})
	wg.Wait() // Just waiting for them means that they are going to be called. This is our main "assertion"
}
