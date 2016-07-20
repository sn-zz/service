// sn - https://github.com/sn
package main

import (
    "fmt"
    "time"
)

type Session struct {
    Id       uuid      `json:"sessionId"`
    UserId   uuid      `json:"userId"`
    Expires  time.Time `json:"expires"`
}

type Sessions []Session

const SessionTime = 86400 // one day

var sessions Sessions

func CreateSession(userId uuid) (Session, error) {
    if id := GenerateUuid(); id != "" {
        s := Session{Id:id,UserId:userId,Expires:time.Now().Add(SessionTime)}
        sessions = append(sessions, s)
        return s, nil
    }
    return Session{}, fmt.Errorf("Could not generate UUID")
}

func GetSession(id uuid) Session {
    for _, s := range sessions {
        if s.Id == id {
            return s
        }
    }
    return Session{}
}

func FindSession(hash string) Session {
    for _, s := range sessions {
        if GenerateSha1Hash(string(s.Id)) == hash {
            return s
        }
    }
    return Session{}
}

func UpdateSessionTime(id uuid) error {
    for i, s := range sessions {
        if s.Id == id {
            sessions[i].Expires = time.Now().Add(SessionTime)
        }
    }
    return fmt.Errorf("Could not find session")
}

func CleanSessions() {
    for i, s := range sessions {
        if time.Now().After(s.Expires) {
            sessions = append(sessions[:i], sessions[i+1:]...)
        }
    }
}
