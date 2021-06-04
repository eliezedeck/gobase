# gobase

This project is a Golang library for making server apps as productive as possible. There are many packages in this
module, and you are free to use only what you want.

## Packages

Usually, you will use the library in a Go Module context.

### config

Typically, configurations are taken from environment variables. This makes changes in runtime easier, for example for
containerized apps (Docker, Kubernetes, ...).

The package already contains a few frequent configurations that you can use.

### database/mongodb

`mongodb.Access` provides a quick connect and access to MongoDB collections.

### logging

Zap (from Uber) is the default logger used in this library. This package provides an instance called `L` which you can
import in any context (for example as dot package) for quick `L.Info(...)`.

### logging/alert

The `alert` package allows to integrate any alerting system. It already contains a default system for sending alert to
PushOver (Android and iOS push notifications), but you need to set the tokens before it will be functional.

### http

Functions for dealing with HTTP requests.

### random

For generating random stuffs, including a fast random string gen.

### validation

Allows you to make validation easier with the pre-built `V` and `M` for Validator and Mold instances respectively.
