package dk

type dkController struct {
	Master      string
	Node        string
	StopCh      chan int
	HostNameCmd string
}

var Dk dkController
