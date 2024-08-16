package gui

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/kevinburke/ssh_config"
)

type host_info struct {
	Name string
	HostName string
	User string
}

var hosts = []host_info{}

func getHosts() ([]host_info) {

	f, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))

	if err != nil {
		log.Fatal("Failed to ssh_config: ", err)
	}

	cfg, err := ssh_config.Decode(f)

	if err != nil {
		log.Fatal("Failed to decode ssh_config: ", err)
	}

	for _, host := range cfg.Hosts {
		if len(host.Patterns) > 0 {
			// used to show in the gui and to connect via 'ssh pattern'
			name := host.Patterns[0].String()
			if name == "*" {
				continue
			}

			// this info is only used to show in the gui
			hostname, _ := cfg.Get(name, "HostName")
			user, _ := cfg.Get(name, "User")
			hosts = append(hosts, host_info{Name: name, HostName: hostname, User: user})
		}
	}

	sort.Slice(hosts, func(i, j int) bool {
		return hosts[i].Name < hosts[j].Name
	})
	  
	return hosts
}