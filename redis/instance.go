package redis

import "net"

type Instance struct {
	ID       string		`json:"id"`
	Host     string		`json:"host"`
	Port     int 		`json:"port"`
	Password string		`json:"password"`
	MaxMemoryInMB int 	`json:"maxmemory"`
	MaxClientConnections int  `json:"maxclients"`
}

func (instance Instance) Address() *net.TCPAddr {
	return &net.TCPAddr{
		IP:   net.ParseIP(instance.Host),
		Port: instance.Port,
	}
}
