package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"FiveHundred/models"
	"FiveHundred/server"
)

const TOTAL_ROUNDS = 3

func main() {
	players, err := readPlayers()
	if err != nil {
		fmt.Printf("Failed to read the players")
		return
	}
	fmt.Printf("Players: %#v\n", players)

	tables := []string{"T1", "T2", "T3"}
	fmt.Printf("Tables: %#v\n", tables)

	state := models.NewState(TOTAL_ROUNDS, players, tables)

	srvr := server.NewServer(5000, state)
	err = srvr.Serve()
	if err != nil {
		fmt.Printf("Failed to start server because %#v\n", err)
		return
	}
}

func readPlayers() ([]string, error) {
	file, err := os.Open("players.txt")
	if err != nil {
		return nil, err
	}

	players := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		players = append(players, strings.TrimSpace(scanner.Text()))
	}

	return players, nil
}
