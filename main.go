package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

var concurrency int
var customDNS string
var dnsPort string
var verbose bool

func worker(ips <-chan string, resolver *net.Resolver, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()
	for ip := range ips {
		hostnames, err := resolver.LookupAddr(context.Background(), ip)
		if err != nil {
			continue
		}
		for _, hostname := range hostnames {
			result := ""
			if verbose {
				result = fmt.Sprintf("%s: %s", ip, strings.TrimSuffix(hostname, "."))
			} else {
				result = fmt.Sprintf("%s", strings.TrimSuffix(hostname, "."))
			}
			results <- result
		}
	}
}

func main() {
	flag.StringVar(&customDNS, "dns", "8.8.8.8", "Custom DNS resolver address (ip only)")
	flag.IntVar(&concurrency, "c", 20, "Set the concurrency level")
	flag.StringVar(&dnsPort, "port", "53", "DNS server port")
	flag.BoolVar(&verbose, "v", false, "See IP and Hostname as output")
	flag.Parse()

	ips := make(chan string, concurrency)
	results := make(chan string, concurrency)

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return net.Dial("udp", customDNS+":"+dnsPort)
		},
	}

	var wg sync.WaitGroup

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go worker(ips, resolver, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// Handle either a single IP or CIDR block from stdin
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			input := scanner.Text()
			if strings.Contains(input, "/") {
				// CIDR Block
				ip, ipnet, err := net.ParseCIDR(input)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Invalid CIDR block: %s\n", input)
					continue
				}
				for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incrementIP(ip) {
					ips <- ip.String()
					if ip.Equal(lastAddress(ipnet)) {
						break
					}
				}
			} else {
				// Single IP
				ips <- input
			}
		}
		close(ips)
	}()

	for result := range results {
		fmt.Println(result)
	}
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func lastAddress(n *net.IPNet) net.IP {
	var last net.IP
	for i := range n.IP {
		last = append(last, n.IP[i]|^n.Mask[i])
	}
	return last
}
