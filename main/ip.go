package main

import (
    "net"
    "fmt"
    "os"
    "flag"
    "strings"
)

var ipv6 = flag.Bool("6", false, "use ipv6")
var ignored = flag.String("e", "", "ignored interfaces")

func main() {
    flag.Parse()
    if !flag.Parsed() {
        flag.Usage()
    }
    getFace := flag.Arg(0)

    ifaces, err := net.Interfaces()
    if err != nil {
        os.Exit(1)
    }

    excludes := strings.Split(*ignored, ",")
L:
    for _, face := range ifaces {
        show := func(iface net.Interface) {
            if addresss, err := iface.Addrs(); err != nil {

            } else {
                for _, address := range addresss {
                    if ipnet, ok := address.(*net.IPNet); ok &&
                        !ipnet.IP.IsLoopback() &&
                        !ipnet.IP.IsInterfaceLocalMulticast() &&
                        !ipnet.IP.IsMulticast() &&
                        !ipnet.IP.IsUnspecified() {
                        if ipnet.IP.To4() != nil {
                            fmt.Println(ipnet.IP.String())
                        } else if *ipv6 {
                            fmt.Println(ipnet.IP.String())
                        }
                    }
                }
            }
        }
        if getFace != "" {
            if getFace == face.Name {
                show(face)
            }
        } else {
            for _, v := range excludes {
                if ( v == face.Name ) {
                    continue L
                }
            }
            show(face)
        }
    }
}
