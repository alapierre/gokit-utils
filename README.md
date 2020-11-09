# gokit-utils
Go Kit Utils for microservices

### Eureka registration

```go
package main

import (
	"github.com/alapierre/gokit-utils/eureka"
)

func main() {
    eurekaInstance, err := eureka.New().
        Default(8080, "api/schedule").
        Register("http://localhost:8761/eureka", "scheduler-service")
    
    if err != nil {
      panic(err)
    }
    
    defer eurekaInstance.Deregister()
}
```
