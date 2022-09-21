package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/csaldiasdev/distws/internal/agent"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const defaultPortValue = 0

var (
	repoRpcPort = flag.Uint("repoRpcPort", defaultPortValue, "a valid port")
	hubRpcPort  = flag.Uint("hubRpcPort", defaultPortValue, "a valid port")
	raftPort    = flag.Uint("raftPort", defaultPortValue, "a valid port")
	httpPort    = flag.Uint("httpPort", defaultPortValue, "a valid port")
	serfPort    = flag.Uint("serfPort", defaultPortValue, "a valid port")
	serfMember  = flag.String("member", "", "a serf address")
	prettylog   = flag.Bool("prettylog", false, "a boolean")
)

func main() {
	if *prettylog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	if *repoRpcPort == defaultPortValue {
		log.Error().
			Uint("configValue", *repoRpcPort).
			Msg("Invalid o missing RepoRpcPort value")

		os.Exit(1)
	}

	if *hubRpcPort == defaultPortValue {
		log.Error().
			Uint("configValue", *hubRpcPort).
			Msg("Invalid o missing HubRpcPort value")

		os.Exit(1)
	}

	if *raftPort == defaultPortValue {
		log.Error().
			Uint("configValue", *raftPort).
			Msg("Invalid o missing RaftPort value")

		os.Exit(1)
	}

	if *httpPort == defaultPortValue {
		log.Error().
			Uint("configValue", *httpPort).
			Msg("Invalid o missing HttpPort value")

		os.Exit(1)
	}

	if *serfPort == defaultPortValue {
		log.Error().
			Uint("configValue", *serfPort).
			Msg("Invalid o missing SerfPort value")

		os.Exit(1)
	}

	var sMembers []string = nil

	if *serfMember != "" {
		sMembers = []string{*serfMember}
	}

	config := agent.AgentConfiguration{
		HttpPort:          *httpPort,
		RepositoryRpcPort: *repoRpcPort,
		HubRpcPort:        *hubRpcPort,
		RaftPort:          *raftPort,
		SerfPort:          *serfPort,
		SerfMembers:       sMembers,
	}

	a, err := agent.NewAgent(config)

	if err != nil {
		log.Error().
			Err(err).
			Msg("Error creating agent")
		os.Exit(1)
	}

	go a.Run()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Info().Msg("Program exited")
}

func init() {
	// Define customized usage output
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, "Retrieve your external IP.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n    %s [flags]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Flags:")
		fmt.Fprintf(os.Stderr, "  -h help\n    \tshow this usage message\n")
		flag.PrintDefaults()
	}

	// Parse CLI Flags
	flag.Parse()
}
