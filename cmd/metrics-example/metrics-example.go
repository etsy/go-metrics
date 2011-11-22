package main

import (
	"log"
	"github.com/rcrowley/go-metrics"
	"os"
	"rand"
//	"syslog"
	"time"
)

const fanout = 10

func main() {

	r := metrics.NewRegistry()

	c := metrics.NewCounter()
	r.RegisterCounter("foo", c)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				c.Dec(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				c.Inc(47)
				time.Sleep(400e6)
			}
		}()
	}

	g := metrics.NewGauge()
	r.RegisterGauge("bar", g)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				g.Update(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				g.Update(47)
				time.Sleep(400e6)
			}
		}()
	}

	hc := metrics.NewHealthcheck(func(h metrics.Healthcheck) {
		if 0 < rand.Intn(2) {
			h.Healthy()
		} else {
			h.Unhealthy(os.NewError("baz"))
		}
	})
	r.RegisterHealthcheck("baz", hc)

	s := metrics.NewExpDecaySample(1028, 0.015)
//	s := metrics.NewUniformSample(1028)
	h := metrics.NewHistogram(s)
	r.RegisterHistogram("bang", h)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				h.Update(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				h.Update(47)
				time.Sleep(400e6)
			}
		}()
	}

	m := metrics.NewMeter()
	r.RegisterMeter("quux", m)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				m.Mark(19)
				time.Sleep(300e6)
			}
		}()
		go func() {
			for {
				m.Mark(47)
				time.Sleep(400e6)
			}
		}()
	}

	t := metrics.NewTimer()
	r.RegisterTimer("hooah", t)
	for i := 0; i < fanout; i++ {
		go func() {
			for {
				t.Time(func() { time.Sleep(300e6) })
			}
		}()
		go func() {
			for {
				t.Time(func() { time.Sleep(400e6) })
			}
		}()
	}

	metrics.RegisterRuntimeMemStats(r)
	go func() {
		t := time.NewTicker(5e9)
		for 0 < <-t.C { metrics.CaptureRuntimeMemStats(r, true) }
	}()

	metrics.Log(r, 60, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))

/*
	w, err := syslog.Dial("unixgram", "/dev/log", syslog.LOG_INFO, "metrics")
	if nil != err { log.Fatalln(err) }
	metrics.Syslog(r, 60, w)
*/

}
