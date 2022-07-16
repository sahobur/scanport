package main

import (
	"fmt"
	conf "scanport/internal/config"
	"strconv"
	//"strings"
	"time"

	//"github.com/derekparker/delve/pkg/config"
	//"bytes"
	"database/sql"
	"log"

	//"os/exec"
	//"strings"
	//"math/big"

	//"github.com/BurntSushi/toml"
	//"os"

	"github.com/go-sql-driver/mysql"
	//"github.com/go-sql-driver/mysql"
	g "github.com/gosnmp/gosnmp"

	//"io/ioutil"
	s "strings"
)

const snmptimeout time.Duration = 1000000

// OIDs describes wthernet port
const (
	ifOperStatus      string = "1.3.6.1.2.1.2.2.1.8"
	ifSpeed           string = "1.3.6.1.2.1.31.1.1.1.15"
	ifDuplex          string = "1.3.6.1.2.1.10.7.2.1.19"
	ifName            string = "1.3.6.1.2.1.31.1.1.1.1"
	ifStatusDlink3028 string = "1.3.6.1.4.1.171.11.63.6.2.2.1.1.5"
	ifStatusDlink3526 string = "1.3.6.1.4.1.171.11.64.1.2.4.4.1.6"
)

// Host struct
type Hosts struct {
	id        int16  // device id in DB
	ip        string // ip address
	community string // snmp commmunity string
	Descr     string // description
}

// Ethernet Interface struct
type Interfaces struct {
	InterfacesName   string // name
	InterfacesDuplex uint64 // duplex
	InterfacesSpeed  uint64 // speed
	InterfacesStatus uint64 // state
}

func processSpecDlink(ip string, community string, model string) {
	//fmt.Println("IP: ",ip,"  DLINK 3028  ")
	g.Default.Community = community
	g.Default.Target = ip
	g.Default.Timeout = snmptimeout
	g.Default.Retries = 5
	err := g.Default.Connect()
	if err != nil {
		fmt.Print("host:=", ip, " ")
		log.Println("Connect() err: ", err)
	}
	oidd := ""
	if model == "3028" {
		oidd = ifStatusDlink3028
	}
	if model == "3526" {
		oidd = ifStatusDlink3526
	}
	resultOperStatus, err2 := g.Default.BulkWalkAll(oidd)
	if err2 != nil {
		fmt.Printf("Walk Error(ifOperstate): %v\n", err2)
		log.Println(" --ip: ", ip, " community: ", community)
		return
	}
	if model == "3028" {
		for _, r := range resultOperStatus {
			aoid := s.Split(r.Name, ".")
			ifindex, err := strconv.Atoi(aoid[16])
			if err != nil {
				panic("error string conv")
			}
			portstate := g.ToBigInt(r.Value).Uint64()

			if portstate == 2 || portstate == 3 || portstate == 4 {
				switch portstate {
				//2 10h
				//3 10f
				//4 100h
				// 5 100f
				case 2:
					fmt.Println("IP: ", ip, " IF: ", ifindex, "  DUPLEX/SPEED  HALF / 10")
				case 3:
					fmt.Println("IP: ", ip, " IF: ", ifindex, "  DUPLEX/SPEED  FULL / 10")
				case 4:
					fmt.Println("IP: ", ip, " IF: ", ifindex, "  DUPLEX/SPEED  HALF / 100")
				}

			}

		}
	}
	if model == "3526" {
		for _, r := range resultOperStatus {
			aoid := s.Split(r.Name, ".")
			ifindex, err := strconv.Atoi(aoid[16])
			if err != nil {
				panic("error string conv")
			}
			portstate := g.ToBigInt(r.Value).Uint64()

			if portstate == 3 || portstate == 4 || portstate == 5 {
				switch portstate {
				//2 10h
				//3 10f
				//4 100h
				// 5 100f
				case 3:
					fmt.Println("IP: ", ip, " IF: ", ifindex, "  DUPLEX/SPEED  HALF / 10")
				case 4:
					fmt.Println("IP: ", ip, " IF: ", ifindex, "  DUPLEX/SPEED  FULL / 10")
				case 5:
					fmt.Println("IP: ", ip, " IF: ", ifindex, "  DUPLEX/SPEED  HALF / 100")
				}

			}
			//fmt.Println("Port: ",ifindex,"  State: ",r.Value)

			//fmt.Println("Name OID: ", r.Name, "  Duplex: ", duplex)

		}
	}

}
func processStandart(ip string, community string) {
	ifs := make([]*Interfaces, 0)
	g.Default.Community = community
	g.Default.Target = ip
	g.Default.Timeout = snmptimeout
	g.Default.Retries = 4
	g.Default.MaxRepetitions = 20

	err := g.Default.Connect()
	//var ifindex []int16  = 0
	if err != nil {
		fmt.Print("host:=", ip, " ")
		log.Println("Connect() err: ", err)
	}
	resultOperStatus, err2 := g.Default.BulkWalkAll(ifOperStatus)
	if err2 != nil {
		fmt.Printf("Walk Error(Operstatus): %v\n", err2)
		log.Println(" --ip: ", ip, " community: ", community)
		return
	}
	resultDuplex, err3 := g.Default.BulkWalkAll(ifDuplex)
	if err3 != nil {
		fmt.Printf("Walk Error(ifDuplex): %v\n", err3)
		log.Println(" --ip: ", ip, " community: ", community)
		return
	}
	resultSpeed, err4 := g.Default.BulkWalkAll(ifSpeed)
	if err4 != nil {
		fmt.Printf("Walk Error(ifSpeed): %v\n", err4)
		log.Println(" --ip: ", ip, " community: ", community)
		return

	}
	resultName, err5 := g.Default.BulkWalkAll(ifName)
	if err4 != nil {
		fmt.Printf("Walk Error(ifName): %v\n", err5)
		log.Println(" --ip: ", ip, " community: ", community)
		return
	}

	// get duplex
	i := 0
	arrifindex := s.Split(resultDuplex[0].Name, ".")
	startIfindex, _ := strconv.Atoi(arrifindex[12])
	for _, r := range resultDuplex {
		I := new(Interfaces)
		I.InterfacesStatus = g.ToBigInt(r.Value).Uint64()
		ifs = append(ifs, I)
		ifs[i].InterfacesDuplex = g.ToBigInt(r.Value).Uint64()
		i++

		//fmt.Println("Name OID: ", r.Name, "  Duplex: ", duplex)

	}
	endIfindex := startIfindex + len(ifs) - 1

	//fmt.Println("func standart, ip: ",ip," Start index: ", startIfindex, "  End ifindex: ", endIfindex)

	// get oper status of port
	i = 0
	for _, r := range resultOperStatus {
		aoid := s.Split(r.Name, ".")
		ifindex, err := strconv.Atoi(aoid[11])
		if err != nil {
			panic("error string conv")
		}
		if ifindex >= startIfindex && ifindex <= endIfindex {
			ifs[i].InterfacesStatus = g.ToBigInt(r.Value).Uint64()
			i++
		} else {
			continue
		}

	}
	// get if name
	i = 0
	for _, r := range resultName {
		aoid := s.Split(r.Name, ".")
		ifindex, err := strconv.Atoi(aoid[12])
		if err != nil {
			panic("error string conv")
		}
		if ifindex >= startIfindex && ifindex <= endIfindex {
			ifs[i].InterfacesName = string(r.Value.([]byte))

			//fmt.Println("I: ", i, "  Value: ", string(r.Value.([]byte)))
			i++
		} else {
			continue
		}
	}
	// get speed
	i = 0
	for _, r := range resultSpeed {
		aoid := s.Split(r.Name, ".")
		ifindex, err := strconv.Atoi(aoid[12])
		if err != nil {
			panic("error string conv")
		}
		if ifindex >= startIfindex && ifindex <= endIfindex {
			ifs[i].InterfacesSpeed = g.ToBigInt(r.Value).Uint64()
			i++
		} else {
			continue
		}
	}
	//fmt.Println(ifs)
	for _, r := range ifs {
		if r.InterfacesStatus == 1 && (r.InterfacesDuplex == 2 || r.InterfacesSpeed == 10) {
			duplex := "UNK"
			if r.InterfacesDuplex == 2 {
				duplex = "HALF"
			}
			if r.InterfacesDuplex == 3 {
				duplex = "FULL"
			}
			fmt.Println("IP: ", ip, " IF: ", r.InterfacesName, "  DUPLEX/SPEED", duplex, "/", r.InterfacesSpeed)

		}
	}

}

