package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func main() {
	// counter
	requestTotal := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "request_total",
		Help: "request total",
	})

	codeStatus := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "status_code_total",
		Help: "status code total",
	}, []string{"status"})
	// guage
	// 固定label
	cpu := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "cpu",
		Help:        "cpu total",
		ConstLabels: prometheus.Labels{"a": "xxx"},
	})

	// 非固定label
	disk := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "disk",
		Help: "disk total",
	}, []string{"mount", "name"})
	// historgram
	url := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "request_time",
		Help:    "request time",
		Buckets: prometheus.LinearBuckets(0, 3, 5),
	}, []string{"url"})
	// summary
	requestSummary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name:       "requestSummary",
		Help:       "requestSummary",
		Objectives: map[float64]float64{0.5: 0.05, 0.99: 0.001, 0.9: 0.01},
	}, []string{"url"})

	// metrice_name{lable=lable_value} metrice_value

	// 有lable
	// label/label_value value固定/不固定
	// 无lable => label都是空 => 固定

	cpu.Set(2)

	disk.WithLabelValues("c:", "xxx").Set(100)
	disk.WithLabelValues("e:", "yyyy").Set(200)

	// 注册指标信息
	prometheus.MustRegister(cpu)
	prometheus.MustRegister(disk)
	prometheus.MustRegister(requestTotal)
	prometheus.MustRegister(codeStatus)
	prometheus.MustRegister(url)
	prometheus.MustRegister(requestSummary)

	requestTotal.Add(10)
	codeStatus.WithLabelValues("200").Add(100)
	codeStatus.WithLabelValues("500").Add(2)
	url.WithLabelValues("/aaaaa").Observe(6)
	requestSummary.WithLabelValues("/aaaa").Observe(6)

	// 值的修改, 修改的时间 => 触发（时间、事件触发）
	// 时间触发
	go func() {
		for range time.Tick(time.Second) {
			disk.WithLabelValues("c:", "xxxxx").Set(float64(rand.Int()))
		}
	}()

	// 事件触发
	// 业务请求
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		requestTotal.Inc()
		codeStatus.WithLabelValues(strconv.Itoa(rand.Intn(5) * 100)).Add(1)
		url.WithLabelValues(request.URL.Path).Observe(float64(rand.Intn(20)))
		requestSummary.WithLabelValues(request.URL.Path).Observe(float64(rand.Intn(20)))
		fmt.Fprint(writer, "hi")
	})

	// 在metice接口访问时
	call := prometheus.NewCounterFunc(prometheus.CounterOpts{
		Name: "xxx",
		Help: "xxx",
	}, func() float64 {
		fmt.Println("call")
		return rand.Float64()
	})

	prometheus.MustRegister(call)

	// 暴露
	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":9999", nil)
}
