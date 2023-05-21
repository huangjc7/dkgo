package dk

import "log"

func (d *dkController) Run() {
	ipList := []string{d.Master, d.Node}
	d.initOS(ipList)
	d.deployKubernetesMaster(ipList[:1])
	d.deployKubernetesNode()

}

func (d *dkController) initOS(ipList []string) {
	d.CheckOS()
	d.CheckRootUser()
	d.OffSwap(ipList)
	d.OffSelinux(ipList)
	d.OffFirewalld(ipList)
	d.setHostName(ipList)
	d.editHosts(ipList)
	d.yumRepo(ipList)
}

func (d *dkController) deployKubernetesMaster(ipList []string) {
	log.Println("马上开始安装相关组件 这可能需要一些时间 等耐心等候")
	d.runCmd("yum -y install docker-ce kubelet-1.27.2 kubeadm-1.27.2 kubectl-1.27.2 && "+
		"systemctl enable docker kubelet && "+
		"systemctl start docker", ipList)
	log.Println("docker | kubernetes完成安装")
	//TODO 安装cri流程未写完
	d.runCmd("kubeadm init "+
		"--image-repository registry.cn-hangzhou.aliyuncs.com/google_containers "+
		"--pod-network-cidr=10.244.0.0/16 --service-cidr=10.96.0.0/12 "+
		" --cri-socket=unix:///var/run/cri-dockerd.sock "+
		" --v=5", ipList)
}

func (d *dkController) deployKubernetesNode() {

}
