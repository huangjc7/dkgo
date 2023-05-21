package cmd

import (
	"github.com/huangjc7/dkgo/dk"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dkgo",
		Short: "快速部署实验环境kubernetes集群",
		Run: func(cmd *cobra.Command, args []string) {
			//dk.Dk.StopCh = make(chan int, 1)
			dk.Dk.Run()
			//<-dk.Dk.StopCh
		},
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}

func init() {
	//rootCmd.AddCommand(rootCmd)
	rootCmd.Flags().StringVar(&dk.Dk.Master, "master", "", "master ip")
	rootCmd.Flags().StringVar(&dk.Dk.Node, "node", "", "node ip")
}
