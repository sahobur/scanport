package main

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strconv"
	s "strings"
	"time"

	conf "github.com/sahobur/scanport/internal/config"
	"github.com/sahobur/scanport/internal/entity"

	"github.com/go-sql-driver/mysql"
	g "github.com/gosnmp/gosnmp"
)

func GetDlinkIfState(ip string, community string) {
	g.Default.Community = community
	g.Default.Target = ip
	g.Default.Timeout = entity.SNMPTimeout
	g.Default.Retries = 5
	err := g.Default.Connect()
	if err != nil {
		fmt.Print("host:=", ip, " ")
		log.Println("Connect() err: ", err)
	}

	oidd := ""
	model := GetDLinkModel(ip, community)

	if model == "3028" {
		oidd = entity.IfStatusDlink3028

	}
	if model == "3526" {
		oidd = entity.IfStatusDlink3526
	}

	resultOperStatus, err := g.Default.BulkWalkAll(oidd)
	if err != nil {
		fmt.Printf("Walk Error(ifOperstate): %v\n", err)
		log.Println(" --ip: ", ip, " community: ", community)
		return
	}

	if model == "3028" {
		for _, r := range resultOperStatus {
			aoid := s.Split(r.Name, ".")
			ifindex, err := strconv.Atoi(aoid[16])
			if err != nil {
				panic(err)
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
				panic(err)
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

		}
	}

}

func GetStandartIfState(ip string, community string) {
	ifs := make([]*entity.Interfaces, 0)
	g.Default.Community = community
	g.Default.Target = ip
	g.Default.Timeout = entity.SNMPTimeout
	g.Default.Retries = 4
	g.Default.MaxRepetitions = 20

	err := g.Default.Connect()
	//var ifindex []int16  = 0
	if err != nil {
		fmt.Print("host:=", ip, " ")
		log.Println("Connect() err: ", err)
	}

	resultOperStatus, err := g.Default.BulkWalkAll(entity.IfOperStatus)
	if err != nil {
		fmt.Printf("Walk Error(Operstatus): %v\n", err)
		log.Println(" --ip: ", ip, " community: ", community)
		return
	}

	resultDuplex, err := g.Default.BulkWalkAll(entity.IfDuplex)
	if err != nil {
		fmt.Printf("Walk Error(ifDuplex): %v\n", err)
		log.Println(" --ip: ", ip, " community: ", community)
		return
	}

	resultSpeed, err := g.Default.BulkWalkAll(entity.IfSpeed)
	if err != nil {
		fmt.Printf("Walk Error(ifSpeed): %v\n", err)
		log.Println(" --ip: ", ip, " community: ", community)
		return

	}
	resultName, err := g.Default.BulkWalkAll(entity.IfName)
	if err != nil {
		fmt.Printf("Walk Error(ifName): %v\n", err)
		log.Println(" --ip: ", ip, " community: ", community)
		return
	}

	// get duplex
	i := 0
	arrifindex := s.Split(resultDuplex[0].Name, ".")
	startIfindex, _ := strconv.Atoi(arrifindex[12])
	for _, r := range resultDuplex {
		I := new(entity.Interfaces)
		I.InterfacesStatus = g.ToBigInt(r.Value).Uint64()
		ifs = append(ifs, I)
		ifs[i].InterfacesDuplex = g.ToBigInt(r.Value).Uint64()
		i++
	}

	endIfindex := startIfindex + len(ifs) - 1

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
			panic(err)
		}
		if ifindex >= startIfindex && ifindex <= endIfindex {
			ifs[i].InterfacesName = string(r.Value.([]byte))

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
			panic(err)
		}
		if ifindex >= startIfindex && ifindex <= endIfindex {
			ifs[i].InterfacesSpeed = g.ToBigInt(r.Value).Uint64()
			i++
		} else {
			continue
		}
	}

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

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

// isStdOID returns truie if dev has std oid
func devType(descr string) string {

	match, err := regexp.MatchString("Cisco|S2328|DES-3200-10|DES-3200-28|D-Link DES-3200-28|DES-1210-28|DGS-3120-24SC|DGS-3700-12G|ES-2024A|ES-3124|ES-3148|ISCOM2110|ISCOM2128|MES1124|MES-1024|MES-1124|MES-2124|MES2124|MES1024|MES3124|ROS|SNR-S2940|SNR-S2950-24G", descr)
	if err != nil {
		panic(err)
	}

	if match {
		return "STD"
	}

	match, err = regexp.MatchString("DES-3028|DES-3526", descr)
	if err != nil {
		panic(err)
	}

	if match {
		return "DL3028"
	}

	return "UNK"

}

// GetDLinkModel return model of DLink dev

func GetDLinkModel(ip string, community string) string {

	var snmpinstance g.GoSNMP
	snmpinstance.Community = community
	snmpinstance.Port = 161
	snmpinstance.MaxOids = 80
	snmpinstance.Version = 0x1
	snmpinstance.Target = ip
	snmpinstance.Timeout = 2 * time.Second
	snmpinstance.Retries = 4
	err := snmpinstance.Connect()
	checkErr(err)
	defer snmpinstance.Conn.Close()
	result, err := snmpinstance.Get(entity.SysObjOid)
	checkErr(err)
	if err == nil {
		for _, variable := range result.Variables {
			//fmt.Printf("%d: oid: %s ", i, variable.Name)

			// the Value of each variable returned by Get() implements
			// interface{}. You could do a type switch...
			switch variable.Type {

			case g.OctetString:
				fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
			case g.ObjectIdentifier:
				res := variable.Value.(string)
				if res == entity.DL3028 {
					return "3028"
				}
				if res == entity.DL3526 {
					return "3526"
				}

			}
		}
	}
	return "NODLINK"
}

func DBConnect(cfg *conf.Config) (*sql.DB, error) {
	dbconn := mysql.Config{
		User:                 cfg.DBuser,
		Passwd:               cfg.DBpass,
		Net:                  "tcp",
		Addr:                 cfg.DBhost + ":" + cfg.DBport,
		DBName:               cfg.Database,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", dbconn.FormatDSN())
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetHosts(db *sql.DB) []*entity.Hosts {
	rows, err := db.Query("SELECT * from communities")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	hst := make([]*entity.Hosts, 0)
	for rows.Next() {
		h := new(entity.Hosts)
		err := rows.Scan(&h.ID, &h.IP, &h.Community, &h.Descr)
		if err != nil {
			log.Fatal(err)
		}
		hst = append(hst, h)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return hst
}

func main() {

	cfg := conf.GetConfig()
	db, err := DBConnect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	hst := GetHosts(db)
	for _, h := range hst {
		if h.Community != "" {
			dtype := devType(h.Descr)
			if dtype == "STD" {
				fmt.Println("IP: ", h.IP, "  ", h.Community)
				GetStandartIfState(h.IP, h.Community)
			}

			if dtype == "DL3028" {
				GetDlinkIfState(h.IP, h.Community)
			}

			if dtype == "UNK" {
				fmt.Println("IP: ", h.IP, "  UNKNOWN DEVICE")
			}

		}
	}

	fmt.Println("Done.")
}
