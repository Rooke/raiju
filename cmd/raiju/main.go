package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/nyonson/raiju"
	"github.com/peterbourgon/ff/v3/ffcli"
)

// version is set by build tools during linking
var version = "undefined"

func main() {
	cmdLog := log.New(os.Stderr, "raiju: ", 0)

	rootFlagSet := flag.NewFlagSet("raiju", flag.ExitOnError)
	verbose := rootFlagSet.Bool("v", false, "increase log verbosity")

	btc2satCmd := &ffcli.Command{
		Name:       "btc2sat",
		ShortUsage: "raiju btc2sat <btc>",
		ShortHelp:  "Convert bitcoins to satoshis",
		Exec: func(_ context.Context, args []string) error {
			if len(args) != 1 {
				return errors.New("btc2sat only takes one arg")
			}

			if *verbose {
				cmdLog.Printf("converting %s btc to sats", args[0])
			}

			btc, err := strconv.ParseFloat(args[0], 64)
			if err != nil {
				return fmt.Errorf("unable to parse arg: %s", args[0])
			}

			raiju.PrintBtc2sat(btc)
			return nil
		},
	}

	spanCmd := &ffcli.Command{
		Name:       "span",
		ShortUsage: "raiju span <pubkey>",
		ShortHelp:  "Span network graph from node",
		Exec: func(_ context.Context, args []string) error {
			if len(args) != 1 {
				return errors.New("span only takes one arg")
			}

			host := "localhost"
			// take empty for client default
			tlsPath := "/home/lightning/.lnd/tls.cert"
			macDir := "/home/lightning/.lnd/data/chain/bitcoin/mainnet"

			basicClient, err := lndclient.NewBasicClient(host, tlsPath, macDir, "mainnet")

			if err != nil {
				return err
			}

			getInfo := lnrpc.GetInfoRequest{}
			info, err := basicClient.GetInfo(context.Background(), &getInfo)

			if err != nil {
				return err
			}

			cmdLog.Printf("connected to %s", info.GetAlias())

			return nil
		},
	}

	versionCmd := &ffcli.Command{
		Name:       "version",
		ShortUsage: "raiju version",
		ShortHelp:  "Version of raiju",
		Exec: func(_ context.Context, args []string) error {
			if len(args) != 0 {
				return errors.New("version does not take any args")
			}

			fmt.Fprintln(os.Stdout, version)
			return nil
		},
	}

	root := &ffcli.Command{
		ShortUsage:  "raiju [flags] <subcommand>",
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{btc2satCmd, spanCmd, versionCmd},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		// no need to output redundant message
		if err != flag.ErrHelp {
			cmdLog.Fatalln(err)
		} else {
			os.Exit(1)
		}
	}
}
