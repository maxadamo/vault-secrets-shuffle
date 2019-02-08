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

	ConfigParams := make(map[string]string)
	ConfigParams["puppetdb_host"] = cfg.Section("vault_secrets").Key("puppetdb_host").String()
	ConfigParams["puppetdb_port"] = cfg.Section("vault_secrets").Key("puppetdb_port").String()
	ConfigParams["vault_host"] = cfg.Section("vault_secrets").Key("vault_host").String()
	ConfigParams["vault_token"] = cfg.Section("vault_secrets").Key("vault_token").String()

	return ConfigParams
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
	var vaultHTTPProto = "https"
	var maxDigits = vaultParams["max_digits"]
	var minDigits = vaultParams["min_digits"]
	var minSymbols = vaultParams["min_symbols"]
	var maxSymbols = vaultParams["max_symbols"]
	if vaultParams["vault_ssl"] != "true" {
		vaultHTTPProto = fmt.Sprintf("http")
	}

	// var HostPass map[string]string
	// HostPass = make(map[string]string)

	pathArg := vaultParams["vault_path"]
	vaultCFG := api.DefaultConfig()
	vaultCFG.Address = fmt.Sprintf("%v://%v:%v", vaultHTTPProto, vaultParams["vault_host"], vaultParams["vault_port"])

	var err error
	vClient, err := api.NewClient(vaultCFG)
	if err != nil {
		log.Fatal(err)
	}

	vClient.SetToken(vaultParams["vault_token"])
	vault := vClient.Logical()

	for _, host := range allHosts {
		min := 2
		max := 6
		rndDig := rand.Intn(maxDigits-minDigits) + minDigits
		rndSym := rand.Intn(maxSymbols-minSymbols) + minSymbols
		pass, _ := password.Generate(10, rndDig, 0, false, false)
		//HostPass[host] = pass

		secret := make(map[string]interface{})
		secret["value"] = pass
		HostpathArg := fmt.Sprintf("%v/%v", pathArg, host)

		_, err = vault.Write(HostpathArg, secret)
		if err != nil {
			log.Fatal(err)
		}

		s, err := vault.Read(HostpathArg)
		if err != nil {
			log.Fatal(err)
		}
		if s == nil {
			log.Fatal("secret was nil")
		}

		log.Printf("%#v", *s)
	}

}
