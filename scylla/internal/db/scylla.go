package db

import (
	"github.com/gocql/gocql"
	"log"
	"results/succ"
	"time"
)

var Session *gocql.Session

func InitScylla(hosts []string, keyspace string) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 5 * time.Second
	cluster.ConnectTimeout = 10 * time.Second
	cluster.ProtoVersion = 4
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())

	var err error
	Session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%s\n", succ.ScyllaConnected)
}
