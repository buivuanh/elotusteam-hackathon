package main

import (
	"context"
	"fmt"
	"os"

	"github.com/buivuanh/elotusteam-hackathon/utils"
)

func main() {
	//argsWithProg := os.Args
	//dbConnectString := argsWithProg[1]
	dbConnectString := os.Getenv("DATABASE_URL")

	dbPool, err := utils.NewConnectionPool(context.Background(), dbConnectString)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
	defer dbPool.Close()

	if err = utils.Up(context.Background(), dbConnectString); err != nil {
		fmt.Fprintf(os.Stderr, "migrate up failed: %v\n", err)
		os.Exit(1)
	}

}
