package main

import (
	"fmt"
	"strconv"
	"strings"
	"database/sql"
	"log"
	"os"
	_ "github.com/go-sql-driver/mysql"
	g "github.com/gosnmp/gosnmp"

	//"io/ioutil"
	desc "strings"
)

var ifOperStatus string = "1.3.6.1.2.1.2.2.1.8"
var ifSpeed string = "1.3.6.1.2.1.31.1.1.1.15"
var ifDuplex string = "1.3.6.1.2.1.10.7.2.1.19"
var ifName string = "1.3.6.1.2.1.31.1.1.1.1"
var ifStatusDlink3028 string = "1.3.6.1.4.1.171.11.63.6.2.2.1.1.5"
var ifStatusDlink3526 string = "1.3.6.1.4.1.171.11.64.1.2.4.4.1.6"

type Hosts struct {
	id        int16
	ip        string
	community string
	Descr     string
}
type Interfaces struct {
	InterfacesName   string
	InterfacesDuplex uint64
	InterfacesSpeed  uint64
	InterfacesStatus uint64
}

func processCisco(ip string, community string) {
	//fmt.Println("IP: ", h.ip, "  Device model: Cisco")
	ifs := make([]*Interfaces, 0)
	g.Default.Community = community
	g.Default.Target = ip
	err := g.Default.Connect()
	//var ifindex []int16  = 0
	if err != nil {
		fmt.Print("host:=", ip, " ")
		log.Println("Connect() err: ", err)
	}
	resultOperStatus, err2 := g.Default.BulkWalkAll(ifOperStatus)
	if err2 != nil {
		fmt.Printf("Walk Error: %v\n", err)
		return
	}
	resultDuplex, err3 := g.Default.BulkWalkAll(ifDuplex)
	if err3 != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}
	resultSpeed, err4 := g.Default.BulkWalkAll(ifSpeed)
	if err4 != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}
	resultName, err4 := g.Default.BulkWalkAll(ifName)
	if err4 != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}

	// get duplex

	for i, r := range resultDuplex {
		I := new(Interfaces)
		I.InterfacesStatus = g.ToBigInt(r.Value).Uint64()
		ifs = append(ifs, I)
		ifs[i].InterfacesDuplex = g.ToBigInt(r.Value).Uint64()
		i++
	}

	// get oper status of port
	i := 0
	for _, r := range resultOperStatus {
		aoid := strings.Split(r.Name, ".")
		ifindex, err := strconv.Atoi(aoid[11])
		if err != nil {
			panic("error string conv")
		}
		if ifindex > 9999 && ifindex < 10400 {
			//fmt.Println(ifindex)
			ifs[i].InterfacesStatus = g.ToBigInt(r.Value).Uint64()
			i++
		} else {
			continue
		}

	}
	// get if name
	i = 0
	for _, r := range resultName {
		aoid := strings.Split(r.Name, ".")
		ifindex, err := strconv.Atoi(aoid[12])
		if err != nil {
			panic("error string conv")
		}
		if ifindex > 9999 && ifindex < 10400 {
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
		aoid := strings.Split(r.Name, ".")
		ifindex, err := strconv.Atoi(aoid[12])
		if err != nil {
			panic("error string conv")
		}
		if ifindex > 9999 && ifindex < 10400 {
			//fmt.Println("speed:")
			ifs[i].InterfacesSpeed = g.ToBigInt(r.Value).Uint64()
			//fmt.Println("Name OID: ", r.Name, "  Value: ", r.Value)
			i++
		} else {
			continue
		}
	}
	//fmt.Println(ifs)
	for _, r := range ifs {
		if r.InterfacesStatus == 1 && (r.InterfacesDuplex == 2 || r.InterfacesSpeed == 10) {
			fmt.Println("IP: ", ip, " IF: ", r.InterfacesName, " STATUS: ", r.InterfacesStatus, "  DUPLEX/SPEED", r.InterfacesDuplex, "/", r.InterfacesSpeed)
		}
	}
	//panic("STOP")
}
func processHuaweiS23(ip string, community string) {
	ifs := make([]*Interfaces, 0)
	g.Default.Community = community
	g.Default.Target = ip
	g.Default.Timeout = 9000000000
	g.Default.Retries = 5
	err := g.Default.Connect()
	//var ifindex []int16  = 0
	if err != nil {
		fmt.Print("host:=", ip, " ")
		log.Println("Connect() err: ", err)
	}
	resultOperStatus, err2 := g.Default.BulkWalkAll(ifOperStatus)
	if err2 != nil {
		fmt.Printf("Walk Error(OperStatus): %v\n", err2)
		log.Println(" --ip: ", ip, " community: ", community)
		return
	}
	resultDuplex, err3 := g.Default.BulkWalkAll(ifDuplex)
	if err3 != nil {
		fmt.Printf("Walk Error(ifDuplex): %v\n", err3)
		log.Println(" --ip: ", ip, " community: ", community)
	}
	resultSpeed, err4 := g.Default.BulkWalkAll(ifSpeed)
	if err4 != nil {
		fmt.Printf("Walk Error(ifSpeed): %v\n", err4)
		log.Println(" --ip: ", ip, " community: ", community)
	}
	resultName, err5 := g.Default.BulkWalkAll(ifName)
	if err4 != nil {
		fmt.Printf("Walk Error(ifName): %v\n", err5)
		log.Println(" --ip: ", ip, " community: ", community)
	}

	// get duplex
	i := 0
	arrifindex := strings.Split(resultDuplex[0].Name, ".")

	startIfindex, _ := strconv.Atoi(arrifindex[12])
	for _, r := range resultDuplex {
		I := new(Interfaces)
		I.InterfacesStatus = g.ToBigInt(r.Value).Uint64()
		ifs = append(ifs, I)
		ifs[i].InterfacesDuplex = g.ToBigInt(r.Value).Uint64()
		i++

		//fmt.Println("Name OID: ", r.Name, "  Duplex: ", duplex)

	}
	endIfindex := len(ifs)
	//fmt.Println("Start index: ", startIfindex, "  End ifindex: ", endIfindex)
	//ints:=i
	// get oper status of port
	i = 0
	for _, r := range resultOperStatus {
		aoid := strings.Split(r.Name, ".")
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
		aoid := strings.Split(r.Name, ".")
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
		aoid := strings.Split(r.Name, ".")
		ifindex, err := strconv.Atoi(aoid[12])
		if err != nil {
			panic("error string conv")
		}
		if ifindex >= startIfindex && ifindex <= endIfindex {
			//fmt.Println("speed:")
			ifs[i].InterfacesSpeed = g.ToBigInt(r.Value).Uint64()
			//fmt.Println("Name OID: ", r.Name, "  Value: ", r.Value)
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
	//panic("STOP")
}
func processSpecDlink(ip string, community string, model string) {
	//fmt.Println("IP: ",ip,"  DLINK 3028  ")
	g.Default.Community = community
	g.Default.Target = ip
	g.Default.Timeout = 9000000000
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
			aoid := strings.Split(r.Name, ".")
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
			//fmt.Println("Port: ",ifindex,"  State: ",r.Value)

			//fmt.Println("Name OID: ", r.Name, "  Duplex: ", duplex)

		}
	}
	if model == "3526" {
		for _, r := range resultOperStatus {
			aoid := strings.Split(r.Name, ".")
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
	g.Default.Timeout = 10000000000
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
	arrifindex := strings.Split(resultDuplex[0].Name, ".")
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
		aoid := strings.Split(r.Name, ".")
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
		aoid := strings.Split(r.Name, ".")
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
		aoid := strings.Split(r.Name, ".")
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
	//panic("STOP")
}

func main() {
	//sysDescr := []string{".1.3.6.1.2.1.1.1.0"}
	//i:=1
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

	for _, h := range hst {
		//fmt.Println("Hosts processed:",i,"Host ",h.ip," community ",h.community," Sysdesc ",h.Descr)
		if h.community != "" {
			if desc.Contains(h.Descr, "Cisco") {
				//processCisco(h.ip, h.community)
				processStandart(h.ip, h.community)
			}

			if desc.Contains(h.Descr, "S2328") {

				//processHuaweiS23(h.ip, h.community)
				processStandart(h.ip, h.community)
			}

			if desc.Contains(h.Descr, "DES-3028") {
				processSpecDlink(h.ip, h.community, "3028")
			}

			if desc.Contains(h.Descr, "DES-3200-10") {
				//fmt.Println("IP: ", h.ip, "  Device model: DES-3200-10")
				processStandart(h.ip, h.community)
			}
			if desc.HasPrefix(h.Descr, "DES-3200-28") {
				//fmt.Println("IP:", h.ip, "DES-3200-28")
				processStandart(h.ip, h.community)
			}

			if desc.HasPrefix(h.Descr, "D-Link DES-3200-28") {
				//fmt.Println("IP:", h.ip, "D-Link DES-3200-28")
				processStandart(h.ip, h.community)
			}

			if desc.Contains(h.Descr, "DES-1210-28") {
				processStandart(h.ip, h.community)

			}

			if desc.Contains(h.Descr, "DES-2108") {
				fmt.Println("IP: ", h.ip, " UNABLE GET SNMP PORT STATE FOR DES-2108")
				//processDES3028(h.ip, h.community)
			}

			if desc.Contains(h.Descr, "DES-3526") {
				//fmt.Println("IP: ", h.ip, "  Device model: DES-3526")
				processSpecDlink(h.ip, h.community, "3526")
			}

			if desc.Contains(h.Descr, "DGS-3120-24SC") {
				//	fmt.Println("IP: ", h.ip, "  Device model: DGS-3120-24SC")
				processStandart(h.ip, h.community)
			}
			if desc.Contains(h.Descr, "DGS-3700-12G") {
				//	fmt.Println("IP: ", h.ip, "  Device model: DGS-3700-12G")
				processStandart(h.ip, h.community)
			}

			if desc.Contains(h.Descr, "ES-2024A") {
				//	fmt.Println("IP: ", h.ip, "  Device model: Zyxel ES-2024A")
				processStandart(h.ip, h.community)

			}

			if desc.Contains(h.Descr, "ES-3124") {
				//fmt.Println("IP: ", h.ip, "  Device model: Zyxel ES-3124")
				processStandart(h.ip, h.community)
			}

			if desc.Contains(h.Descr, "ES-3148") {
				//	fmt.Println("IP: ", h.ip, "  Device model: Zyxel ES-3148")
				processStandart(h.ip, h.community)

			}

			if desc.Contains(h.Descr, "ISCOM2110") {
				//	fmt.Println("IP: ", h.ip, "  Device model: ISCOM2110")
				processStandart(h.ip, h.community)

			}
			if desc.Contains(h.Descr, "ISCOM2128") {
				//	fmt.Println("IP: ", h.ip, "  Device model: ISCOM2128")
				processStandart(h.ip, h.community)

			}
			if desc.Contains(h.Descr, "MES-1024") {
				//	fmt.Println("IP: ", h.ip, "  Device model: MES-1024 v < 1.1.30")
				processStandart(h.ip, h.community)

			}
			if desc.Contains(h.Descr, "MES-1124") || desc.Contains(h.Descr, "MES1124") {
				//	fmt.Println("IP: ", h.ip, "  Device model: MES-1124")
				processStandart(h.ip, h.community)

			}
			if desc.Contains(h.Descr, "MES-2124") || desc.Contains(h.Descr, "MES2124") {
				//	fmt.Println("IP: ", h.ip, "  Device model: MES-2124")
				processStandart(h.ip, h.community)

			}
			if desc.Contains(h.Descr, "MES1024") {
				//	fmt.Println("IP: ", h.ip, "  Device model: MES-1024 version > 1.1.30")
				processStandart(h.ip, h.community)

			}
			if desc.Contains(h.Descr, "MES3124") {
				//	fmt.Println("IP: ", h.ip, "  Device model: MES 3124")
				processStandart(h.ip, h.community)

			}
			if desc.Contains(h.Descr, "ROS") {
				//	fmt.Println("IP: ", h.ip, "  Device model: Risecom ROS 28 port")
				processStandart(h.ip, h.community)

			}

			if desc.Contains(h.Descr, "SNR-S2940") {
				//	fmt.Println("IP: ", h.ip, "  Device model: SNR-S2940")
				processStandart(h.ip, h.community)

			}
			if desc.Contains(h.Descr, "SNR-S2950-24G") {
				//	fmt.Println("IP: ", h.ip, "  Device model: SNR-S2950-24G")
				processStandart(h.ip, h.community)

			}
			if desc.Contains(h.Descr, "SNR-S2960-24G") {
				//	fmt.Println("IP: ", h.ip, "  Device model: SNR-S2960-24G")
				processStandart(h.ip, h.community)

			}
			if h.Descr == "" {
				fmt.Println("IP: ", h.ip, "  UNKNOWN DEVICE")
			}
			//i++
		}
	}
	fmt.Println("Done.")
}
