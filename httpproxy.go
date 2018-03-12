package main

import (
	"net"
	"fmt"
	"bytes"
	"net/url"
	"strings"
	"io"
)

func main() {
	l, e := net.Listen("tcp", ":80")
	if e != nil {
		return
	}
	for {
		client, e := l.Accept()
		if e != nil {
			continue
		}
		handleConnect(client)
	}
}

func handleConnect(client net.Conn) {
	var buf [1024]byte
	n, e := client.Read(buf[:])

	if e != nil {
		return
	}

	var method,host,address string
	fmt.Sscanf(string(buf[:bytes.IndexByte(buf[:],'\n')]),"%s%s",&method,&host)
	parseAddrUrl, e := url.Parse(host)
	if e != nil {
		return
	}

	if(parseAddrUrl.Opaque == "443"){
		address = parseAddrUrl.Host + ":443"
	} else {
		if strings.Index(parseAddrUrl.Host,":") == -1 {
			address = parseAddrUrl.Host + "80"
		} else {
			address = parseAddrUrl.Host
		}
	}

	server, e := net.Dial("tcp", address)
	if e != nil {
		return
	}
	if method == "CONNECT"{
		fmt.Fprint(client,"HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(buf[:n])
	}

	go io.Copy(server,client)
	io.Copy(client,server)
}
