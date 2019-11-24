package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
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
		log.Fatal("Fail to read file: ", err)
	}

	ConfigParams := make(map[string]string)
	ConfigParams["puppetdb_host"] = cfg.Section("vault").Key("puppetdb_host").String()
	ConfigParams["puppetdb_port"] = cfg.Section("vault").Key("puppetdb_port").String()
	ConfigParams["vault_host"] = cfg.Section("vault").Key("vault_host").String()
	ConfigParams["vault_token"] = cfg.Section("vault").Key("vault_token").String()
	ConfigParams["vault_ssl"] = cfg.Section("vault").Key("vault_ssl").String()
	ConfigParams["vault_port"] = cfg.Section("vault").Key("vault_port").String()
	ConfigParams["vault_path"] = cfg.Section("vault").Key("vault_path").String()
	ConfigParams["vault_keyname"] = cfg.Section("vault").Key("vault_keyname").String()
	ConfigParams["min_digits"] = cfg.Section("vault").Key("min_digits").String()
	ConfigParams["max_digits"] = cfg.Section("vault").Key("max_digits").String()
	ConfigParams["min_symbols"] = cfg.Section("vault").Key("min_symbols").String()
	ConfigParams["max_symbols"] = cfg.Section("vault").Key("max_symbols").String()
	ConfigParams["pass_lenght"] = cfg.Section("vault").Key("pass_lenght").String()

	return ConfigParams
}

// queryPuppetDB queries the puppetdb for all hosts
func queryPuppetDB(puppetdbhost string, puppetdbport int) []string {
	hostSlice := make([]string, 0)
	client := puppetdb.NewClient(puppetdbhost, puppetdbport, true)
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

// writeSecrets upload secrets to Vault
func writeSecrets(pwlenght string, maxdigit string, mindigit string, maxsymbol string, minsymbol string, vaulturl string, vaulttoken string, allhosts []string, patharg string, vaultkeyname string, debuginfo bool, keystore string) {
	vaultCFG := api.DefaultConfig()
	vaultCFG.Address = fmt.Sprintf(vaulturl)
	var err error
	vClient, err := api.NewClient(vaultCFG)
	if err != nil {
		log.Fatal(err)
	}

	vClient.SetToken(vaulttoken)
	vault := vClient.Logical()

	for _, host := range allhosts {
		hostUnquoted := strings.Replace(host, "\"", "", -1)
		intpwLenght, err := strconv.Atoi(pwlenght)
		intmaxDigits, err := strconv.Atoi(maxdigit)
		intminDigits, err := strconv.Atoi(mindigit)
		intmaxSymbols, err := strconv.Atoi(maxsymbol)
		intminSymbols, err := strconv.Atoi(minsymbol)

		rand.Seed(time.Now().UnixNano())
		rndDig := intminDigits + rand.Intn(intmaxDigits-intminDigits+1)
		rndSym := intminSymbols + rand.Intn(intmaxSymbols-intminSymbols+1)

		pass, _ := password.Generate(intpwLenght, rndDig, rndSym, false, false)

		secret := make(map[string]interface{})
		HostpathArg := fmt.Sprintf("/%v/data/%v/%v", patharg, hostUnquoted, vaultkeyname)
		if keystore == "1" {
			HostpathArg = fmt.Sprintf("/%v/%v/%v", patharg, hostUnquoted, vaultkeyname)
			secret["value"] = pass
		} else {
			secret["data"] = map[string]interface{}{
				"value": pass,
			}
		}
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
		if debuginfo == true {
			log.Printf("password %v for %v stored as vault:%v", pass, hostUnquoted, HostpathArg)
		} else {
			log.Printf("changed password for %v", hostUnquoted)
		}
	}
}

func main() {

	usage := `Vault Secrets Shuffler:
  - iterates all VMs registered in PuppetDB
  - generate generate random secrets different for each host
  - upload the secrets to vault.

Usage:
  vault-secrets-shuffle --config CONFIG [--kv=kv] [--debug]
  vault-secrets-shuffle (-h | --help)

Options:
  -h --help            Show this screen.
  -c, --config=CONFIG  Config file.
  -k, --kv=kv          Keystore Version. [default: 2]
  -d, --debug          Print password and full key path (OPTIONAL)`

	arguments, _ := docopt.Parse(usage, nil, true, "vault-secrets-shuffle 1.1", false)
	debugInformation := false
	if arguments["--debug"] == true {
		debugInformation = true
	}
	kv := arguments["--kv"].(string)
	if kv != "1" && kv != "2" {
		log.Fatal("Error: KeyStore version can only be 1 or 2")
	}
	vaultParams := ReadParams(arguments["--config"].(string))
	puppetDBport, err := strconv.Atoi(vaultParams["puppetdb_port"])
	if err != nil {
		log.Fatalf("Error: puppetdb_port must be an integer: %s", err)
	}
	allHosts := queryPuppetDB(vaultParams["puppetdb_host"], puppetDBport)
	pwLenght := vaultParams["pass_lenght"]
	maxDigits := vaultParams["max_digits"]
	minDigits := vaultParams["min_digits"]
	minSymbols := vaultParams["min_symbols"]
	maxSymbols := vaultParams["max_symbols"]
	vaultKEYName := vaultParams["vault_keyname"]
	vaultToken := vaultParams["vault_token"]
	vaultHTTPProto := fmt.Sprintf("http")
	if vaultParams["vault_ssl"] == "true" {
		vaultHTTPProto = fmt.Sprintf("https")
	}
	vaultURL := fmt.Sprintf("%v://%v:%v", vaultHTTPProto, vaultParams["vault_host"], vaultParams["vault_port"])
	pathArg := vaultParams["vault_path"]

	writeSecrets(
		pwLenght, maxDigits, minDigits, maxSymbols,
		minSymbols, vaultURL, vaultToken, allHosts,
		pathArg, vaultKEYName, debugInformation, kv)

}
