package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"time"

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
	ConfigParams["vault_ssl"] = cfg.Section("vault_secrets").Key("vault_ssl").String()
	ConfigParams["vault_port"] = cfg.Section("vault_secrets").Key("vault_port").String()
	ConfigParams["vault_path"] = cfg.Section("vault_secrets").Key("vault_path").String()
	ConfigParams["vault_keyname"] = cfg.Section("vault_secrets").Key("vault_keyname").String()
	ConfigParams["min_digits"] = cfg.Section("vault_secrets").Key("min_digits").String()
	ConfigParams["max_digits"] = cfg.Section("vault_secrets").Key("max_digits").String()
	ConfigParams["min_symbols"] = cfg.Section("vault_secrets").Key("min_symbols").String()
	ConfigParams["max_symbols"] = cfg.Section("vault_secrets").Key("max_symbols").String()
	ConfigParams["pass_lenght"] = cfg.Section("vault_secrets").Key("pass_lenght").String()

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
	var pwLenght = vaultParams["pass_lenght"]
	var maxDigits = vaultParams["max_digits"]
	var minDigits = vaultParams["min_digits"]
	var minSymbols = vaultParams["min_symbols"]
	var maxSymbols = vaultParams["max_symbols"]
	var vaultKEYName = vaultParams["vault_keyname"]
	vaultHTTPProto := fmt.Sprintf("http")
	if vaultParams["vault_ssl"] == "true" {
		vaultHTTPProto = fmt.Sprintf("https")
	} else {
		vaultHTTPProto = fmt.Sprintf("http")
	}

	pathArg := vaultParams["vault_path"]
	vaultCFG := api.DefaultConfig()
	if vaultParams["vault_port"] != "443" {
		vaultCFG.Address = fmt.Sprintf("%v://%v:%v", vaultHTTPProto, vaultParams["vault_host"], vaultParams["vault_port"])
	} else {
		vaultCFG.Address = fmt.Sprintf("%v://%v", vaultHTTPProto, vaultParams["vault_host"])
	}
	var err error
	vClient, err := api.NewClient(vaultCFG)
	if err != nil {
		log.Fatal(err)
	}

	vClient.SetToken(vaultParams["vault_token"])
	vault := vClient.Logical()

	for _, host := range allHosts {
		intpwLenght, err := strconv.Atoi(pwLenght)
		intmaxDigits, err := strconv.Atoi(maxDigits)
		intminDigits, err := strconv.Atoi(minDigits)
		intmaxSymbols, err := strconv.Atoi(maxSymbols)
		intminSymbols, err := strconv.Atoi(minSymbols)

		rand.Seed(time.Now().UnixNano())
		rndDig := intminDigits + rand.Intn(intmaxDigits-intminDigits+1)
		rndSym := intminSymbols + rand.Intn(intmaxSymbols-intminSymbols+1)

		pass, _ := password.Generate(intpwLenght, rndDig, rndSym, false, false)
		//data := make(map[string]map[string]interface{})
		//data["secret"]["value"] = pass
		secret := make(map[string]interface{})
		secret["value"] = pass
		HostpathArg := fmt.Sprintf("/%v/data/%v/%v", pathArg, host, vaultKEYName)
		//_, err = vault.Write(HostpathArg, data)
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
		log.Printf("/secret/data/%v/%v/%v", pathArg, host, vaultKEYName)
		log.Printf("%#v", *s)
	}

}
