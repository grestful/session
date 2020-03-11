package session

import (
	"errors"
	"fmt"
	"github.com/grestful/utils"
)

type ISession interface {
	Close() bool
	Destroy(sid string) bool
	Gc(maxLeftTime int64) bool
	Open(savePath string) bool
	Read(sid string) map[string]interface{}
	Write(sid string, data map[string]interface{}) bool
	Error(sid string) error
	SetPrefix(keyPrefix string)
}

type IUserSession interface {
	SetData(data *utils.MapReader) error
	GetData() (data *utils.MapReader, err error)
	GetUserId() int64
	GetProperty(name string) interface{}
	GetAuthName() string
	SetSessionHandler(session ISession) error
	SetUserIdKey(idName string)
}

type UserSession struct {
	UserId   int64            `json:"userId"`
	Property *utils.MapReader `json:"property"`
	Sid      string           `json:"sid"`
	flag     uint8            //0 未读 1 已读 2 已写
	idName   string
	handler  ISession
}

func GetNewUserSession(sid string, session ISession) IUserSession {
	return &UserSession{
		UserId:   0,
		Property: utils.NewMapperReader(make(map[string]interface{})),
		Sid:      sid,
		handler:  session,
	}
}

func (s *UserSession) SetUserIdKey(idName string) {
	s.idName = idName
}

func (s *UserSession) getUserIdKey() string {
	if s.idName == "" {
		return "userId"
	}

	return s.idName
}

func (s *UserSession) SetData(data *utils.MapReader) error {
	s.UserId = data.ReadWithInt64(s.getUserIdKey(), 0)
	if s.UserId > 0 {
		s.Property = data
		s.flag = 2
		s.handler.Write(s.Sid, s.Property.GetValue())
	}
	return errors.New(fmt.Sprintf("no found id in data %v", data.GetValue()))
}

func (s *UserSession) GetData() (data *utils.MapReader, err error) {
	s.checkRead()
	return s.Property, nil
}

func (s *UserSession) GetUserId() int64 {
	s.checkRead()
	return s.UserId
}

func (s *UserSession) GetProperty(name string) interface{} {
	s.checkRead()
	return s.Property.ReadWithInterface(name, nil)
}

func (s *UserSession) GetAuthName() string {
	return "cookie"
}

func (s *UserSession) SetSessionHandler(session ISession) error {
	s.handler = session
	return nil
}

func (s *UserSession) checkRead() bool {
	if s.flag != 1 {
		if s.Sid == "" {
			return false
		}
		s.flag = 1
		s.Property = utils.NewMapperReader(s.handler.Read(s.Sid))
	}
	return true
}
