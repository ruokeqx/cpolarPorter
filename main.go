package main

import (
	"fmt"
	"time"

	"github.com/ruokeqx/cpolarPorter/alidns"
	"github.com/ruokeqx/cpolarPorter/cpolar"
	"github.com/ruokeqx/cpolarPorter/env"
)

// Environment variables names.
const (
	envNamespace = "DOMAIN_"

	EnvDomain = envNamespace + "DOMAIN"
	EnvRR     = envNamespace + "RR"
)

func main() {
	domain := env.GetOrFile(EnvDomain)
	RR := env.GetOrFile(EnvRR)

	pm := NewPortManager()

	dnsProvidor, err := alidns.NewDNSProvider()
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	cc := cpolar.NewCpolarConnector()
	for {
		err := cc.Login()
		if err != nil {
			fmt.Print(err.Error())
			time.Sleep(time.Hour)
			continue
		}
	tunnel:
		tunnels, err := cc.Tunnels()
		if err != nil {
			fmt.Print(err.Error())
			continue
		}

		if pm.Update(tunnels) {
			// map Marshal should not fail
			mb, _ := pm.Marshal()
			fmt.Println(string(mb))

			err = dnsProvidor.CleanUp(domain, RR)
			if err != nil {
				fmt.Print(err.Error())
				return
			}
			err = dnsProvidor.Present(domain, RR, string(mb))
			if err != nil {
				fmt.Print(err.Error())
				return
			}
			fmt.Println("success")
		}

		time.Sleep(time.Hour)
		goto tunnel
	}
}
