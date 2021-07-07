package metricHandler

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"log"
	"os/exec"
	"shell-exporter/pkgs/configParse"
	"strconv"
	"strings"
	"time"
)

func newCounter(config *configParse.MetricConfig) {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name:        config.MetricName,
		Help:        config.MetricHelp,
		ConstLabels: config.Labels,
	})

	// 注册metric
	prometheus.MustRegister(counter)

	// 根据配置的取值间隔，执行metric操作，counter的值为脚本取值的累加值
	for range time.Tick(time.Duration(config.MetricInterval) * time.Second) {
		counter.Add(cmdRun(config.Script))
	}
}

func newGauge(config *configParse.MetricConfig) {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        config.MetricName,
		Help:        config.MetricHelp,
		ConstLabels: config.Labels,
	})

	// 注册metric
	prometheus.MustRegister(gauge)

	// 根据配置的取值间隔，执行metric操作，gauge的值为脚本取值的当前值
	for range time.Tick(time.Duration(config.MetricInterval) * time.Second) {
		gauge.Set(cmdRun(config.Script))
	}
}

func cmdRun(cmd string) float64 {
	cmdObj := exec.Command("sh", "-c", cmd)

	output, err := cmdObj.Output()
	if err != nil {
		log.Printf("执行命令:%s 失败！Error: %v\n", cmd, err)
		// 执行错误，返回-999
		return float64(-999)
	}

	//将output []byte转为float64
	ret, err := strconv.ParseFloat(strings.TrimSuffix(string(output), "\n"), 64)
	if err != nil {
		log.Printf("执行结果转换float64失败，原始值：%s, error: %v\n", strings.TrimSuffix(string(output), "\n"), err)
		// 执行错误，返回-999
		return float64(-999)
	}

	log.Printf("执行命令:%s 成功，结果为:%f\n", cmd, ret)
	return ret
}

func Metric(configFile string, defaultMetric bool) {
	if !defaultMetric {
		prometheus.Unregister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
		prometheus.Unregister(collectors.NewGoCollector())
	}

	metricConfigList := configParse.GetMetricConfigList(configFile)
	for _, cf := range metricConfigList {
		log.Printf("metric配置：%v\n", cf)

		if cf.MetricType == "counter" {
			log.Printf("启动metric：%s\n", cf.MetricName)
			go newCounter(cf)
		} else {
			log.Printf("启动metric：%s\n", cf.MetricName)
			go newGauge(cf)
		}
	}
}
