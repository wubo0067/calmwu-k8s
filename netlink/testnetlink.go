/*
 * @Author: calm.wu
 * @Date: 2019-04-08 15:16:22
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-04-08 17:36:35
 */

package main

import (
	cryptoRand "crypto/rand"
	"flag"
	"log"
	"net"
	"runtime"

	"github.com/containernetworking/plugins/pkg/ns"
	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

var paramNsName = flag.String("nsname", "", "")
var paramNewNs = flag.Bool("newns", false, "")

func generateRandomPrivateMacAddr() (string, error) {
	buf := make([]byte, 6)
	_, err := cryptoRand.Read(buf)
	if err != nil {
		return "", err
	}

	// Set the local bit for local addresses
	// Addresses in this range are local mac addresses:
	// x2-xx-xx-xx-xx-xx , x6-xx-xx-xx-xx-xx , xA-xx-xx-xx-xx-xx , xE-xx-xx-xx-xx-xx
	buf[0] = (buf[0] | 2) & 0xfe

	hardAddr := net.HardwareAddr(buf)
	return hardAddr.String(), nil
}

func getCurrNsPath() string {
	curNS, err := ns.GetCurrentNS()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%s", curNS.Path())
	return curNS.Path()
}

func listLinkInNS(nsname string) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var origns netns.NsHandle

	if nsname == "curr" {
		origns, _ = netns.Get()
	} else {
		origns, _ = netns.GetFromName(nsname)
	}

	defer origns.Close()

	log.Printf("orgins is %d name:%s uniqueID:%s\n",
		origns, origns.String(), origns.UniqueId())

	//path := getCurrNsPath()

	netlinkHandle, err := netlink.NewHandleAt(origns)
	if err != nil {
		log.Panic(err.Error())
	}

	linkList, err := netlinkHandle.LinkList()
	if err != nil {
		log.Panic(err.Error())
	}

	for i := range linkList {
		link := linkList[i]
		log.Printf("%d: Type:%s linkAttrs:%#v\n", i, link.Type(), link.Attrs())

		addrs, err := netlinkHandle.AddrList(link, netlink.FAMILY_ALL)
		if err != nil {
			log.Println(err)
			continue
		}

		routes, err := netlinkHandle.RouteList(link, netlink.FAMILY_ALL)
		if err != nil {
			log.Println(err)
			continue
		}

		for j := range addrs {
			addr := &addrs[j]
			log.Printf("%d:%d addr:%s\n", i, j, addr.String())
		}

		for k := range routes {
			route := &routes[k]
			log.Printf("%d:%d route:%s\n", i, k, route.String())
		}

		log.Println("-------------------------------------\n\n")
	}
}

func createNewNs() {
	netNS, err := ns.NewNS()
	if err != nil {
		log.Printf("NewNS failed! reason:%s\n", err.Error())
		return
	}

	// NewNs path:/var/run/netns/cni-0531baf9-7589-f42e-9863-09dae359fa1e
	// 这个会按规范创建
	log.Printf("NewNs path:%s\n", netNS.Path())
}

func main() {
	flag.Parse()

	if len(*paramNsName) > 0 {
		listLinkInNS(*paramNsName)
	}

	if *paramNewNs {
		createNewNs()
	}
}
