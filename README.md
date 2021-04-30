# Cased Proxy

![Build Status](https://github.com/cased/cased-go-proxy/workflows/cased-go-proxy/badge.svg)

cased-go-proxy is a service you can deploy as a proxy to publish your Cased
audit events when you are unable to distribute the publish key associated
with your account.

## Standalone server

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy?template=https://github.com/cased/cased-go-proxy)

## Adding to your existing application

An alternative way to proxy Cased audit events is to include the same Go HTTP
handler used in the standalone server in your existing Go application or API.

**Server**

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cased/cased-go-proxy/handlers"
)

func main() {
	http.HandleFunc("/", handlers.AuditEvents)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

**Client**

Next, you'll need to configure the URL Cased SDKs publish audit events to. All
Cased SDKs support configuring the URL which audit events are published to using
the `CASED_PUBLISH_URL` environment variable.

```
CASED_PUBLISH_URL=https://cased-proxy.herokuapp.com go run main.go
```

Each Cased SDK has documentation on how to configure the publish URL
programmatically.

## Contributing

1. Fork it ( https://github.com/cased/cased-go-proxy/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
