package driver

import (
	"time"

	"gopkg.in/mgo.v2"
)

type MongoStore struct {
	Host      []string
	User      string
	Pwd       string
	PoolLimit int
	Timeout   time.Duration
	Session   *mgo.Session
}

func NewMongoStore(host []string, username, password string, timeout int) *MongoStore {
	m := MongoStore{
		Host:      host,
		User:      username,
		Pwd:       password,
		PoolLimit: 2048,
		Timeout:   time.Duration(timeout) * time.Second,
	}
	return &m
}

func (m *MongoStore) GetSession(direct, fastFail bool) error {
	dialInfo := mgo.DialInfo{}
	dialInfo.Addrs = m.Host
	dialInfo.Direct = direct
	dialInfo.Username = m.User
	dialInfo.Password = m.Pwd
	dialInfo.PoolLimit = m.PoolLimit
	dialInfo.Timeout = m.Timeout
	dialInfo.Source = "admin"
	dialInfo.FailFast = fastFail
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		return err
	}
	// defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	m.Session = session
	return nil
}

func (m *MongoStore) DBStore(dbname string, direct, fastFail bool) (*mgo.Database, error) {
	if m.Session == nil {
		err := m.GetSession(direct, fastFail)
		if err != nil {
			return nil, err
		}
	}
	return m.Session.DB(dbname), nil
}
