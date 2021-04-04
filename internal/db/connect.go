package db

import (
	"fmt"

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
		panic(err)
	}
	fmt.Println("Connected to Cassandra. All systems go!")
	Session = session
}

func GetSession() *gocql.Session {
	return Session
}
