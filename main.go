package main

import (
	"fmt"
	"os"
	"time"

	"github.com/monitoring-forge/flagrun"
	"github.com/montanaflynn/stats"
)

var version string

type Opt struct {
	Host      string `long:"host" description:"Target IP address to ping" required:"true"`
	Timeout   int    `long:"timeout" default:"1000" description:"timeout millisec per ping"`
	Interval  int    `long:"interval" default:"10" description:"sleep millisec after every ping"`
	Count     int    `long:"count" default:"10" description:"Count Sending ping"`
	KeyPrefix string `long:"key-prefix" description:"Metric key prefix" required:"true"`
	Version   bool   `short:"v" long:"version" description:"Show version"`
}

func (opt *Opt) Run(_ []string) (any, int) {
	ra, pinger, err := resolveIPAddrAndPinger(opt.Host)

	if err != nil {
		errorNow := uint64(time.Now().Unix())
		fmt.Printf("pinging.%s_rtt_count.success\t%f\t%d\n", opt.KeyPrefix, 0.0, errorNow)
		fmt.Printf("pinging.%s_rtt_count.error\t%f\t%d\n", opt.KeyPrefix, float64(opt.Count), errorNow)
		return err, flagrun.CRITICAL
	}
	defer pinger.Close()

	rtts := []float64{}
	succeeded := float64(0)
	failed := float64(0)

	// preflight
	_, err = pinger.Ping(ra, time.Millisecond*time.Duration(opt.Timeout))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error in preflight: %v\n", err)
	}

	for i := 0; i < opt.Count; i++ {
		time.Sleep(time.Millisecond * time.Duration(opt.Interval))
		rtt, err := pinger.Ping(ra, time.Millisecond*time.Duration(opt.Timeout))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in ping: %v\n", err)
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
			fmt.Fprintf(os.Stderr, "error in calculating average: %v\n", err)
		}
		min, err := stats.Min(rtts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in calculating min: %v\n", err)
		}
		max, err := stats.Max(rtts)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in calculating max: %v\n", err)
		}
		percentile90, err := stats.Percentile(rtts, 90)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error in calculating 90th percentile: %v\n", err)
		}

		fmt.Printf("pinging.%s_rtt_ms.max\t%f\t%d\n", opt.KeyPrefix, max, now)
		fmt.Printf("pinging.%s_rtt_ms.min\t%f\t%d\n", opt.KeyPrefix, min, now)
		fmt.Printf("pinging.%s_rtt_ms.average\t%f\t%d\n", opt.KeyPrefix, mean, now)
		fmt.Printf("pinging.%s_rtt_ms.90_percentile\t%f\t%d\n", opt.KeyPrefix, percentile90, now)
	}
	return "", flagrun.OK
}

func main() {
	os.Exit(flagrun.Go(&Opt{}, flagrun.Version(version)))
}
