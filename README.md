negroni-statsd
==============

Statsd middleware for Negroni.

This middleware will track:

- time spent in a request.
- status codes of the responses.

Example
-------

    func main() {
        n := negroni.New(statsd.NewMiddleware("localhost:1234", "jobandtalent"))
    }
