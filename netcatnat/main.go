package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	nw "github.com/getlantern/nattywad"
	udt "github.com/jbenet/go-udtwrapper/udt"
)

var waddellAddr string
var serverID string
var VERBOSE bool
var TRACE bool

func init() {
	TRACE = os.Getenv("TRACE") == "true"
	flag.StringVar(&waddellAddr, "waddell", "", "waddell signaling service address (required)")
	flag.StringVar(&serverID, "id", "", "id of the peer to connect to (optional)")
}

func Connect(id string) {
	nclient := nw.Client{
		DialWaddell:       dialWaddel,
		OnSuccess:         onClientSuccess,
		OnFailure:         onClientFailure,
		KeepAliveInterval: time.Second * 10,
	}

	sp := &nw.ServerPeer{ID: id, WaddellAddr: waddellAddr}
	nclient.Configure([]*nw.ServerPeer{sp})
}

func Serve() {
	nserver := nw.Server{OnSuccess: onServerSuccess}
	nserver.Configure(waddellAddr)
}

func dialWaddel(addr string) (net.Conn, error) {
	return net.Dial("tcp", waddellAddr)
}

func onClientSuccess(info *nw.TraversalInfo) {
	trace("Traversal Succeeded: %s\n", info)
	trace("Peer Country: %s\n", info.Peer.Extras["country"])
	trace("Peer ID: %s\n", info.Peer.ID)

	laddr := info.LocalAddr
	raddr := info.RemoteAddr
	log("connecting %s to %s\n", laddr, raddr)
	go func() {
		if err := netcat(laddr, raddr, false); err != nil {
			log("failed to netcat: %s\n", err)
		}
	}()
}

func onClientFailure(info *nw.TraversalInfo) {
	trace("Traversal Failed: %s\n", info)
	trace("Peer Country: %s\n", info.Peer.Extras["country"])
	trace("Peer ID: %s\n", info.Peer.ID)
	log("failed to connect to %s\n", info.Peer.ID)
}

func onServerSuccess(laddr, raddr *net.UDPAddr) bool {
	log("connecting %s to %s\n", laddr, raddr)
	go func() {
		if err := netcat(laddr, raddr, true); err != nil {
			log("failed to netcat: %s\n", err)
		}
	}()
	return true
}

func netcat(laddr, raddr *net.UDPAddr, listen bool) error {
	laddr2 := udt.WrapUDPAddr(laddr)
	raddr2 := udt.WrapUDPAddr(raddr)

	var conn net.Conn
	var err error
	if listen {
		l, err2 := udt.ListenUDT("udt", laddr2)
		if err2 != nil {
			return err
		}
		log("listening at %s\n", laddr2)
		conn, err = l.Accept()
	} else {
		log("dialing from %s to %s\n", laddr2, raddr2)
		conn, err = udt.DialUDT("udt", laddr2, raddr2)
	}
	if err != nil {
		return err
	}

	go io.Copy(conn, os.Stdin)
	go io.Copy(os.Stdout, conn)

	log("\n>> connected! enter text below to transmit <<\n\n")

	return nil
}

func main() {
	flag.Parse()

	if waddellAddr == "" {
		usage()
	}
	log("using waddell server: %s\n", waddellAddr)

	if len(serverID) > 1 {
		log("attempting to connect to: %s\n", serverID)
		Connect(serverID)
	} else {
		Serve()
	}

	// wait until we exit.
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT,
		syscall.SIGTERM, syscall.SIGQUIT)
	<-sigc
}

func usage() {
	t := `usage: %s --waddell <address> [--id <id>]

Connects via waddell signaling server at <address>.
No --id puts %s in listening mode, it will output its <id>.
Specify --id <id> to connect to the corresponsing process.
`
	procname := os.Args[0]
	fmt.Fprintf(os.Stderr, t, procname, procname)
	os.Exit(1)
}

func log(s string, vals ...interface{}) {
	fmt.Fprintf(os.Stderr, s, vals...)
}

func trace(s string, vals ...interface{}) {
	if TRACE {
		log(s, vals...)
	}
}
