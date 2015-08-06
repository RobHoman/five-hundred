package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-martini/martini"

	"FiveHundred/models"
)

type Server struct {
	martini *martini.Martini
	port    int
	state   *models.State
}

func NewServer(port int, state *models.State) *Server {
	mart := martini.New()
	router := martini.NewRouter()

	mart.Action(router.Handle)

	server := &Server{
		martini: mart,
		port:    port,
		state:   state,
	}

	router.Get("/state", server.HTTPGetState)
	router.Post("/state", server.HTTPPostScore)

	return server
}

func (server *Server) Serve() error {
	err := http.ListenAndServe(fmt.Sprintf(":%d", server.port), server.martini)
	if err != nil {
		return err
	}
	return nil
}

type StateResponse struct {
	Round       int
	TotalRounds int
	Scores      []*PlayerScore
	Rounds      []*models.Round
}

type PlayerScore struct {
	Player string
	Scores []int
}

func (server *Server) HTTPGetState(w http.ResponseWriter, req *http.Request) {
	server.respondWithState(w, req)
}

type PostScoreReq struct {
	North   string
	South   string
	West    string
	East    string
	NSScore int
	WEScore int
}

func (server *Server) HTTPPostScore(w http.ResponseWriter, req *http.Request) {
	var psReq PostScoreReq
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		respondWithErr(w, err)
		return
	}

	err = json.Unmarshal(body, &psReq)
	if err != nil {
		respondWithErr(w, err)
		return
	}

	state := server.state
	err = state.ApplyScore(psReq.North, psReq.South, psReq.West, psReq.East, psReq.NSScore, psReq.WEScore)
	if err != nil {
		respondWithErr(w, err)
		return
	}

	if state.CurrentRound().Finished() {
		state.AdvanceRound()
	}

	server.respondWithState(w, req)
}

func (server *Server) respondWithState(w http.ResponseWriter, req *http.Request) {
	state := server.state
	resp := &StateResponse{
		Round:       state.CurrentRoundNumber(),
		TotalRounds: state.TotalRounds,
		Scores:      make([]*PlayerScore, 0),
		Rounds:      state.Rounds,
	}

	// Get all the scores
	for _, player := range state.Players {
		scores := state.Scores(player)
		resp.Scores = append(resp.Scores, &PlayerScore{
			Player: player,
			Scores: scores,
		})
	}

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		respondWithErr(w, err)
		return
	}
}

func respondWithErr(w http.ResponseWriter, err error) {
	w.Write([]byte(fmt.Sprintf("%#v", err)))
}
