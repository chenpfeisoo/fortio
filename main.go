package main

import (
	"bufio"
	"io/ioutil"
	"k8s.io/klog"
	"os"
	"os/exec"
	"strings"
)

type Testdata map[string]string

const (
	FortioServerUrl = "http://fortioserver:8080/echo"
	fortioclient    = "fortioclient"
	code            = "\"Code 200\""
	target          = "\"# target\""
	target50        = "50%"
	target75        = "75%"
	target90        = "90%"
	target99        = "99%"
	target999       = "99.9%"
	time            = "5s"
	//每一组测试数据跑的次数
	count = 3
)

func main() {
	fortio(makeTestdata())
}
func fortio(data Testdata) {
	if data == nil {
		panic("测试数据构造不完整，请重新配置")
	}
	for args1, args2 := range data {
		klog.Infof("测试组合 connection :%v ，Qps :%v：\n",args1,args2)
		for i := 1; i <= count; i++ {
			kubectlcmd := "kubectl exec " + getPodname(fortioclient) + " -- fortio load -c " + args1 + " -qps " + args2 + " -t " + time + " " + FortioServerUrl
			output, err := exec.Command("/bin/sh", "-c", kubectlcmd).CombinedOutput()
			if err != nil {
				panic("fortio测试失败" + err.Error())
			}
			//执行日志写入log
			writeLog(string(output))
		}
		//分析结果
		show()
		deletelog()
	}
}

func getPodname(podregex string) string {
	kubectlcmd := "kubectl get pod -A | grep " + podregex + "| awk '{ print $2  }'"
	output, err := exec.Command("/bin/sh", "-c", kubectlcmd).CombinedOutput()
	if err != nil {
		panic("访问k8s失败，检查集群是否正常运行，是否开启外网" + err.Error())
	}
	return strings.Replace(string(output), "\n", "", -1)
}

//构造测试数据
func makeTestdata() Testdata {
	return Testdata{
		"1":  "1000",
		"2":  "1000",
		"4":  "1000",
		"8":  "1000",
		"16": "1000",
		"32": "1000",
		"64": "1000",
	}
}

//结果展示
func show() {
	klog.Info("测试结果如下\n")
	result(target50)
	result(target75)
	result(target90)
	result(target99)
	result(target999)
}

//执行数据分析
func result(metric string) {
	a := func(metric string) {
		cmd := "cat ./log | grep \"# target\"  | grep \"" + metric + "\"| awk '{sum+=$NF}END{print \"" + metric + " Average = \", sum/NR}'"
		//klog.Info(cmd)
		output, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
		if err != nil {
			panic(cmd + "执行失败" + err.Error())
		}
		klog.Info(string(output))
	}

	switch metric {
	case target50:
		a(target50)
	case target75:
		a(target75)
	case target90:
		a(target90)
	case target99:
		a(target99)
	case target999:
		a(target999)
	}

}

//过滤数据
func grepStr(str string) string {
	cmd := "cat " + logfile() + " | grep " + str
	output, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
	klog.Infof("cmd: %v", cmd)
	if err != nil {
		panic("执行" + cmd + "失败" + err.Error())
	}
	if string(output) != "" {
		return string(output)
	}
	panic("fortio测试失败")
}

//写日志
func writeLog(str string) {
	_, err := ioutil.ReadFile(logfile())
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		_, err = os.Create(logfile())
		if err != nil {
			panic("创建文件失败" + err.Error())
		}
	}
	filePath := logfile()
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic("打开文件失败" + err.Error())
	}
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.WriteString(str)
	write.Flush()
}

func logfile() string {
	return "./log"
}

func deletelog() {
	err := os.Remove(logfile())
	if err != nil {
		panic(err.Error())
	}
}
