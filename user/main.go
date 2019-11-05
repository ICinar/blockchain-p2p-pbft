package main

import (
        "bufio"
        "context"
        "crypto/rand"
        "encoding/json"
        "flag"
        "fmt"
        "io"
        "log"
        "time"
        "strings"
	"os"

        "mpse/blockchain"

        "github.com/davecgh/go-spew/spew"
        "github.com/libp2p/go-libp2p"
        "github.com/libp2p/go-libp2p-crypto"
        "github.com/libp2p/go-libp2p-net"
        "github.com/libp2p/go-libp2p-peerstore"
        "github.com/multiformats/go-multiaddr"
         "github.com/libp2p/go-libp2p-host"
)

var Blockchain []blockchain.Block
var TX blockchain.TX

func handleStream(s net.Stream) {
        log.Println("Got a new stream!")
                                          
	       // Create a buffer stream for non blocking read and write.
        rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

    go readData(rw)
	go sendData(rw)

        // stream 's' will stay open until you close it (or the other side closes it).
}
func readData(rw *bufio.ReadWriter) {
        for {
                str, _ := rw.ReadString('\n')

                if str == "" {
                        return
                }
                if str != "\n" {
                        chain := make([]blockchain.Block, 0)
                        if err := json.Unmarshal([]byte(str), &chain); err != nil {
                                log.Fatal(err)
                        }

                        if len(chain) > len(Blockchain) {
                                Blockchain = chain
                                bytes, err := json.MarshalIndent(Blockchain, "", "  ")
                                if err != nil {

                                        log.Fatal(err)
                                }
                                // Green console color:         \x1b[32m
                                // Reset console color:         \x1b[0m
                                fmt.Printf("\x1b[32m%s\x1b[0m> ", string(bytes))
                        }

                }

        }
}

func sendData(rw *bufio.ReadWriter) {

        go func() {
        for {
                time.Sleep(5 * time.Second)

                bytes, err := json.Marshal(Blockchain)
                if err != nil {
                        log.Println(err)
                }



                rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
                rw.Flush()


                }
							        }()


        stdReader := bufio.NewReader(os.Stdin)

        for {
                fmt.Print("> ")
                sendData, err := stdReader.ReadString('\n')
                //sendData := inp
				

                sendData = strings.Replace(sendData, "\n", "", -1)
               TX.TokenTran = sendData
               
			newBlock := blockchain.GenerateBlock(Blockchain[len(Blockchain)-1],TX)

                if blockchain.IsBlockValid(newBlock, Blockchain[len(Blockchain)-1]) {
                        Blockchain = append(Blockchain, newBlock)
                }

                bytes, err := json.Marshal(Blockchain)
                if err != nil {
                        log.Println(err)
                }

	        spew.Dump(Blockchain)


                rw.WriteString(fmt.Sprintf("%s\n", string(bytes)))
                rw.Flush()
        }

}


func makeHost(sourcePort *int) (host.Host,error){


        var r io.Reader

        r = rand.Reader

        prvKey, _, err:= crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
        if err != nil {
                panic(err)
        }

        // 0.0.0.0 will listen on any interface device.
        sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", *sourcePort))
	      // libp2p.New constructs a new libp2p Host.
        // Other options can be added here.
        host, err := libp2p.New(
                context.Background(),
                libp2p.ListenAddrs(sourceMultiAddr),
                libp2p.Identity(prvKey),
        )

        if err != nil {
                panic(err)
        }

       return host,err

}

func main() {
        t := time.Now()	
		signature := "abceufwlf423"
		pubKey := "123456789sdjfsdje"
		TX := blockchain.TX{t.String(),signature,pubKey,""}
	
		genesisBlock := blockchain.Block{}
        genesisBlock = blockchain.Block{t.String(), "", blockchain.CalculateHash(genesisBlock),TX}
        Blockchain = append(Blockchain, genesisBlock)   
	
	
        sourcePort := flag.Int("sp", 0, "Source port number")
        dest := flag.String("d", "", "Destination multiaddr string")
        //help := flag.Bool("help", false, "Display help")
        //debug := flag.Bool("debug", false, "Debug generates the same node ID on every execution")
        flag.Parse()
        
	
        host, err :=makeHost(sourcePort)
        if err != nil {
                log.Fatalln(err)
        }
        
         if *dest == "" {
                host.SetStreamHandler("/chat/1.0.0",handleStream)


                // Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
                var port string
                for _, la := range host.Network().ListenAddresses() {
                        if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
                                port = p
                                break
                        }
                }

                if port == "" {
                        panic("was not able to find actual local port")
                }

                fmt.Printf("Run 'go run main.go -d /ip4/127.0.0.1/tcp/%v/p2p/%s' on another console.\n", port, host.ID().Pretty())
                fmt.Println("You can replace 141.100.74.x with public IP as well.")
                fmt.Printf("\nWaiting for incoming connection\n\n")

                // Hang forever
                <-make(chan struct{})

        } else {
                        fmt.Println("This node's multiaddresses:")
                        for _, la := range host.Addrs() {
                                fmt.Printf(" - %v\n", la)
                        }
                        fmt.Println()

                        // Turn the destination into a multiaddr.
                        maddr, err := multiaddr.NewMultiaddr(*dest)
                        if err != nil {
                                log.Fatalln(err)
                        }

                        // Extract the peer ID from the multiaddr.
                        info, err := peerstore.InfoFromP2pAddr(maddr)
                        if err != nil {
                                log.Fatalln(err)
                        }

                        // Add the destination's peer multiaddress in the peerstore.
                        // This will be used during connection and stream creation by libp2p.
                        host.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

                        // Start a stream with the destination.
                        // Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
                        s, err := host.NewStream(context.Background(), info.ID, "/chat/1.0.0")
                        if err != nil {
                                panic(err)
                        }

                        // Create a buffered stream so that read and writes are non blocking.
                        rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

                        // Create a thread to read and write data.
                        
                 //log.Printf("Test")
						go sendData(rw)
                        go readData(rw)

                        // Hang forever.
                        select {}
                }
        
}






