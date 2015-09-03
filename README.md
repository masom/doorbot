# doorbot

Doorbot API server

Similar to https://envoy.co

It maintains a list of people, doors, devices, authentications ( api keys ).

A device (ex: iPad ) would act as a receptionist and allow a visitor to "ping" someone they would like to meet.

This is my first Golang application that goes beyond the "hello world" example.

There is probably some terrible stuff in there ( read: not idiomatic ) but it pretty much works.

Released under the MIT license

#### Features
- Multi-tenancy using subdomains
- Database migrations using goose
- 12-Factor app ( Heroku and ENV config support)
- PostgreSQL database using Gorp
- Structured logging
- Tests
