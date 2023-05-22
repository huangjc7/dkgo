package dk

import "log"

func (d *dkController) deployKubernetesMaster(ipList []string) {
	log.Println("马上开始安装相关组件 这可能需要一些时间 等耐心等候")
	d.runCmd("yum -y install docker-ce kubelet-1.27.2 kubeadm-1.27.2 kubectl-1.27.2 && "+
		"systemctl enable docker kubelet && "+
		"systemctl start docker", ipList)
	log.Println("docker | kubernetes完成安装")
	//安装cri-dockerd
	d.runCmd("\\cp -f cri-dockerd /usr/local/bin/ && chmod 0755 /usr/local/bin/cri-dockerd", ipList)
	d.runCmd("cat > /usr/lib/systemd/system/cri-dockerd.service <<EOF\n[Unit]\nDescription=CRI Interface for Docker Application Container Engine\nDocumentation=https://docs.mirantis.com\nAfter=network-online.target firewalld.service docker.service\nWants=network-online.target\nRequires=cri-docker.socket\n\n[Service]\nType=notify\nExecStart=/usr/local/bin/cri-dockerd --container-runtime-endpoint fd://\nExecReload=/bin/kill -s HUP $MAINPID\nTimeoutSec=0\nRestartSec=2\nRestart=always\nStartLimitBurst=3\nStartLimitInterval=60s\nLimitNOFILE=infinity\nLimitNPROC=infinity\nLimitCORE=infinity\nTasksMax=infinity\nDelegate=yes\nKillMode=process\n\n[Install]\nWantedBy=multi-user.target\nEOF", ipList)
	d.runCmd("cat > /usr/lib/systemd/system/cri-dockerd.socket <<EOF\n[Unit]\nDescription=CRI Docker Socket for the API\nPartOf=cri-docker.service\n[Socket]\nListenStream=%t/cri-dockerd.sock\nSocketMode=0660\nSocketUser=root\nSocketGroup=docker\n[Install]\nWantedBy=sockets.target\nEOF", ipList)
	d.runCmd("systemctl enable cri-docker.service && systemctl enable --now cri-docker.socket", ipList)

	d.runCmd("kubeadm init "+
		//"--image-repository registry.cn-hangzhou.aliyuncs.com/google_containers "+
		"--pod-network-cidr=10.244.0.0/16 --service-cidr=10.96.0.0/12 "+
		" --cri-socket=unix:///var/run/cri-dockerd.sock"+
		" --v=5", ipList)
	d.runCmd("mkdir -p $HOME/.kube &&  \\cp -f /etc/kubernetes/admin.conf $HOME/.kube/config && chown $(id -u):$(id -g) $HOME/.kube/config", ipList)
}

func (d *dkController) deployKubernetesNode() {

}
