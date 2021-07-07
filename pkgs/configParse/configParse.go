package configParse

import (
	"encoding/json"
	"gopkg.in/ini.v1"
	"log"
)

type MetricConfig struct {
	MetricName     string
	Script         string
	Labels         map[string]string
	MetricType     string
	MetricHelp     string
	MetricInterval int64
}

func GetMetricConfigList(configPath string) (ret []*MetricConfig) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Fatal("加载配置文件失败", err)
	}

	for _, sec := range cfg.SectionStrings() {
		if sec == "DEFAULT" {
			continue
		}
		log.Printf("开始加载%s配置段\n", sec)

		labelsMap := make(map[string]string)
		if cfg.Section(sec).HasKey("labels") {
			labels := cfg.Section(sec).Key("labels").String()
			err := json.Unmarshal([]byte(labels), &labelsMap)
			if err != nil {
				log.Printf("label进行json解析失败，%v\n", err)
				continue
			}
		}

		mc := &MetricConfig{
			MetricName:     sec,
			Script:         cfg.Section(sec).Key("script").MustString("default"),
			MetricType:     cfg.Section(sec).Key("metricType").In("gauge", []string{"gauge", "counter"}),
			MetricHelp:     cfg.Section(sec).Key("metricHelp").MustString(sec),
			MetricInterval: cfg.Section(sec).Key("metricInterval").MustInt64(300),
			Labels:         labelsMap,
		}
		ret = append(ret, mc)
	}
	return
}
