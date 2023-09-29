# ippoippophotography-com-functions-contact

Structured Go code to provide functionality to a DigitalOcean serverless functions.

## Overview

DigitalOcean provides [Servless Functions](https://docs.digitalocean.com/products/functions/), which are also available via their [App Platform](https://docs.digitalocean.com/products/app-platform/) product.
Functions can be written in Go, which is my preferred choice for backend implementation.
The [expected source structure](https://docs.digitalocean.com/products/functions/how-to/structure-projects/) and [build process](https://docs.digitalocean.com/products/functions/reference/build-process/) seems to prevent a well structured Go program.

```text
example-project
├── packages
│   ├── example-package-1
│   │   ├── example-function-a.php
│   │   ├── example-function-b
│   │   │   ├── package.json
│   │   │   └── example.js
│   │   └── example-function-c
│   │       └── index.php
│   └── example-package-2
│       └── example-function
│           ├── requirements.txt
│           ├── __main__.py
│           └── example.py
└── project.yml
```

The precise requirements for a custom `build.sh` are not clear. Most samples seem to rely on a single `.go` file for each function. The `doctl serverless deploy` output showed that each `.go` source file was being compiled seperatly, and this sub-package imports were failing.
This was also causing problems with `*_test.go` files too.

Since I want to avoid a single `.go` file, I am extracting all functionality into a seperate repository per servless function.

### Solution

[My Website (monorepo) repository: ippoippophotography-com-src](https://github.com/ippoippo/ippoippophotography-com-src) source code will be updated to include `functions-src`.

```text
functions-src
├── packages
│   └── contact
│       ├── go.mod // `require`s `ippoippophotography-com-functions-contact`
│       └── contact.go (package main) // DigitalOcean function: Depends on ippoippophotography-com-functions-contact. Calls `contactform.Execute()`
├── packages-test
│   └── contact
│       ├── go.mod // `require`s `ippoippophotography-com-functions-contact`
│       └── main.go (package main) // Locally executable file, executed via `functions-contact-test-main-run` Task: Simulates the DigitalOcean function
└── project.yml // Configures package/contact to be deployed to DigitalOcean
```

[This repository: ippoippophotography-com-functions-contact](https://github.com/ippoippo/ippoippophotography-com-functions-contact) will be a versioned/tagged source code containing the actual structured implementation and unit tests.

```text
.
├── api
│   ├── api_test.go
│   └── api.go // Define the API (request/response) for the Serverless Function, and also the `contactform` `Execute()` function
├── configuration
│   ├── configuration_test.go
│   └── configuration.go // Load configuration for third party APIs, such SendGrid
├── contactform
│   ├── contactform_test.go
│   └── contactform.go // Main "executable", that is configured. Exposes an `Execute()` function to be called from the DigitalOcean function
├── mailer
│   ├── mailer_test.go
│   └── mailer.go // Constructs an email message from the contact form request, and calls SendGrid
├── validation
│   ├── validator_test.go
│   └── validator.go // Validates the request from DigitalOcean
└── go.mod
```

## Version naming convention

This repo use [SemVer](https://semver.org).

- MAJOR version when you make incompatible API changes
- MINOR version when you add functionality in a backward compatible manner
  - My addition: This will also include version upgrades of Go (eg. 1.20 to 1.21)
- PATCH version when you make backward compatible bug fixes
  - My addition: This will include library updates.
