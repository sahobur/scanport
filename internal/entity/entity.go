package entity

import (
	"time"
)

const SNMPTimeout time.Duration = 1000000000

// OIDs describes wthernet port
const (
	IfOperStatus      string = "1.3.6.1.2.1.2.2.1.8"
	IfSpeed           string = "1.3.6.1.2.1.31.1.1.1.15"
	IfDuplex          string = "1.3.6.1.2.1.10.7.2.1.19"
	IfName            string = "1.3.6.1.2.1.31.1.1.1.1"
	IfStatusDlink3028 string = "1.3.6.1.4.1.171.11.63.6.2.2.1.1.5"
	IfStatusDlink3526 string = "1.3.6.1.4.1.171.11.64.1.2.4.4.1.6"
)

// Host struct
type Hosts struct {
	ID        int16  // device id in DB
	IP        string // ip address
	Community string // snmp commmunity string
	Descr     string // description
}

// Ethernet Interface struct
type Interfaces struct {
	InterfacesName   string // name
	InterfacesDuplex uint64 // duplex
	InterfacesSpeed  uint64 // speed
	InterfacesStatus uint64 // state
}

const DL3028 = ".1.3.6.1.4.1.171.10.63.6"
const DL3526 = ".1.3.6.1.4.1.171.10.64.1"

var SysObjOid = []string{"1.3.6.1.2.1.1.2.0"}
