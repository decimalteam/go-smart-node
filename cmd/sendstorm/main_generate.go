package main

import (
	"fmt"
	"os"

	"bitbucket.org/decimalteam/go-smart-node/sdk/wallet"
	"github.com/spf13/cobra"
)

func cmdGenerate() *cobra.Command {
	var mnemonicsCount int

	var cmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate mnemonics and save them to file",
		Run: func(cmd *cobra.Command, args []string) {
			err := cmd.Flags().Parse(args)
			if err != nil {
				fmt.Println(err)
				return
			}
			mnemonicsFile, err := cmd.Flags().GetString(mnemonicsFlag)
			if err != nil {
				fmt.Println(err)
				return
			}
			f, err := os.Create(mnemonicsFile)
			if err != nil {
				fmt.Printf("create mnemonic file error: %s\n", err.Error())
				return
			}
			defer f.Close()
			fmt.Printf("create %d mnemonics\n", mnemonicsCount)
			for i := 0; i < mnemonicsCount; i++ {
				mn, err := wallet.NewMnemonic("")
				if err != nil {
					fmt.Printf("create mnemonic error: %s\n", err.Error())
					continue
				}
				_, err = fmt.Fprintf(f, "%s\n", mn.Words())
				if err != nil {
					fmt.Printf("write to mnemonic file error: %s\n", err.Error())
				}
			}
		},
	}

	cmd.PersistentFlags().IntVar(&mnemonicsCount, "count", 10, "count of mnemonics")

	return cmd
}
