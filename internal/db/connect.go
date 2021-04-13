package db

import (
	"chatter/internal/logger"

	"github.com/gocql/gocql"
)

var (
	Session *gocql.Session
)

func Init() {
	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "chatter"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Connected to cassandra db")
	Session = session
}

func GetSession() *gocql.Session {
	return Session
}
