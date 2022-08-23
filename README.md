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

If the environment variable `SEQ_URL` is set, it will be used as the URL to a SEQ instance. Consequently, logs will also
be sent to SEQ for your convenience. If you have a specific SEQ API Token, you can set it using `SEQ_API_TOKEN`.

### logging/alert

The `alert` package allows to integrate any alerting system. It already contains a default system for sending alert to
PushOver (Android and iOS push notifications), but you need to set the tokens before it will be functional.

### http and web

Functions for dealing with HTTP requests.

- `web.UploadImageToAzureBlob()` — Takes the image from the user from the HTTP request, and uploads it to Azure Blob
  storage. If the file is of invalid size or type, and error will be immediately returned using `web.Error()`

### random

For generating random stuffs, including a fast random string gen.

### validation

Allows you to make validation easier with the pre-built `V` and `M` for Validator and Mold instances respectively.

- `ValidateJSONBody` and optionally returns a recreated bytes of the JSON from the structure
