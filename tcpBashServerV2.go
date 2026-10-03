package main

import (
	"fmt"
//	"io"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/creack/pty"
	util "github.com/prr123/utility/utilLib"
)

var dbg bool

func handleConnection(conn net.Conn) {
	inp := make([]byte, 1024)
	out := make([]byte, 1024)


	defer conn.Close()
	fmt.Printf("[+] New connection from %s\n", conn.RemoteAddr().String())

	// 1. Fire up a brand-new bash command instance
	cmd := exec.Command("/bin/bash", "-i")

	// 2. Start the command inside a Pseudoterminal (PTY)
	ptyMaster, err := pty.Start(cmd)
	if err != nil {
		log.Printf("[-] Failed to start PTY: %v\n", err)
		return
	}
	defer ptyMaster.Close()
	defer cmd.Process.Kill() // Make sure bash closes if the client disconnects

	n, err := conn.Read(inp)
	if err != nil {log.Fatalf("error -- read conn: %v\n", err)}
	fmt.Printf("conn read[%d]: %s\n", n, inp)

	_, err = ptyMaster.Write(inp[:n])
	if err != nil {log.Fatalf("error -- pty write: %v\n", err)}

	no, err :=ptyMaster.Read(out)
	if err != nil {log.Fatalf("error -- pty read: %v\n", err)}
	fmt.Printf(" resp [%d] %s", no, out[:no])

	inp = []byte("\n")
	_, err = ptyMaster.Write(inp[:1])
	if err != nil {log.Fatalf("error -- pty write nl: %v\n", err)}

	no, err = ptyMaster.Read(out)
	if err != nil {log.Fatalf("error -- pty read ll resp: %v\n", err)}
	fmt.Printf(" resp2 [%d] %s", no, out[:no])

	no, err = ptyMaster.Read(out)
	if err != nil {log.Fatalf("error -- pty read ll resp: %v\n", err)}
	fmt.Printf(" resp3 [%d] %s\n", no, out[:no])

/*
	// 3. Pipe data bidirectionally between the TCP socket and the PTY master
	// Create channels to track when the copying routines finish
	done := make(chan struct{}, 2)

	// Copy data coming from the user (TCP) into the PTY (Bash)
	go func() {
		n, err := io.Copy(ptyMaster, conn)
		if err != nil {log.Fatalf("error from client: %v\n", err)}
		log.Printf("rec %d\n", n)
		done <- struct{}{}
	}()

	// Copy data outputted from the PTY (Bash) back to the user (TCP)
	go func() {
		m, err := io.Copy(conn, ptyMaster)
		if err != nil {log.Fatalf("error to client: %v\n", err)}
		log.Printf("sent %d\n", m)
		done <- struct{}{}
	}()

	// Wait for either stream to drop out
	<-done
*/

	fmt.Printf("[-] Connection closed for %s\n", conn.RemoteAddr().String())
}

func main() {

    numArgs := len(os.Args)

    flags:=[]string{"dbg","port", "serv"}

    useStr := "/serv=adr /port=<portnum> [/dbg]"
    helpStr := fmt.Sprintf("help: tcp server V2\n")

    if numArgs > len(flags)+1 {
        fmt.Println("too many arguments in cl!")
        fmt.Println("usage: %s %s\n", os.Args[0], useStr)
        os.Exit(1)
    }

   if numArgs == 1 {
        fmt.Printf("usage is: %s\n", useStr)
        fmt.Printf("%s\n", helpStr)
        os.Exit(1)
    }

   if numArgs == 2 {
        if os.Args[1] == "help" {
            fmt.Printf("usage is: %s\n", useStr)
            fmt.Printf("%s\n", helpStr)
            os.Exit(1)
        }
    }


    flagMap, err := util.ParseFlags(os.Args, flags)
    if err != nil {log.Fatalf("util.ParseFlags: %v\n", err)}

    dbg = false
    _, ok := flagMap["dbg"]
    if ok {dbg = true}

    portStr := ""
    pval, ok := flagMap["port"]
    if ok {
        if pval.(string) == "none" {log.Fatalf("error -- no port number provided with /port flag!")}
        portStr = pval.(string)
    }

    sadrStr := ""
    sval, ok := flagMap["serv"]
    if ok {
        if sval.(string) == "none" {log.Fatalf("error -- no serve IP provided with /serve flag!")}
        sadrStr = sval.(string)
    }

    servAdr := sadrStr + ":" + portStr
    if dbg {
        fmt.Printf(" serv IP:  %s\n", sadrStr)
        fmt.Printf(" port:     %s\n", portStr)
        fmt.Printf(" serv Adr: %s\n", servAdr)
    }

	listener, err := net.Listen("tcp", servAdr)
	if err != nil {
		log.Fatalf("[-] Failed to bind server: %v\n", err)
	}
	defer listener.Close()

	log.Printf("*** TCP PTY Server listening on %s\n", servAdr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error -- listener accept error: %v\n", err)
			continue
		}

		// Handle each interactive shell safely in its own goroutine
		go handleConnection(conn)
	}
}

