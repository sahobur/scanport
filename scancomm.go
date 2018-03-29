package main

import (
	"github.com/soniah/gosnmp"
	//"github.com/derekparker/delve/pkg/config"
	//"bytes"
	"database/sql"
	"fmt"
	"log"
	//"os/exec"
	//"strings"

	//"github.com/BurntSushi/toml"
	//"github.com/soniah/gosnmp"
	_ "github.com/go-sql-driver/mysql"
	//"os"
	//"io/ioutil"
	//"strings"
)

type Hosts struct {
	id        int16
	ip        string
	community string
	Descr     string
}

func main() {
	sysDescr := []string{".1.3.6.1.2.1.1.1.0"}
	comms := []string{"eltex12","public", "Public", "Etthkpi12",  "selrktjgnsdkl"}
	// if process data from db
	//bs, err := ioutil.ReadFile("hosts")
	//i/f err != nil {
	//	log.Fatalln(err)
//	}
	//str := string(bs)
	//hosts := strings.Split(str, "\n")

	//var hostsSlice []string = hosts[0:]
	//fmt.Println(hostsSlice)

	db, err := sql.Open("mysql", "gonet:gonetpas@tcp(172.16.25.96:3306)/network")
	defer db.Close()
	rows, err := db.Query("SELECT * from communities")
	defer rows.Close()
	hst := make([]*Hosts, 0)
	for rows.Next() {
		h := new(Hosts)
		err := rows.Scan(&h.id, &h.ip, &h.community, &h.Descr)
		if err != nil {
			log.Fatal(err)
		}
		hst = append(hst, h)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
/*
	for n, gg := range hst {
		hostsSlice[n] = gg.ip
		fmt.Printf("%d, %s, %s, %s\n", gg.id, gg.ip, gg.community, gg.Descr)
		fmt.Println("---------------------------------------------------------------------------")
	}
	//fmt.Println(hostsSlice)
	//panic("dddddddddddddddddddd")
	*/

	for _, comm := range comms {
		fmt.Println("#####################  ",comm,"  #########################")
		for _, h := range hst {
			
			fmt.Println("h.ip = "+h.ip+" h.commm = ", h.community, " Descr = "+h.Descr)
			if h.community == "" || (h.community != "" && h.Descr == "") {
				gosnmp.Default.Target = h.ip
				if h.community == "" {
					gosnmp.Default.Community = comm
				} else {
					gosnmp.Default.Community = h.community
				}

				gosnmp.Default.Timeout = 390000000
				gosnmp.Default.Retries = 4
				err := gosnmp.Default.Connect()
				if err != nil {
					fmt.Print("host:=", h.ip, " ")
					log.Println("Connect() err: ", err)

				}

				result, err2 := gosnmp.Default.Get(sysDescr)
				if err2 != nil {
					fmt.Println("Get() err: host: ", h.ip, " Comm: ", gosnmp.Default.Community, " Error: ", err2)
						continue
				}
				//hostsSlice = append(hostsSlice[:n], hostsSlice[n+1:]...)
				//hostsSlice[n] = "0"
				//fmt.Println("IP: ", h, "  community ", comm)
				for i, variable := range result.Variables {
					fmt.Printf("%d: oid: %s ", i, variable.Name)

					// the Value of each variable returned by Get() implements
					// interface{}. You could do a type switch...
					switch variable.Type {
					case gosnmp.OctetString:
						fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
						if h.community == "" {
							_, err := db.Exec("UPDATE communities set community=\"" + comm + "\"where ip = \"" + h.ip + "\"")
							h.community = comm
							fmt.Println("updating community ip: ",h.ip," community: ",h.community)
							//panic("wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww")
							fmt.Println("--update result: ",err)
							if err != nil {
								panic(err)
							}
						}
						if h.Descr == "" {
							_, err := db.Exec("UPDATE communities set sysdescr=\"" + string(variable.Value.([]byte)) + "\"where ip=\"" + h.ip + "\"")
							if err != nil {
								panic(err)
							}
						}
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
	println("Host that no snmp community:")
	var validhost int = 0
	var nosnmphost int = 0
	for _, h := range hst {
		if h.community == "" {
			nosnmphost++
			//println("ip: ", h.ip)
		} else {
			validhost++
		}
	}
	fmt.Println("Total hosts  : ", len(hst))
	fmt.Println("Guessed SNMP : ", validhost)
	fmt.Println("Wrong SNMP   : ", nosnmphost)
}
