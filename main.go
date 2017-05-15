package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/kubernetes-incubator/service-catalog/contrib/pkg/broker/server"
	"github.com/kubernetes-incubator/service-catalog/pkg"
	"github.com/prydonius/mariadb-broker/controller"
)

var options struct {
	Port int
}

func init() {
	flag.IntVar(&options.Port, "port", 8005, "use '--port' option to specify the port for broker to listen on")
	flag.Parse()
}

func main() {
	if flag.Arg(0) == "version" {
		fmt.Printf("%s/%s\n", path.Base(os.Args[0]), pkg.VERSION)
		return
	}

	server.Start(options.Port, controller.CreateController())
}
