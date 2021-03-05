# fortio
fortio压测 istio，数据分析工具 </br>
需要提前部署好fortio server和 clinet </br>
代码构造运行如下命令 
```shell
    kubectl exec   fortioclientpod -- fortio load -c  args1   -qps   args2  -t  time   FortioServerUrl
```
测试数据 ,每组测试数据执行三次，算均值 </br>
```shell
connection :  qps
		"1":  "1000",
		"2":  "1000",
		"4":  "1000",
		"8":  "1000",
		"16": "1000",
		"32": "1000",
		"64": "1000",
```
结果展示</br>
![image](https://user-images.githubusercontent.com/18147157/110070854-d9b4f600-7db5-11eb-90c2-e23e6d6022e0.png)
