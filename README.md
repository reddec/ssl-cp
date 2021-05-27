## Features


### CA key encryption

### CA key non-exposed by-default

### REST API

### CLI API

TBD

### Different storage backend

- sqlite (default)
- postgres

### GitOps

Initialization from pre-defined configuration that could
be stored in SVC.

#### Pre-build Go client library

```go
package main

import (
	"context"
	
	"github.com/reddec/ssl-cp/api"
	"github.com/reddec/ssl-cp/api/client"
)

func main() {
	remote := client.New("https://ssl-cp.example.com")
	cert, err := remote.CreateCertificate(context.Background(), api.Subject{
		Name:    "my.example.com",
		Days:    365,
		Domains: []string{"my.example.com", "srv2.example.com"},
	})
	// ...
}

```

### Simple installation

Just binary. Click and run. No additional config by-default.

### Complete web UI

### Mobile friendly