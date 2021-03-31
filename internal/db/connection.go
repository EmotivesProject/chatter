package db

import (
	"fmt"

	"github.com/gocql/gocql"
)

func Init() *gocql.Session {
	cluster := gocql.NewCluster("cassandra")
	cluster.Keyspace = "chatter"
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Cassandra. All systems go!")
	return session
}
