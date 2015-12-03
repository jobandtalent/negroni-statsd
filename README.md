negroni-statsd
==============

[![GoDoc](https://godoc.org/github.com/jobandtalent/negroni-statsd?status.svg)](http://godoc.org/github.com/jobandtalent/negroni-statsd)

[Statsd](https://github.com/etsy/statsd) middleware for Negroni.

At jobandtalent we use statsd for tracking a lot of event that happen with our service, this middleware will help you to track some of those that can be globally tracked. At the moment we are tracking two things:

- time spent in a request. We will track the time between the requests get to the middleware until it cames back. Please, check the [lavarel documentation](https://mattstauffer.co/blog/laravel-5.0-middleware-filter-style) to see how middlewares work on Negroni.

For example, a request to `/api/users` will be tracked (using jt.yourservice as prefix) with the key: `jt.yourservice.api.users`.

- status codes of the responses. This stat is just a counter of the status code received and it will be tracked as: `jt.yourservice.request.200`, or `201`, etc...

Take care of...
---------------

- All the requests sent to statsd are going to be sent using UDP so you can loss some of them in favour of performance.
- The requests are going to be sent asynchronously without any waiting group, or similar. This means that if you stop your service ungracefully it's possible that some of them were not sent yet.

Using it
--------

    func main() {
        statsdURI := "localhost:1234"
        prefix := "jt.yourservice"
        n := negroni.New(statsd.NewMiddleware(stastdURI, prefix))
    }
