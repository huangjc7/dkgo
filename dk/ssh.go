package dk

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

func (d *dkController) runCmd(cmd string, ipList []string) {
	sshUser := "root"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("无法获取家目录: %v", err)
	}

	privateKeyPath := homeDir + "/.ssh/id_rsa"

	privateKey, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatalf("无法读取私钥文件: %v", err)
	}

	// 解析私钥
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		log.Fatalf("无法解析私钥: %v", err)
	}

	config := ssh.ClientConfig{
		Timeout:         time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}

	for _, v := range ipList {
		//dial 获取ssh client
		sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", v), &config)
		if err != nil {
			log.Fatal("创建ssh client 失败", err)
		}

		//创建ssh-session
		session, err := sshClient.NewSession()
		if err != nil {
			log.Fatal("创建ssh client session", err)
		}

		defer session.Close()

		//执行远程命令
		cmdInfo, err := session.CombinedOutput(cmd)
		if err != nil {
			log.Fatal("远程命令执行失败", err)
			os.Exit(-1)
		}
		if len(cmdInfo) != 0 {
			log.Println(string(cmdInfo))
		}
	}
}
