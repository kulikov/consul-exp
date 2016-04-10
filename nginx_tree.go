package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/hashicorp/consul/api"
)

func main() {
	// Get a new client
	consul, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}

	re := regexp.MustCompile(":")
	output := make(map[string]map[string][]map[string]interface{})

	servs, _, err := consul.Catalog().Services(&api.QueryOptions{})
	if err != nil {
		panic(err)
	}

	for s, tags := range servs {
		isHost := false
		for _, tag := range tags {
			if strings.HasPrefix(tag, "host:") {
				isHost = true
				break
			}
		}

		if isHost {
			services, _, err := consul.Health().Service(s, "", true, &api.QueryOptions{})
			if err != nil {
				panic(err)
			}

			for _, serv := range services {
				attrs := make(map[string]string)
				for _, t := range serv.Service.Tags {
					tt := re.Split(t, 2)
					if len(tt) > 1 {
						attrs[tt[0]] = tt[1]
					}
				}

				address := serv.Node.Address
				if serv.Service.Address != "" {
					address = serv.Service.Address
				}

				location, ok := attrs["location"]
				if ok == false {
					location = "/"
				}

				if _, ok := output[attrs["host"]]; ok == false {
					output[attrs["host"]] = make(map[string][]map[string]interface{}, 0)
				}

				if _, ok := output[attrs["host"]][location]; ok == false {
					output[attrs["host"]][location] = make([]map[string]interface{}, 0)
				}

				output[attrs["host"]][location] = append(output[attrs["host"]][location], map[string]interface{}{
					"address": address,
					"port":    serv.Service.Port,
					"staging": attrs["staging"],
				})
			}
		}
	}

	// cleanup empty or only live staging routes
	for _, hs := range output {
		for _, ls := range hs {
			if len(ls) == 1 && ls[0]["staging"] != "stage" {
				delete(ls[0], "staging")
			}
		}
	}

	out, _ := json.MarshalIndent(output, "", "  ")
	fmt.Fprintln(os.Stdout, string(out))
	os.Exit(0)
}
