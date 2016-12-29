Ambari on Runv & Local repo
===========================


### 设置本地 HDP repo

鉴于客户网络安全问题，我们需要在客户内网环境部署一套 HDP repo， 假设我们在机器 192.168.1.211 上部署该repo，则过程如下：

#### 1. 安装必要软件

>```bash
yum install -y wget vim httpd yum-utils createrepo
```

#### 2. 启动 httpd 服务

>```bash
service httpd start
```

#### 3. 配置 HDP 和 HDP-UTILS repo

>我们可以下述命令从官网下载相应的tar, 并解压到相应目录

>```bash
mkdir /var/www/html/hdp
cd /var/www/html/hdp
wget http://public-repo-1.hortonworks.com/HDP/centos7/2.x/updates/2.3.2.0/HDP-2.3.2.0-centos7-rpm.tar.gz
wget http://public-repo-1.hortonworks.com/HDP-UTILS-1.1.0.20/repos/centos7/HDP-UTILS-1.1.0.20-centos7.tar.gz
tar -xzvf HDP-2.3.2.0-centos7-rpm.tar.gz
tar -xzvf HDP-UTILS-1.1.0.20-centos7.tar.gz
rm HDP-2.3.2.0-centos7-rpm.tar.gz
rm HDP-UTILS-1.1.0.20-centos7.tar.gz
```

> 现在我们用浏览器访问链接 http://192.168.1.211/hdp 可以看到相应的文件了。

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

#### 8. 部署 ambari server， 假设我们在 192.168.1.20 机器上部署 ambari server，则我们可以直接用pod文件 [ambari-server.json](./ambari-server.json)

>这个pod文件表明会启动一个

>* 4vcpu，6096M内存
>* 固定IP为192.168.1.110
>* 并挂载物理机目录 /ambari/ambari110/log 到 /var/log/ambari-server/ , /ambari/ambari110/pgsql 到 /var/lib/pgsql/
>* 使用gateway 192.168.1.20

> 的 VM.

> 使用命令 `hyperctl run -p ambari110.json` 启动VM实例。

#### 9. 部署 agent node，与上述第8步类似， 我们需要编辑 [ambari-agent.json](./ambari-agent.json) 来依次启动多个 agent node。

>这个pod文件表明会启动一个

>* 4vcpu，6096M 内存
>* 固定IP为192.168.1.112
>* 并挂载物理机目录 /ambari/ambari112 到 /lib/modules
>* 使用gateway 192.168.1.20

> 的 VM.

### 进入 ambari server 设置

#### 1. 浏览器访问 192.168.1.110:8080

>* 用户名： admin
>* 密码: admin

#### 2. Launch Install Wizard -> 创建集群（create cluster）

#### 3. 在 **Select Stack** 步

>选择
