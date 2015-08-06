package models

import (
	"fmt"
	"math/rand"
	"time"
)

type State struct {
	TotalRounds int
	Players     []string
	Tables      []string
	Rounds      []*Round
}

func NewState(totalRounds int, players, tables []string) *State {
	firstRound := NewRound(players, tables, nil)
	return &State{
		TotalRounds: totalRounds,
		Players:     players,
		Tables:      tables,
		Rounds:      []*Round{firstRound},
	}
}

func (s *State) CurrentRound() *Round {
	return s.Rounds[len(s.Rounds)-1]
}

func (s *State) CurrentRoundNumber() int {
	return len(s.Rounds)
}

// Move to the next round. Fails if some scores aren't
// received yet.
func (s *State) AdvanceRound() error {
	if !s.CurrentRound().Finished() {
		return fmt.Errorf("Current round not finished.")
	}

	newRound := NewRound(s.Players, s.Tables, s.Rounds)

	s.Rounds = append(s.Rounds, newRound)

	return nil
}

func (s *State) Scores(player string) []int {
	scores := make([]int, 0)
	for _, round := range s.Rounds {
		// Search seatings
		found := false
		for _, seating := range round.Seatings {
			if seating.Finished && (seating.North == player || seating.South == player) {
				found = true
				scores = append(scores, seating.NSScore)
			} else if seating.Finished && (seating.West == player || seating.East == player) {
				found = true
				scores = append(scores, seating.WEScore)
			}
		}
		if !found {
			break
		}
	}

	return scores
}

// ApplyScore applies the score to the table with the given players in the
// current round. If that table isn't found, an error is returned.
func (s *State) ApplyScore(north, south, west, east string, nsScore, weScore int) error {
	seating := s.CurrentRound().FindSeating(north, south, west, east)
	if seating == nil {
		return fmt.Errorf("Couldn't find a seating with these players: N %s, S %s, W %s, E %s", north, south, west, east)
	}

	seating.Finished = true
	seating.NSScore = nsScore
	seating.WEScore = weScore
	if nsScore >= weScore {
		seating.NSWins = true
	}

	return nil
}

type Round struct {
	Seatings []*Seating
}

func NewRound(players, tables []string, pastRounds []*Round) *Round {
	if pastRounds == nil || len(pastRounds) == 0 {
		return &Round{
			Seatings: newRandomSeatings(players, tables),
		}
	}

	// There have been past rounds, use the latest round to find the next
	return &Round{
		// TODO: switch to a losers move seating
		// Seatings: newLosersMoveSeating(players, tables, lastRound),
		Seatings: newRandomSeatings(players, tables),
	}
}

// Finished is true when all the Seatings are finished.
func (r *Round) Finished() bool {
	finished := true
	for _, seating := range r.Seatings {
		finished = finished && seating.Finished
	}
	return finished
}

func (r *Round) FindSeating(north, south, west, east string) *Seating {
	for _, seating := range r.Seatings {
		if seating.North == north && seating.South == south &&
			seating.West == west && seating.East == east {
			return seating
		}
	}
	return nil
}

type Seating struct {
	TableName string
	North     string
	South     string
	West      string
	East      string

	Finished bool
	NSScore  int
	WEScore  int
	NSWins   bool
}

func newRandomSeatings(players, tables []string) []*Seating {
	t := time.Now()
	rand.Seed(int64(t.Nanosecond())) // no shuffling without this line

	for i := len(players) - 1; i > 0; i-- {
		j := rand.Intn(i)
		players[i], players[j] = players[j], players[i]
	}

	var seatings []*Seating
	for i := 0; i < len(tables); i++ {
		j := i * 4

		newSeating := &Seating{
			TableName: tables[i],
			North:     players[j],
			South:     players[j+1],
			West:      players[j+2],
			East:      players[j+3],
		}

		seatings = append(seatings, newSeating)
	}

	return seatings
}
