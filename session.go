package session

import (
	"strconv"
)

type ISession interface {
	Close () bool
	Destroy(sid string)  bool
	Gc(maxLeftTime int64)  bool
	Open(savePath string)  bool
	Read(sid string) map[string]string
	Write(sid string, data map[string]string)  bool
	Error(sid string) error
}

type IUserSession interface {
	SetData(data map[string]string) error
	GetData() (data map[string]string, err error)
	GetUserId() int64
	GetProperty(name string) interface{}
	GetAuthName() string
	SetSessionHandler(session ISession) error
}

type UserSession struct {
	UserId		int64	`json:"user_id"`
	Property    map[string]string `json:"property"`
	Sid         string  `json:"sid"`
	flag        uint8   //0 未读 1 已读 2 已写
	handler     ISession
}

func GetNewUserSession(sid string, session ISession) IUserSession {
	return &UserSession{
		UserId:   0,
		Property: make(map[string]string),
		Sid:      sid,
		handler:  session,
	}
}

func (s *UserSession) SetData(data map[string]string) error {
	if id,ok := data["user_id"]; ok {
		//s.UserId = base.String2Int64(id, 0)
		s.UserId,_ = strconv.ParseInt(id, 10, 64)
	}
	s.Property = data
	s.flag = 2
	s.handler.Write(s.Sid, s.Property)
	return nil
}

func (s *UserSession) GetData() (data map[string]string, err error) {
	s.checkRead()
	return s.Property,nil
}

func (s *UserSession) GetUserId() int64 {
	s.checkRead()
	return s.UserId
}

func (s *UserSession) GetProperty(name string) interface{} {
	s.checkRead()
	if v,ok := s.Property[name]; ok {
		return v
	}
	return nil
}

func (s *UserSession) GetAuthName() string {
	return "cookie"
}

func (s *UserSession)  SetSessionHandler(session ISession) error {
	return nil
}

func (s *UserSession) checkRead() bool {
	if s.flag != 1 {
		if s.Sid == "" {
			return false
		}
		s.flag = 1
		s.Property = s.handler.Read(s.Sid)
	}
	return true
}