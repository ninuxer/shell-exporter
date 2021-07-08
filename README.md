## 功能说明
- 1、将shell脚本执行结果输出为prometheus exporter，以供prometheus进行采集并监控
- 2、如果脚本执行失败，则metric默认值为-999
- 3、支持附加label信息
- 4、输出指标类型默认为gauge，支持gauge、counter
    - gauge的计算逻辑为取脚本的执行结果(标准输出)作为其值
    - counter的计算逻辑为取脚本的执行结果(标准输出)作为值，进行累加

## 参数说明
- --listenAddr string 指定exporter程序监听的端口，默认0.0.0.0
- --listenPort string 指定exporter程序监听的端口，默认9527
- --metricConfigFile string  指定需要监控的指标的配置文件路径，默认./config.ini
- --defaultMetric bool 是否展示Exporter默认输出的go和process指标信息，默认为true,即默认展示

## 配置文件格式说明
```ini
[sh_exporter_test] # 配置段名称，作为metric名称
script=/path/to/script.sh # 执行的脚本，可以是命令，也可以是自定义脚本
labels={"label1": "aaa", "label2": "222"} # 指标附加label, 可选，默认空
metricType=gauge # 指标类型，支持gauge、counter，可选，默认gauge
metricHelp="XXX" # 指标帮助信息，描述指标功能，可选，默认帮助信息为指标名称
meticInterval=300 # 指标采集间隔，单位秒，默认300
[xxx]
[yyy]
...
[zzz]
```

## 执行方式（Linux系统）
- 已编译的二进制
```shell
wget https://github.com/ninuxer/shell-exporter/blob/master/cmd/shell-exporter && chmod +x shell-exporter
参考https://github.com/ninuxer/shell-exporter/blob/master/example/config.ini进行配置文件编写

运行: ./shell-exporter --metricConfigFile=./config.ini --listenPort=9527 --defaultMetric=false 

```

- 源码编译
```shell
git clone git@github.com:ninuxer/shell-exporter.git

cd shell-exporter
go build

参考https://github.com/ninuxer/shell-exporter/blob/master/example/config.ini进行配置文件编写

运行: ./shell-exporter --metricConfigFile=./config.ini --listenPort=9527 --defaultMetric=false 
```

## Docker镜像
- Dockerfile(先下载编译好的shell-exporter，或自己手动编译后使用)
> 注意：如果使用自己编译的版本(在CentOS或Ubuntu系统上编译)，使用动态链接方式编译了一个使用了 GLIBC 库的程序生成的二进制，但是 alpine 镜像中没有 GLIBC 库而是用的 MUSL LIBC 库，这样就会导致该二进制文件无法被执行，可以在编译时，加上CGO_ENABLED=0进行编译，或在alpine中编译一个依赖MUSL LIBC的版本后使用
```dockerfile
FROM alpine:3.14
WORKDIR /shell-exporter
RUN sed -i s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@ /etc/apk/repositories && \
    apk add mysql-client curl iproute2 tzdata --no-cache && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone
ENV TZ Asia/Shanghai
COPY ./shell-exporter .
CMD ["./shell-exporter","--metricConfigFile=/shell-exporter/conf/config.ini", "--listenPort=9527", "--defaultMetric=false"]
```

- Docker方式启动(注意：将配置和脚本挂载进容器。自己使用时，注意脚本内容是否在容器内可正常执行)
```shell
docker run -d --name test-shell-exporter -v /tmp/dockerfiles/shell-exporter/examples/conf:/shell-exporter/conf -v /tmp/dockerfiles/shell-exporter/examples/scripts:/shell-exporter/scripts/ ninuxer/shell-exporter:v1.0
```

