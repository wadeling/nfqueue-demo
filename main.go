package main

import (
	"fmt"
	knsenter "github.com/kata-containers/runtime/virtcontainers/pkg/nsenter"
	"github.com/subgraph/go-nfnetlink/nfqueue"
	"log"
	"os"
	"sync"
)

func nsenter(pid int, fn func() error) error {
	netns := knsenter.Namespace{PID: pid, Type: knsenter.NSTypeNet}
	return knsenter.NsEnter([]knsenter.Namespace{netns}, fn)
}

func openNfqueue(pid int, nfqueueNum uint16) (ch <-chan *nfqueue.NFQPacket, err error) {
	err = nsenter(pid, func() error {
		nfq := nfqueue.NewNFQueue(nfqueueNum)
		log.Printf("new nf queue %v end", nfqueueNum)
		ch, err = nfq.Open()
		if err != nil {
			log.Printf("open nf queue failed.%v", err)
		}
		return err
	})
	return
}

type MockContainer struct {
	Pid        int
	NfqueueNum uint16
}

var (
	containers = map[string]MockContainer{
		"wade-test-1": {Pid: 2034, NfqueueNum: 1},
		"wade-test-2": {Pid: 2110, NfqueueNum: 2},
	}
)

func testLocal() {
	q := nfqueue.NewNFQueue(1)

	ps, err := q.Open()
	if err != nil {
		fmt.Printf("Error opening NFQueue: %v\n", err)
		os.Exit(1)
	}
	defer q.Close()

	for p := range ps {
		networkLayer := p.Packet.NetworkLayer()
		ipsrc, ipdst := networkLayer.NetworkFlow().Endpoints()

		transportLayer := p.Packet.TransportLayer()
		tcpsrc, tcpdst := transportLayer.TransportFlow().Endpoints()

		log.Printf("A new tcp connection will be established: %s:%s -> %s:%s\n",
			ipsrc, tcpsrc, ipdst, tcpdst)
		err = p.Accept()
		if err != nil {
			log.Printf("nfqueue accept err:%v", err)
		}
	}
}

func testContainer() {
	wg := sync.WaitGroup{}
	for k, v := range containers {
		wg.Add(1)
		go func(name string, pid int, quenum uint16) {
			defer wg.Done()
			log.Printf("do container %s,pid %v,que num %v", name, pid, quenum)
			ps, err := openNfqueue(pid, quenum)
			if err != nil {
				log.Printf("open nfqueue err:%v", err)
				return
			}

			// do network work
			for p := range ps {
				networkLayer := p.Packet.NetworkLayer()
				ipsrc, ipdst := networkLayer.NetworkFlow().Endpoints()

				transportLayer := p.Packet.TransportLayer()
				tcpsrc, tcpdst := transportLayer.TransportFlow().Endpoints()

				log.Printf("A new tcp connection will be established: %s:%s -> %s:%s\n",
					ipsrc, tcpsrc, ipdst, tcpdst)
				err = p.Accept()
				if err != nil {
					log.Printf("nfqueue accept err:%v", err)
				}
			}
		}(k, v.Pid, v.NfqueueNum)
	}
	wg.Wait()
}

func main() {
	log.Printf("start")

	if len(os.Args) != 2 {
		fmt.Printf("Usage: %v [local|container]", os.Args[0])
		return
	}

	if os.Args[1] == "local" {
		testLocal()
	} else if os.Args[1] == "container" {
		testContainer()
	} else {
		fmt.Printf("unsupport param.")
	}
	log.Printf("end")
}
