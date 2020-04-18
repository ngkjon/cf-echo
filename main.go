package main

import(
"fmt"
"flag"
"net"
"time"
"os"
"log"
"golang.org/x/net/icmp"
"golang.org/x/net/ipv4"
)

var usage = "test doc"
var ListenAddr = "0.0.0.0"
var sent = 0
var received = 0
func main() {
  delay := flag.Duration("d",100000000,"sending delay")
  timeout := flag.Duration("t",10000000000,"sending delay")

  flag.Parse()
  host := flag.Arg(0)

  fmt.Println("delay:", *delay)
  fmt.Println("timeout:", *timeout)
	fmt.Println("host:", host)
  fmt.Println("listening:", ListenAddr)

  for{
    Ping(host,*timeout)
    time.Sleep(*delay)
  }
}

func Ping(addr string, timeout time.Duration)(*net.IPAddr, time.Duration, error){
  //start listener
  icmpListener, err := icmp.ListenPacket("ip4:icmp", ListenAddr)
  if err != nil {
      return nil, 0, err
  }
  defer icmpListener.Close()
  dst, err := net.ResolveIPAddr("ip4", addr)
  if err != nil {
      return nil, 0, err
  }

  //create msg
  icmpMsg := icmp.Message{
      Type: ipv4.ICMPTypeEcho, Code: 0,
      Body: &icmp.Echo{
          ID: os.Getpid() & 0xffff, Seq: 1, //<< uint(seq),
          Data: []byte(""),
      },
  }
  icmpMsgB, err := icmpMsg.Marshal(nil)
  if err != nil {
      return dst, 0, err
  }

  //echo send
  start := time.Now()
  _, err = icmpListener.WriteTo(icmpMsgB, dst)
  if err != nil {
      return dst, 0, err
  } 
  sent += 1
  //await response
  response := make([]byte, 1500)

  err = icmpListener.SetReadDeadline(time.Now().Add(timeout))
  if err != nil {
      return dst, 0, err
  }
  _, peer, err := icmpListener.ReadFrom(response)

  if err != nil {
      return dst, 0, err
  }
  received += 1
  loss := sent - received
  //calculate latency
  rtt := time.Since(start)

  log.Printf("Ping %s (%s): %s Sent:%d Loss:%d \n", peer, dst, rtt,sent,loss)
  return dst, rtt, nil
}

