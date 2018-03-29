package main

import (
	"fmt"

	//"github.com/derekparker/delve/pkg/config"
	//"bytes"
	"database/sql"
	"log"
	//"os/exec"
	//"strings"

	//"github.com/BurntSushi/toml"
	//"github.com/soniah/gosnmp"
	_ "github.com/go-sql-driver/mysql"
	//"os"
	//"io/ioutil"
	desc "strings"
)
/*
 ///////////////cisco //////////////////
 iface speed  1.3.6.1.2.1.31.1.1.1.15
 iface duplex 1.3.6.1.2.1.10.7.2.1.19   3-full 2-half 1-down/auto
 ifname       1.3.6.1.2.1.31.1.1.1.1
 operstatus   1.3.6.1.2.1.2.2.1.8     1 up 2 do
/////////////// risecom  //////////////
  
iface duplex 1.3.6.1.2.1.10.7.2.1.19.
 iface spedd  1.3.6.1.2.1.31.1.1.1.15
  fname       1.3.6.1.2.1.31.1.1.1.1
 operstatus   1.3.6.1.2.1.2.2.1.8     1 up 2 do

/////////// Eltex MES    ////////////////////
iface speed   1.3.6.1.2.1.31.1.1.1.15
 iface duplex 1.3.6.1.2.1.10.7.2.1.19.
  fname       1.3.6.1.2.1.31.1.1.1.1
  port index FE - 1..48, GE 49..100 te-105..108
  operstatus   1.3.6.1.2.1.2.2.1.8     1 up 2 do


/////////// SNR  /////////////////
iface speed   1.3.6.1.2.1.31.1.1.1.15
 iface duplex 1.3.6.1.2.1.10.7.2.1.19.   2-half/down 3 full
  fname       1.3.6.1.2.1.31.1.1.1.1
  operstatus   1.3.6.1.2.1.2.2.1.8     1 up 2 do


///// DLINK DES 3526 ///////////////////
port state   1.3.6.1.4.1.171.11.64.1.2.4.4.1.6
  other(0),
 empty(1),
 link-down(2),
 half-10Mbps(3),
 full-10Mbps(4),
 half-100Mbps(5),
 full-100Mbps(6),
 half-1Gigabps(7),
 full-1Gigabps(8),
 full-10Gigabps(9)

////  DLINK sysDescr = "DES-3200-28" ///
iface speed   1.3.6.1.2.1.31.1.1.1.15
 iface duplex 1.3.6.1.2.1.10.7.2.1.19.   2-half/down 3 full
  fname       1.3.6.1.2.1.31.1.1.1.1
  operstatus   1.3.6.1.2.1.2.2.1.8     1 up 2 do


/// DLink  sysDescr="D-Link DES-3200-28"

example out iso.3.6.1.4.1.171.11.113.1.3.2.2.1.1.5.6.100 = INTEGER: 5
    6 port 5= 100full                              |


.1.3.6.1.4.1.171.11.113.1.3.2.2.1.1.5
 empty(1),
 link-down(2),
 half-10Mbps(3),
 full-10Mbps(4),
 half-100Mbps(5),
 full-100Mbps(6),
 half-1Gigabps(7),
 full-1Gigabps(8),
 full-10Gigabps(9)


///// Dlink sysDescr="D-Link DES-3028" ////

.1.3.6.1.4.1.171.11.63.6.2.2.1.1.5
 empty(1),
 link-down(2),
 half-10Mbps(3),
 full-10Mbps(4),
 half-100Mbps(5),
 full-100Mbps(6),
 half-1Gigabps(7),
 full-1Gigabps(8),
 full-10Gigabps(9)

/// S2328  Huawei
iface speed   1.3.6.1.2.1.31.1.1.1.15
 iface duplex 1.3.6.1.2.1.10.7.2.1.19.   2-half/down 3 full
  fname       1.3.6.1.2.1.31.1.1.1.1
  operstatus   1.3.6.1.2.1.2.2.1.8     1 up 2 down
 phy port count from 5, 5 is ethernet 0/0/1, 6 is ethernet 0/0/2 etc







*/
type Hosts struct {
	id        int16
	ip        string
	community string
	Descr     string
}

