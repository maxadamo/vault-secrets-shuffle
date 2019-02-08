package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"

	"github.com/akira/go-puppetdb"
	"github.com/docopt/docopt-go"
	"github.com/go-ini/ini"
	"github.com/hashicorp/vault/api"
	"github.com/sethvargo/go-password/password"
)

var firstRe = regexp.MustCompile("^(.+) fqdn ")
var secondRe = regexp.MustCompile("}")

// ReadParams Reads info from config file
func ReadParams(configfile string) map[string]string {
	cfg, err := ini.Load(configfile)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	acmeParams := make(map[string]string)
	acmeParams["puppetdb_host"] = cfg.Section("geant_acme").Key("puppetdb_host").String()
	acmeParams["puppetdb_port"] = cfg.Section("geant_acme").Key("puppetdb_port").String()
	acmeParams["vault_host"] = cfg.Section("geant_acme").Key("vault_host").String()
	acmeParams["vault_token"] = cfg.Section("geant_acme").Key("vault_token").String()

	return acmeParams
}

// queryPuppetDB queries the puppetdb for all hosts
func queryPuppetDB(puppetdbhost string, puppetdbport string) []string {
	hostSlice := make([]string, 0)
	PuppetdbURL := fmt.Sprintf("http://%v:%s/pdb/query", puppetdbhost, puppetdbport)
	client := puppetdb.NewClient(PuppetdbURL, true)
	resp, _ := client.FactPerNode("fqdn")

	for _, value := range resp {
		stringName := fmt.Sprint(value)
		stringName = firstRe.ReplaceAllString(stringName, "")
		stringName = secondRe.ReplaceAllString(stringName, "")
		fmt.Sprintln(stringName)

		hostSlice = append(hostSlice, stringName)
	}
	return hostSlice
}

func writeSecret() {

}

func main() {

	usage := `Root password changer:
  - iterates all VMs registered in the PuppetDB
  - generate random passwords for each VM and upload them to vault.

Usage:
  root-password --config CONFIG
  root-password (-h | --help)
	
Options:
  -h --help            Show this screen.
  -c, --config=CONFIG  Config file.`

	arguments, _ := docopt.Parse(usage, nil, true, "root-password 1.0", false)
	vaultParams := ReadParams(arguments["--config"].(string))
	allHosts := queryPuppetDB(vaultParams["puppetdb_host"], vaultParams["puppetdb_port"])

	var HostPass map[string]string
	HostPass = make(map[string]string)

	for _, host := range allHosts {
		min := 2
		max := 6
		rndNum := rand.Intn(max-min) + min
		pass, _ := password.Generate(10, rndNum, 0, false, false)
		HostPass[host] = pass
	}

	//pathArg := vaultParams["vault_path"]
	vaultCFG := api.DefaultConfig()
	vaultCFG.Address = fmt.Sprintf("https://%v", vaultParams["vault_host"])

	var err error
	vClient, err := api.NewClient(vaultCFG)
	if err != nil {
		log.Fatal(err)
	}

	vClient.SetToken(vaultParams["vault_token"])
	//vault := vClient.Logical()

	secret := make(map[string]interface{})
	secret["value"] = "test secret"
	for k, v := range secret {
		fmt.Printf("Hostname[%s] password[%s]\n", k, v)
	}

	/*
		for hostKey, passValue := range HostPass {
			//fmt.Printf("Hostname[%s] password[%s]\n", k, v)
			secret[hostKey] = passValue
		}
		_, err = vault.Write(pathArg, secret)
		if err != nil {
			log.Fatal(err)
		}

		s, err := vault.Read(pathArg)
		if err != nil {
			log.Fatal(err)
		}
		if s == nil {
			log.Fatal("secret was nil")
		}

		log.Printf("%#v", *s)
	*/
}
