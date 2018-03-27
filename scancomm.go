package main

import (
	"github.com/soniah/gosnmp"
	//"github.com/derekparker/delve/pkg/config"
	//"bytes"
	"database/sql"
	"fmt"
	"log"
	//"os/exec"
	"strings"

	//"github.com/BurntSushi/toml"
	//"github.com/soniah/gosnmp"
	_ "github.com/go-sql-driver/mysql"
	//"os"
	"io/ioutil"
	//"strings"
)

func main() {
	sysDescr := []string{".1.3.6.1.2.1.1.1.0"}
	comms := []string{"public", "Etthkpi12", "eltex12", "selrktjgnsdkl"}
	bs, err := ioutil.ReadFile("hosts")
	if err != nil {
		log.Fatalln(err)
	}
	str := string(bs)
	hosts := strings.Split(str, "\n")
	// scan to db
	var hostsSlice []string = hosts[0:]
	fmt.Println(hostsSlice)
	db, err := sql.Open("mysql", "gonet:gonetpas@tcp(172.16.25.96:3306)/network")
	defer db.Close()
	for _, comm := range comms {
		for n, h := range hostsSlice {
			if h != "0" && h != "" {
				gosnmp.Default.Target = h
				gosnmp.Default.Community = comm
				gosnmp.Default.Timeout = 90000000 
				err := gosnmp.Default.Connect()
				if err != nil {
					fmt.Print("host:=",h," ")
					log.Println("Connect() err: %v", err)

				}

				result, err2 := gosnmp.Default.Get(sysDescr)
				if err2 != nil {
					fmt.Println("Get() err: host: ", h," Comm: ",comm, " Error: ", err2)
					continue
				}
				//hostsSlice = append(hostsSlice[:n], hostsSlice[n+1:]...)
				hostsSlice[n] = "0"
				fmt.Println("IP: ", h, "  community ", comm)
				for i, variable := range result.Variables {
					fmt.Printf("%d: oid: %s ", i, variable.Name)

					// the Value of each variable returned by Get() implements
					// interface{}. You could do a type switch...
					switch variable.Type {
					case gosnmp.OctetString:
						fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
					default:
						// ... or often you're just interested in numeric values.
						// ToBigInt() will return the Value as a BigInt, for plugging
						// into your calculations.
						fmt.Printf("number: %d\n", gosnmp.ToBigInt(variable.Value))
					}
				}
			}
		}
	}
	println ("Host that no snmp community:")
	var validhost int = 0
	var nosnmphost int =0
	for _, h:=range hostsSlice {
		if h !="0" {
			nosnmphost++
			println ("ip: ",h) 
		} else {
			validhost++
		 }
	}
	fmt.Println("Total hosts  : ",len(hostsSlice))
	fmt.Println("Guessed SNMP : ",validhost)
	fmt.Println("Wrong SNMP   : ",nosnmphost)
}