func main() {
	//var cfg *conf.Config
	cfg := conf.GetConfig()
	//fmt.Printf("%+v",cfg)
	dbconn := mysql.Config {
		User: 	cfg.DBuser,
		Passwd: cfg.DBpass,
		Net: "tcp",
		Addr: cfg.DBhost+":"+cfg.DBport,
		DBName: cfg.Database,
	}
	db, err := sql.Open("mysql", dbconn.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * from communities")
	if err != nil {
		log.Fatal(err)
	}
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

	for _, h := range hst {
		if h.community != "" {
			if s.Contains(h.Descr, "Cisco") ||                s.Contains(h.Descr, "S2328") ||
				s.Contains(h.Descr, "DES-3200-10") ||         s.HasPrefix(h.Descr, "DES-3200-28") ||
				s.HasPrefix(h.Descr, "D-Link DES-3200-28") || s.Contains(h.Descr, "DES-1210-28") ||
				s.Contains(h.Descr, "DGS-3120-24SC") ||       s.Contains(h.Descr, "DGS-3700-12G") ||
				s.Contains(h.Descr, "ES-2024A") ||            s.Contains(h.Descr, "ES-3124") ||
				s.Contains(h.Descr, "ES-3148") ||             s.Contains(h.Descr, "ISCOM2110") ||
				s.Contains(h.Descr, "ISCOM2128") ||           s.Contains(h.Descr, "MES-1024") ||
				s.Contains(h.Descr, "MES-1124") ||            s.Contains(h.Descr, "MES1124") ||
				s.Contains(h.Descr, "MES-2124") ||            s.Contains(h.Descr, "MES2124") ||
				s.Contains(h.Descr, "MES1024") ||             s.Contains(h.Descr, "MES3124") ||
				s.Contains(h.Descr, "ROS") ||                 s.Contains(h.Descr, "SNR-S2940") ||
				s.Contains(h.Descr, "SNR-S2950-24G") || s.Contains(h.Descr, "SNR-S2960-24G") {

				processStandart(h.ip, h.community)

			}

			if s.Contains(h.Descr, "DES-3028") ||s.Contains(h.Descr, "DES-3526")  {
				processSpecDlink(h.ip, h.community, "3028")
			}

			if h.Descr == "" {
				fmt.Println("IP: ", h.ip, "  UNKNOWN DEVICE")
			}
			
		}
	}
	fmt.Println("Done.")
}
