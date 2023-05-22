package dk

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
