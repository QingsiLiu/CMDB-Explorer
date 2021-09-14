package main

import "github.com/go-kit/kit/metrics/prometheus"

func main() {
	// counter
	// guage
	prometheus.NewGauge()
	// historgram
	// summary

	// metrice_name{lable=lable_value} metrice_value

	// 有lable
	// label/label_value 固定/不固定
	// 无lable => label都是空 => 固定

}
