package main

import (
	"flag"
	"log"
	"net/http"
	"shell-exporter/pkgs/metricHandler"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	listenAddr := flag.String("listenAddress", "0.0.0.0", "Exporter的HTTP Server监听的地址")
	listenPort := flag.String("listenPort", "9527", "Exporter的HTTP Server监听的端口")
	metricConfigFile := flag.String("metricConfigFile", "./config.ini", "指标配置文件目录")
	defaultMetric := flag.Bool("defaultMetric", true, "是否显示Exporter默认输出的go和process指标信息")

	flag.Parse()
	addr := *listenAddr + ":" + *listenPort
	// 处理metric
	metricHandler.Metric(*metricConfigFile, *defaultMetric)

	// 提供web页面
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
