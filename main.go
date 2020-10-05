package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/libvirt/libvirt-go"
)

func main() {
	id := flag.Int("id", 0, "the id the the vm you want to switch in list.json")
	flag.Parse()

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatalln("Unable to connect to libvirt")
	}
	defer conn.Close()

	file, err := ioutil.ReadFile("list.json")
	if err != nil {
		log.Fatalln("Unable to read vm list file")
	}

	var config struct {
		List []string `json:"list"`
	}

	json.Unmarshal(file, &config)

	for _, v := range config.List {
		dom, err := conn.LookupDomainByName(v)
		if err != nil {
			log.Printf("The VM %s does not exist\n", v)
			continue
		}

		defer dom.Free()

		err = dom.Shutdown()
		if err != nil {
			log.Printf("The VM %s was already powered off\n", v)
			continue
		}
		log.Printf("The VM %s has been shut down\n", v)
	}

	dom, err := conn.LookupDomainByName(config.List[*id])
	if err != nil {
		log.Printf("The VM %s could not be started because it did not exist\n", config.List[*id])
		return
	}
	defer dom.Free()

	err = dom.Create()
	if err != nil {
		log.Printf("The VM %s could not be started\n", config.List[*id])
	}
}
