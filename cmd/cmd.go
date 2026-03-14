package cmd

import (
	log "github.com/sirupsen/logrus"

	_ "ad2nx/core/imports"

	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use: "ad2nx",
}

func Run() {
	err := command.Execute()
	if err != nil {
		log.WithField("err", err).Error("Execute command failed")
	}
}
