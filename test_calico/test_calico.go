/*
 * @Author: calm.wu
 * @Date: 2019-04-25 16:43:20
 * @Last Modified by: calm.wu
 * @Last Modified time: 2019-04-25 16:48:49
 */

package main

import (
	"runtime"
	"sync"
	"strconv"
	"fmt"
	"net"
	"context"

	"github.com/projectcalico/libcalico-go/lib/apiconfig"
	client "github.com/projectcalico/libcalico-go/lib/clientv3"
	"github.com/projectcalico/libcalico-go/lib/names"
	uuid "github.com/satori/go.uuid"
	cnet "github.com/projectcalico/libcalico-go/lib/net"
	"github.com/projectcalico/libcalico-go/lib/ipam"
)

type InheritanceInt struct {
	int
}

func (d InheritanceInt) String() {
	fmt.Println(strconv.Itoa(d.int))
}

func parseCIDR() {
	cidr := "10.244.53.128/26"

	netIP, netIPNet, e := net.ParseCIDR(cidr)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	fmt.Printf("%s %v %s\n", netIP.To4().String(), netIPNet, netIPNet.Mask.String())
}

func createCalicoClient() client.Interface {
	calicoConfig := "/etc/calico/calicoctl.cfg"

	clientConfig, err := apiconfig.LoadClientConfig(calicoConfig)
	if err != nil {
		fmt.Printf("LoadClientConfig failed! reason:%s\n", err.Error())
		return nil
	}

	calicoClient, err := client.New(*clientConfig)
	if err != nil {
		fmt.Printf("client.New failed! reason:%s\n", err.Error())
		return nil
	}
	return calicoClient
}

func allocCalicoIP(calicoClient client.Interface) {
	ctx := context.Background()
	v4pools := []cnet.IPNet{}	
	v6pools := []cnet.IPNet{}	
	nodename, _ := names.Hostname()
	assignArgs := ipam.AutoAssignArgs{
		Num4:      1,
		Num6:      0,
		//HandleID:  &handleID,
		Hostname:  nodename,
		IPv4Pools: v4pools,
		IPv6Pools: v6pools,
		//Attrs:     attrs,
	}

	var wg sync.WaitGroup
	af := func(id int) {
		defer wg.Done()
		assignedV4, _, err := calicoClient.IPAM().AutoAssign(ctx, assignArgs)
		if err != nil {
			fmt.Println("AutoAssign failed! reason:%s\n", err.Error())
			return
		}
		var allocIP net.IP 
		allocIP = assignedV4[0].IP
		fmt.Printf("[%d]----------autoassign ip:%s\n", id, allocIP.String())
	
		//_, err = calicoClient.IPAM().ReleaseIPs(ctx, []cnet.IP{cnet.IP{allocIP},})
		_, err = calicoClient.IPAM().ReleaseIPs(ctx, []cnet.IP{cnet.IP{allocIP},})
		if err != nil {
			fmt.Println("ReleaseIPs failed! reason:%s\n", err.Error())
			return		
		}
		fmt.Printf("[%d]++++++++++release ip:%s\n", id, allocIP)
	}
	
	parallelCount := 10
	runtime.GOMAXPROCS(parallelCount)
	wg.Add(parallelCount)
	for i := 0; i < parallelCount; i++ {
		go af(i)
	}

	wg.Wait()
}

func main() {
	di := InheritanceInt{9996}
	di.String()

	fmt.Println(uuid.NewV4().String())
	parseCIDR()
	client := createCalicoClient()
	if client == nil {
		fmt.Println("create failed!")
	}
	allocCalicoIP(client)
}
