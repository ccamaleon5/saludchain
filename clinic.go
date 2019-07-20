package main

import (
	"github.com/ccamaleon5/saludchain/state"
//	cm "github.com/ccamaleon5/saludchain/cmd"
	"os"

	"github.com/tendermint/tendermint/abci/server"
	"github.com/tendermint/tendermint/abci/types"
	mgo "gopkg.in/mgo.v2"

	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
    "github.com/spf13/cobra"
)

func main() {
	var cmdInit = &cobra.Command{
		Use:   "init",
		Short: "Init Saludchain",
		Long: `Init a new Blockchain on Tendermint`,
		Args: cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
		  initStore()
		},
	  }
	
	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdInit)
	rootCmd.Execute()
}

func initStore() error {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))

	// Create the application
	var app types.Application

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	db := session.DB("tendermintdb")

	// Clean the DB on each reboot
	collections := [3]string{"doctors", "patients", "medicalappointments"}

	for _, collection := range collections {
		db.C(collection).RemoveAll(nil)
	}

	app = state.NewJSONStateApplication(db)

	// Start the listener
	srv, err := server.NewServer("tcp://0.0.0.0:36658", "socket", app)
	if err != nil {
		return err
	}
	srv.SetLogger(logger.With("module", "abci-server"))
	if err := srv.Start(); err != nil {
		return err
	}

	// Wait forever
	cmn.TrapSignal(func() {
		srv.Stop()
	})
	return nil
}