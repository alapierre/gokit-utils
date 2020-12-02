package eureka

import (
	"fmt"
	slog "github.com/go-eden/slf4go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd/eureka"
	"github.com/hudl/fargo"
	"net"
	"os"
)

type Client struct {
	HealthCheckUrl string
	StatusPageUrl  string
	HomePageUrl    string
	port           int

	ip       string
	hostname string
	logger   log.Logger
}

//goland:noinspection GoUnusedExportedFunction
func New() *Client {
	return &Client{
		HealthCheckUrl: "",
		StatusPageUrl:  "",
		HomePageUrl:    "",
	}
}

func (c *Client) Default(port int, homePage string) *Client {
	c.DefaultWithLogger(port, homePage, log.NewLogfmtLogger(os.Stderr))
	return c
}

func (c *Client) DefaultWithLogger(port int, homePage string, logger log.Logger) *Client {
	c.logger = logger
	err := c.init(port, homePage, "http")
	if err != nil {
		panic(err)
	}
	return c
}

func (c *Client) init(port int, homePage, proto string) error {

	c.logger = log.With(c.logger, "ts", log.DefaultTimestamp)

	c.port = port
	var err error
	c.ip, err = GetLocalIP()
	if err != nil {
		return err
	}

	c.hostname, err = os.Hostname()

	if err != nil {
		return err
	}

	slog.Info(fmt.Sprintf("IP: %s hostname: %s", c.ip, c.hostname))

	c.HealthCheckUrl = fmt.Sprintf("%s://%s:%d/health", proto, c.ip, port)
	c.StatusPageUrl = fmt.Sprintf("%s://%s:%d/info", proto, c.ip, port)
	c.HomePageUrl = fmt.Sprintf("%s://%s:%d/%s", proto, c.ip, port, homePage)

	return nil
}

func (c *Client) Register(eurekaHost, serviceName string) (*eureka.Registrar, error) {

	var fargoConfig fargo.Config
	// Target Eureka server(s).
	fargoConfig.Eureka.ServiceUrls = []string{eurekaHost}
	// How often the subscriber should poll for updates.
	fargoConfig.Eureka.PollIntervalSeconds = 1

	serviceInstance := &fargo.Instance{
		HostName:         fmt.Sprintf(c.ip),
		InstanceId:       fmt.Sprintf("%s:%s:%d", c.hostname, serviceName, c.port),
		Port:             c.port,
		PortEnabled:      true,
		App:              serviceName,
		IPAddr:           c.ip,
		VipAddress:       serviceName,
		SecureVipAddress: serviceName,
		HealthCheckUrl:   c.HealthCheckUrl,
		StatusPageUrl:    c.StatusPageUrl,
		HomePageUrl:      c.HomePageUrl,
		Status:           fargo.UP,
		DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
		CountryId:        1,
		LeaseInfo:        fargo.LeaseInfo{RenewalIntervalInSecs: 30},
	}

	// Create a Fargo connection and a Eureka registrar.
	fargoConnection := fargo.NewConnFromConfig(fargoConfig)
	registrar1 := eureka.NewRegistrar(&fargoConnection, serviceInstance, log.With(c.logger, "component", "registrar1"))

	// Register one instance.
	registrar1.Register()

	return registrar1, nil
}

func GetLocalIP() (string, error) {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("can't get InterfaceAddrs %v", err)
	}
	for _, address := range address {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("can't find any IP")
}
