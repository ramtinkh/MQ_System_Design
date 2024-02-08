package models

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/websocket"
)

type Subscriber struct {
	UserMap  map[string]*websocket.Conn
	UserList []string
}

func (s *Subscriber) Subscribe(conn *websocket.Conn, key string) error {
	if _, ok := s.UserMap[key]; ok {
		return ErrSubscriberAlreadyExists
	}

	s.UserList = append(s.UserList, key)
	s.UserMap[key] = conn
	return nil
}

func (s *Subscriber) Unsubscribe(key string) error {
	if _, ok := s.UserMap[key]; !ok {
		return ErrSubscriberDoesNotExist
	}

	delete(s.UserMap, key)
	for i, k := range s.UserList {
		if k == key {
			s.UserList = append(s.UserList[:i], s.UserList[i+1:]...)
			break
		}
	}
	return nil
}

func (s *Subscriber) SendData(data Data) error {
	conn := s.UserMap[s.UserList[rand.Intn(len(s.UserList))]]
	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return conn.WriteMessage(websocket.TextMessage, body)
}

func NewSubscriber() *Subscriber {
	return &Subscriber{
		UserMap: make(map[string]*websocket.Conn),
	}
}
