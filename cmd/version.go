package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version  = "TempVersion" //use ldflags replace
	codename = "ad2nx"
	intro    = "A board backend based on multi core"
)

var versionCommand = cobra.Command{
	Use:   "version",
	Short: "Print version info",
	Run: func(_ *cobra.Command, _ []string) {
		showVersion()
	},
}

func init() {
	command.AddCommand(&versionCommand)
}

func showVersion() {
	fmt.Println(`========= ad2nx Version:`)
	fmt.Printf("%s %s 1001 (%s) ==========\n", codename, version, intro)
	//fmt.Printf("Supported cores: %s\n", strings.Join(vCore.RegisteredCore(), ", "))
	// Warning
	//fmt.Println(Warn("This version need board version >= 1.7.0."))
	//fmt.Println(Warn("The version have many changed for config, please check your config file"))
}
