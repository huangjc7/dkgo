package dk

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

func (d *dkController) CheckOS() {
	log.Println("操作系统:", runtime.GOOS)
}

func (d *dkController) CheckRootUser() {
	if os.Geteuid() != 0 {
		log.Fatal("请以 root 用户身份运行此工具")
	}
}

func (d *dkController) OffSwap(ipList []string) {
	d.runCmd("swapoff -a", ipList)
	d.runCmd("sed -i '/swap/s/^/#/' /etc/fstab", ipList)
	log.Println("已禁用swap分区")
}

func (d *dkController) OffSelinux(ipList []string) {
	d.runCmd("setenforce 0", ipList)
	d.runCmd("sed -i 's/SELINUX=enforcing/SELINUX=disabled/' /etc/selinux/config", ipList)
	log.Println("已禁用selinux")
}

func (d *dkController) OffFirewalld(ipList []string) {
	d.runCmd("systemctl stop firewalld", ipList)
	d.runCmd("systemctl disable firewalld", ipList)
	log.Println("已禁用firewalld防火墙")
}

func (d *dkController) setHostName(ipList []string) {
	for _, v := range ipList {
		if v == d.Master {
			d.HostNameCmd = fmt.Sprintf("hostnamectl set-hostname %s.master",
				strings.Replace(v, ".", "-", -1))
			d.runCmd(d.HostNameCmd, ipList[:1])
		} else {
			d.HostNameCmd = fmt.Sprintf("hostnamectl set-hostname %s.node",
				strings.Replace(v, ".", "-", -1))
			d.runCmd(d.HostNameCmd, ipList[1:])
		}
	}
	log.Println("已修改主机名")
}

func (d *dkController) editHosts(ipList []string) {
	hostsCmd := fmt.Sprintf("echo \"%s %s\" >> /etc/hosts", d.Master, "apiserver.cluster.local")
	d.runCmd(hostsCmd, ipList)
	log.Println("域名已经注入hosts文件")
}

func (d *dkController) yumRepo(ipList []string) {
	d.runCmd("yum install -y yum-utils device-mapper-persistent-data lvm2 && "+
		"yum-config-manager --add-repo https://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo && "+
		"sed -i 's+download.docker.com+mirrors.aliyun.com/docker-ce+' /etc/yum.repos.d/docker-ce.repo && "+
		"cat <<EOF > /etc/yum.repos.d/kubernetes.repo\n[kubernetes]\nname=Kubernetes\nbaseurl=https://mirrors.aliyun.com/kubernetes/yum/repos/kubernetes-el7-x86_64/\nenabled=1\ngpgcheck=1\nrepo_gpgcheck=1\ngpgkey=https://mirrors.aliyun.com/kubernetes/yum/doc/yum-key.gpg https://mirrors.aliyun.com/kubernetes/yum/doc/rpm-package-key.gpg\nEOF",
		ipList)
	log.Println("已添加yum源")
}
