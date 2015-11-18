package negronistatsd

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/codegangsta/negroni"
	g2s "github.com/peterbourgon/g2s"
)

const (
	sampleRate float32 = 1.0

	counterKeyFormatter = "%s.request.%d"
	timeKeyFormatter    = "%s%s"
)

// Middleware stores the prefix and the statsd client
type Middleware struct {
	prefix string
	client g2s.Statter
}

// NewMiddleware returns a middleware struct with a configured statsd client
func NewMiddleware(uri, prefix string) *Middleware {
	c, err := g2s.Dial("udp", uri)
	if err != nil {
		log.Println("No statsd server on %s", uri)
		c = nopClient{}
	}
	return &Middleware{client: c, prefix: prefix}
}

func (m Middleware) timeRequest(startTime time.Time, r *http.Request) {
	name := fmt.Sprintf(timeKeyFormatter, m.prefix, strings.Replace(r.RequestURI, "/", ".", -1))
	m.client.Timing(sampleRate, name, time.Since(startTime))
}

func (m Middleware) countResponse(res negroni.ResponseWriter) {
	name := fmt.Sprintf(counterKeyFormatter, m.prefix, res.Status())
	m.client.Counter(sampleRate, name, 1)
}

func (m Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	res := negroni.NewResponseWriter(rw) // TODO: should we create our ResponseWriter wrapper to avoid this?

	defer func(startTime time.Time) {
		go m.timeRequest(startTime, r)
		go m.countResponse(res)
	}(time.Now())

	next(res, r)
}
