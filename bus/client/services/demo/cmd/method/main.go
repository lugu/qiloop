package main

import (
	"github.com/lugu/qiloop/bus/client/services"
	"github.com/lugu/qiloop/bus/session"
)

func main() {
	sess, err := session.NewSession("tcp://localhost:9559")
	if err != nil {
		panic(err)
	}

	proxies := services.Services(sess)

	directory, err := proxies.ServiceDirectory()
	if err != nil {
		panic(err)
	}

	serviceList, err := directory.Services()
	if err != nil {
		panic(err)
	}

	for _, info := range serviceList {
		println("service " + info.Name)
	}
}
