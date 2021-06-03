package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
)

type Config struct {
	RemoteHost string `json:"remote_host"`
	RemotePort int `json:"remote_port"`
	LocalPort int `json:"local_port"`
}

func parseConfig() map[string]Config {
	jsonFile, err := os.Open("port.json")

	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]Config
	json.Unmarshal([]byte(byteValue), &result)

	return result
}



func forward(remote_host string, remote_port int, local_port int) 	{
	host := "0.0.0.0"
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, local_port))
	if err != nil {
		log.Fatalln(err, err.Error())
		os.Exit(0)
	}

	for {
		s_conn, err := l.Accept()
		if err != nil {
			continue
		}
		remote := fmt.Sprintf("%s:%d", remote_host, remote_port)
		d_tcpAddr, _ := net.ResolveTCPAddr("tcp4", remote)
		d_conn, err := net.DialTCP("tcp", nil, d_tcpAddr)
		if err != nil {
			fmt.Println(err)
			s_conn.Write([]byte("can't connect " + remote))
			s_conn.Close()
			continue
		}

		go func() {
			count, _ := io.Copy(s_conn, d_conn)
			log.Println(fmt.Sprintf("read remote %s %d bytes", remote, count))
		}()

		go func() {
			count, _ := io.Copy(d_conn, s_conn)
			log.Println(fmt.Sprintf("read local :%d %d bytes", local_port, count))
		}()
	}
}


func main() {
	ch := make(chan bool, 0)

	config := parseConfig()
	for k, v := range config {
		var msg = fmt.Sprintf("node: %s remote: %s:%d local: :%d", k, v.RemoteHost, v.RemotePort, v.LocalPort)
		log.Println(msg)
		go forward(v.RemoteHost, v.RemotePort, v.LocalPort)
	}
	<-ch
}
