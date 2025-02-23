package remote

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Host struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func HostsFromYaml(filename string) ([]Host, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	hosts := make([]Host, 0)
	err = yaml.Unmarshal(content, &hosts)
	if err != nil {
		return nil, err
	}

	for i, _ := range hosts {
		if hosts[i].Port == 0 {
			hosts[i].Port = 22
		}
		if hosts[i].Username == "" {
			hosts[i].Username = "root"
		}
	}
	return hosts, nil
}
