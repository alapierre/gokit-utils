package eureka

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"os"
	"testing"
	"time"
)

func TestIntegration(t *testing.T) {

	logger := log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestamp)

	registrar1, err := New().
		Default(9099, "api/schedule").
		Register("http://localhost:8761/eureka", "scheduler-service")

	if err != nil {
		panic(err)
	}

	defer registrar1.Deregister()

	fmt.Println("Press 'Enter' to stop...")
	duration := time.Minute * 3
	time.Sleep(duration)
}

func TestIp(t *testing.T) {
	ip, _ := GetLocalIP()
	fmt.Println(ip)
}
