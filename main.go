package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	ping "github.com/digineo/go-ping"
	flags "github.com/jessevdk/go-flags"
	"github.com/montanaflynn/stats"
)

var version string
var commit string

type Opt struct {
	Host      string `long:"host" description:"Target IP address to ping" required:"true"`
	Timeout   int    `long:"timeout" default:"1000" description:"timeout millisec per ping"`
	Interval  int    `long:"interval" default:"10" description:"sleep millisec after every ping"`
	Count     int    `long:"count" default:"10" description:"Count Sending ping"`
	KeyPrefix string `long:"key-prefix" description:"Metric key prefix" required:"true"`
	Version   bool   `short:"v" long:"version" description:"Show version"`
}

func resolveIPAddrAndPinger(addr string) (*net.IPAddr, *ping.Pinger, error) {
	var ra *net.IPAddr
	var pinger *ping.Pinger
	if strings.Index(addr, ":") != -1 {
		r, err := net.ResolveIPAddr("ip6", addr)
		if err != nil {
			return ra, pinger, err
		}
		ra = r
		p, err := ping.New("", "::")
		if err != nil {
			return ra, pinger, err
		}
		pinger = p
	} else {
		r, err := net.ResolveIPAddr("ip4", addr)
		if err != nil {
			return ra, pinger, err
		}
		ra = r
		p, err := ping.New("0.0.0.0", "")
		if err != nil {
			return ra, pinger, err
		}
		pinger = p
	}
	return ra, pinger, nil
}

func (opt *Opt) run() error {
	ra, pinger, err := resolveIPAddrAndPinger(opt.Host)

	if err != nil {
		errorNow := uint64(time.Now().Unix())
		fmt.Printf("pinging.%s_rtt_count.success\t%f\t%d\n", opt.KeyPrefix, 0.0, errorNow)
		fmt.Printf("pinging.%s_rtt_count.error\t%f\t%d\n", opt.KeyPrefix, float64(opt.Count), errorNow)
		return err
	}
	defer pinger.Close()

	rtts := []float64{}
	succeeded := float64(0)
	failed := float64(0)

	// preflight
	_, err = pinger.Ping(ra, time.Millisecond*time.Duration(opt.Timeout))
	if err != nil {
		log.Printf("error in preflight: %v", err)
	}

	for i := 0; i < opt.Count; i++ {
		time.Sleep(time.Millisecond * time.Duration(opt.Interval))
		rtt, err := pinger.Ping(ra, time.Millisecond*time.Duration(opt.Timeout))
		if err != nil {
			log.Printf("%v", err)
			failed++
			continue
		}
		rttMilliSec := float64(rtt.Nanoseconds()) / 1000.0 / 1000.0
		rtts = append(rtts, rttMilliSec)
		succeeded++
	}

	now := uint64(time.Now().Unix())
	fmt.Printf("pinging.%s_rtt_count.success\t%f\t%d\n", opt.KeyPrefix, succeeded, now)
	fmt.Printf("pinging.%s_rtt_count.error\t%f\t%d\n", opt.KeyPrefix, failed, now)
	if len(rtts) > 0 {
		mean, err := stats.Mean(rtts)
		if err != nil {
			log.Printf("error in calculating average: %v", err)
		}
		min, err := stats.Min(rtts)
		if err != nil {
			log.Printf("error in calculating min: %v", err)
		}
		max, err := stats.Max(rtts)
		if err != nil {
			log.Printf("error in calculating max: %v", err)
		}
		percentile90, err := stats.Percentile(rtts, 90)
		if err != nil {
			log.Printf("error in calculating 90th percentile: %v", err)
		}

		fmt.Printf("pinging.%s_rtt_ms.max\t%f\t%d\n", opt.KeyPrefix, max, now)
		fmt.Printf("pinging.%s_rtt_ms.min\t%f\t%d\n", opt.KeyPrefix, min, now)
		fmt.Printf("pinging.%s_rtt_ms.average\t%f\t%d\n", opt.KeyPrefix, mean, now)
		fmt.Printf("pinging.%s_rtt_ms.90_percentile\t%f\t%d\n", opt.KeyPrefix, percentile90, now)
	}
	return nil
}

func main() {
	os.Exit(_main())
}

func _main() int {
	opt := &Opt{}
	psr := flags.NewParser(opt, flags.HelpFlag|flags.PassDoubleDash)
	_, err := psr.Parse()

	if opt.Version {
		if commit == "" {
			commit = "dev"
		}
		fmt.Printf(
			"%s-%s\n%s/%s, %s, %s\n",
			filepath.Base(os.Args[0]),
			version,
			runtime.GOOS,
			runtime.GOARCH,
			runtime.Version(),
			commit)
		return 0
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	err = opt.run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}
