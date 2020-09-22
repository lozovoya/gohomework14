package app

import (
	"encoding/json"
	"github.com/lozovoya/gohomework14_1/cmd/bank/app/dto"
	"github.com/lozovoya/gohomework14_1/pkg/card"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	cardSvc *card.Service
	mux     *http.ServeMux
}

func NewServer(cardSvc *card.Service, mux *http.ServeMux) *Server {
	return &Server{cardSvc: cardSvc, mux: mux}
}

func (s *Server) Init() {
	s.mux.HandleFunc("/getAllCards", s.getAllCards)
	s.mux.HandleFunc("/getHolderCards", s.getHolderCards)
	s.mux.HandleFunc("/addHolderCard", s.addHolderCard)
}

func (s *Server) getAllCards(w http.ResponseWriter, r *http.Request) {
	cards := s.cardSvc.AllCards()
	if len(cards) == 0 {
		log.Println("no cards available")
		err := s.SendReply(w, cards, "no cards available")
		if err != nil {
			log.Println(err)
		}
		return
	}

	err := s.SendReply(w, cards, "")
	if err != nil {
		log.Println(err)
		return
	}
}

func (s *Server) getHolderCards(w http.ResponseWriter, r *http.Request) {

	holderid, err := strconv.Atoi(r.FormValue("holderid"))
	if err != nil {
		log.Println(err)
		return
	}
	cards, err := s.cardSvc.HolderCards(holderid)
	if err != nil {
		s.SendReply(w, nil, err.Error())
		return
	}

	err = s.SendReply(w, cards, "")
	if err != nil {
		log.Println(err)
		return
	}

}

func (s *Server) addHolderCard(w http.ResponseWriter, r *http.Request) {

	holderid, err := strconv.Atoi(r.FormValue("holderid"))
	if err != nil {
		log.Println(err)
		return
	}
	issuer := r.FormValue("issuer")
	image := r.FormValue("image")

	err = s.cardSvc.AddHolderCard(issuer, holderid, image)
	if err != nil {
		s.SendReply(w, nil, err.Error())
		return
	}
	s.SendReply(w, nil, "Card is added")
	return
}

func (s *Server) SendReply(w http.ResponseWriter, cards []*card.Card, message string) (err error) {

	var respBody []byte

	if len(cards) != 0 {
		dtos := make([]*dto.CardDTO, len(cards))
		for i, c := range cards {
			dtos[i] = &dto.CardDTO{
				Id:       c.Id,
				Number:   c.Number,
				Issuer:   c.Issuer,
				HolderId: c.HolderId,
				Type:     c.Type,
			}
		}

		respBody, err = json.Marshal(dtos)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		var dtos = &dto.MessageDTO{
			Message: message,
		}
		respBody, err = json.Marshal(dtos)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(respBody)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