func main() {
	//sysDescr := []string{".1.3.6.1.2.1.1.1.0"}


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
		if h.community !="" {
			if desc.Contains(h.Descr,"Cisco") {
				fmt.Println("IP: ",h.ip,"  Device model: Cisco")

			}
			if desc.Contains(h.Descr,"DES-3028") {
				fmt.Println("IP: ",h.ip,"  Device model: Dlink DES-3028")
			}
			if desc.Contains(h.Descr,"DES-3200-10") {
				fmt.Println("IP: ",h.ip,"  Device model: DES-3200-10")
			}
			if desc.Contains(h.Descr,"DES-3200-28") {
				fmt.Println("IP: ",h.ip,"  Device model: DES-3200-10")
			}
			if desc.Contains(h.Descr,"DES-1210-28") {
				fmt.Println("IP: ",h.ip,"  Device model: DES-1210-28")
			}
			if desc.Contains(h.Descr,"DES-2108") {
				fmt.Println("IP: ",h.ip,"  Device model: DES-2108")
			}
			if desc.Contains(h.Descr,"DES-3526") {
				fmt.Println("IP: ",h.ip,"  Device model: DES-3526")
			}			
		
			if desc.Contains(h.Descr,"DGS-3120-24SC") {
				fmt.Println("IP: ",h.ip,"  Device model: DGS-3120-24SC")
				
			}
			if desc.Contains(h.Descr,"DGS-3700-12G") {
				fmt.Println("IP: ",h.ip,"  Device model: DGS-3700-12G")
				
			}
			if desc.Contains(h.Descr,"ES-2024A") {
				fmt.Println("IP: ",h.ip,"  Device model: Zyxel ES-2024A")
				
			}
			if desc.Contains(h.Descr,"ES-3124") {
				fmt.Println("IP: ",h.ip,"  Device model: Zyxel ES-3124")
				
			}
			if desc.Contains(h.Descr,"ES-3148") {
				fmt.Println("IP: ",h.ip,"  Device model: Zyxel ES-3148")
				
			}
			if desc.Contains(h.Descr,"ISCOM2110") {
				fmt.Println("IP: ",h.ip,"  Device model: ISCOM2110")
				
			}
			if desc.Contains(h.Descr,"ISCOM2128") {
				fmt.Println("IP: ",h.ip,"  Device model: ISCOM2128")
				
			}
			if desc.Contains(h.Descr,"MES-1024") {
				fmt.Println("IP: ",h.ip,"  Device model: MES-1024 v < 1.1.30")
				
			}
			if desc.Contains(h.Descr,"MES-1124") || desc.Contains(h.Descr,"MES1124") {
				fmt.Println("IP: ",h.ip,"  Device model: MES-1124")
				
			}
			if desc.Contains(h.Descr,"MES-2124")||desc.Contains(h.Descr,"MES2124") {
				fmt.Println("IP: ",h.ip,"  Device model: MES-2124")
				
			}
			if desc.Contains(h.Descr,"MES1024") {
				fmt.Println("IP: ",h.ip,"  Device model: MES-1024 version > 1.1.30")
				
			}
			if desc.Contains(h.Descr,"MES3124") {
				fmt.Println("IP: ",h.ip,"  Device model: MES 3124")
				
			}
			if desc.Contains(h.Descr,"ROS") {
				fmt.Println("IP: ",h.ip,"  Device model: Risecom ROS 28 port")
				
			}
			if desc.Contains(h.Descr,"S2328") {
				fmt.Println("IP: ",h.ip,"  Device model: S2328P-EI-AC")
				
			}
			if desc.Contains(h.Descr,"SNR-S2940") {
				fmt.Println("IP: ",h.ip,"  Device model: SNR-S2940")
				
			}
			if desc.Contains(h.Descr,"SNR-S2950-24G") {
				fmt.Println("IP: ",h.ip,"  Device model: SNR-S2950-24G")
				
			}
			if desc.Contains(h.Descr,"SNR-S2960-24G") {
				fmt.Println("IP: ",h.ip,"  Device model: SNR-S2960-24G")
				
			} 
			if h.Descr=="" {
				fmt.Println("IP: ",h.ip,"  UNKNOWN DEVICE")
			}
			
		}
	}
}

