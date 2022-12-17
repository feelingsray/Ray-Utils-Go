package mongo

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

func NewMongoStore(host []string, username string, password string) *MongoStore {
  m := MongoStore{
    Host:      host,
    User:      username,
    Pwd:       password,
    PoolLimit: 4096,
    Timeout:   20 * time.Second,
    Session:   nil,
  }
  return &m
}

func (m *MongoStore) GetSession() error {
  dialInfo := mgo.DialInfo{}
  dialInfo.Addrs = m.Host
  dialInfo.Direct = false
  dialInfo.Username = m.User
  dialInfo.Password = m.Pwd
  dialInfo.PoolLimit = m.PoolLimit
  dialInfo.Timeout = m.Timeout
  dialInfo.Source = "admin"
  session, err := mgo.DialWithInfo(&dialInfo)
  if err != nil {
    return err
  }
  //defer session.Close()
  session.SetMode(mgo.Monotonic, true)
  m.Session = session
  return nil
}

func (m *MongoStore) DBStore(dbname string) (*mgo.Database, error) {
  if m.Session == nil {
    err := m.GetSession()
    if err != nil {
      return nil, err
    }
  }
  return m.Session.DB(dbname), nil
  
}
