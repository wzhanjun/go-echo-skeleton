package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "gen",
	Short: "go-echo-skeleton 代码生成工具",
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}
