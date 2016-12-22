hadoop on Runv
==============

### 部署
假设我们有5台物理机，并且这些物理机**在同一个子网**。

#### 1. 开启 Intel VT 或者 AMD-V virtualization hardware extensions in BIOS

>Intel VT 可能会被 BIOS disable 掉了，执行命令 `cat /proc/cpuinfo | grep vmx svm` ，如果没有输出则表明，VT 被disable掉了，需求重启机器进入BIOS enable VT。

>具体设置方式请参考：https://docs.fedoraproject.org/en-US/Fedora/13/html/Virtualization_Guide/sect-Virtualization-Troubleshooting-Enabling_Intel_VT_and_AMD_V_virtualization_hardware_extensions_in_BIOS.html

#### 2. 安装 KVM 虚拟化组件

>通过如下命令安装 kvm 组件

>```bash
yum install -y qemu-kvm qemu-img virt-manager libvirt-python libvirt-client virt-viewer libvirt virt-install bridge-utils
systemctl start libvirtd
systemctl enable libvirtd
```

>执行下面命令 `lsmod | grep kvm` 如果输出

>```
kvm_intel             162153  0
kvm                   525409  1 kvm_intel
```

>则表明安装成功。

>参考链接： http://www.linuxtechi.com/install-kvm-hypervisor-on-centos-7-and-rhel-7/

#### 3. 安装 hyper

>```bash
curl -sSL https://hypercontainer.io/install | bash
```

>或者参考链接 https://docs.hypercontainer.io/get_started/install/linux.html

#### 4. 配置 linux bridge
>假设物理网卡为 `eth0`, 且主机IP为`192.168.1.20`，子网掩码为`255.255.255.0`, 编辑文件 `vim /etc/sysconfig/network-scripts/ifcfg-eth0`

>至少将下述各值配置好：

>```
TYPE=Ethernet
BOOTPROTO=static
IPADDR=192.168.1.20
ONBOOT=yes
BRIDGE=hyper0
```

> 然后编辑文件 `vim /etc/sysconfig/network-scripts/ifcfg-hyper0` , 将下述各值配置好：

>```
DEVICE=hyper0
TYPE=Bridge
BOOTPROTO=static
IPADDR=192.168.1.20
NETMASK=255.255.255.0
GATEWAY=192.168.1.1
ONBOOT=yes
```

#### 5. 执行 `service network restart`

#### 6. 配置 hyper
> 编辑文件 `vim /etc/hyper/config` , 将其中的相关值设置成如下：

>```
Hypervisor=kvm
BridgeIP=192.168.1.20/24
```

> 然后重启 hyperd ， `service hyperd restart`

#### 7. 在其他几台物理机上依次执行上述 1-6 步，注意要把IP地址换成相应的IP。

#### 8. 部署hdfs name node， 假设我们在 192.168.1.20 机器上部署 namenode，则我们可以直接用pod文件 [nn200.json](./pod/nn200.json)

>这个pod文件表明会启动一个

>* 4vcpu，8192M内存
>* 固定IP为192.168.1.200
>* 并挂载物理机目录 /hadoop/nn200 到 /var/hdfs/namenode
>* 使用gateway 192.168.1.20

> 的 VM.

> 使用命令 `hyperctl run -p pod/nn200.json` 启动VM实例。

#### 9. 部署 hdfs data node，与上述第8步类似， 我们需要编辑 [dn207.json](./pod/dn207.json) 来依次启动多个 data node。

#### 10. 部署 yarn resource manager， 假设我们在 192.168.1.20 机器上部署 resource manager，我们可以直接使用pod文件 [resourcemanager210.json](./pod/resourcemanager210.json)

>这个pod文件表明会启动一个

>* 8vcpu，16384M 内存
>* 固定IP为192.168.1.210
>* 并挂载物理机目录 /hadoop/rm210 到 /tmp/hadoop
>* 使用gateway 192.168.1.20

> 的 VM.

> 使用命令 `hyperctl run -p pod/resourcemanager210.json` 启动VM实例。

#### 11. 部署 yarn node manager，与上述第10步类似， 我们需要编辑 [nodemanager211.json](./pod/nodemanager211.json) 来依次启动多个 node manager。

另外，我们使用的docker镜像的dockerfile 都在 https://github.com/vitan/hadoop/tree/master/centos/7/hadoop


## perf

The cluster statistics:

* 5 bare metals
* KVM
* Linux Bridge
* hdfs on Runv: 4 VCores, 8G * (1 namenode + 3 datanode)/physical machine
* yarn on Runv: 8VCores, 16G * (1 resource manager + 4 node manager) / physical machine

#### in house benchmark:

Terasort in 100G data costed:

```
real    601m25.369s
user    0m56.305s
sys     0m21.853s
```

#### Yahoo benchmark:

In May 2008, Owen O'Malley ran this code on a 910 node cluster and sorted the 10 billion records (1 TB) in 209 seconds (3.48 minutes) to win the annual general purpose (daytona) terabyte sort benchmark.

* 910 nodes
* 4 dual core Xeons @ 2.0ghz per a node
* 4 SATA disks per a node
* 8G RAM per a node
* 1 gigabit ethernet on each node
* 40 nodes per a rack
* 8 gigabit ethernet uplinks from each rack to the core
* Red Hat Enterprise Linux Server Release 5.1 (kernel 2.6.18)
* Sun Java JDK 1.6.0_05-b13

refer link:

* http://www.michael-noll.com/blog/2011/04/09/benchmarking-and-stress-testing-an-hadoop-cluster-with-terasort-testdfsio-nnbench-mrbench/
* http://hadoop.apache.org/docs/r2.7.3/api/org/apache/hadoop/examples/terasort/package-summary.html
* https://github.com/vitan/hadoop/tree/master/centos/7/hadoop

