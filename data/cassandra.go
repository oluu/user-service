package data

import (
	"log"
	"sync"

	"github.com/gocql/gocql"
	"github.com/oluu/user-service/util"
)

const (
	envCassandraKeyspace    = "CASSANDRA_KEYSPACE"
	envCassandraUsername    = "CASSANDRA_USERNAME"
	envCassandraPassword    = "CASSANDRA_PASSWORD"
	envCassandraServiceHost = "CASSANDRA_SERVICE_HOST"
	envCassandraServicePort = "CASSANDRA_SERVICE_PORT"
)

type cassandraRepository struct {
	session *gocql.Session
}

func (r *cassandraRepository) execute(name, q string, args ...interface{}) error {
	log.Println("INFO: executed", name)
	err := r.session.Query(q).Bind(args...).Exec()
	if err != nil {
		log.Println("ERR:", err)
	}
	return err
}

func (r *cassandraRepository) query(name, query string, pageSize int, args ...interface{}) (*gocql.Iter, error) {
	q := r.session.Query(query)
	if pageSize > 0 {
		q.PageState(nil).PageSize(pageSize)
	}
	iter := q.Bind(args...).Iter()
	count := iter.NumRows()
	log.Println("INFO: queries", name, "rows:", count)
	return iter, nil
}

// CassandraSessionManager ...
type CassandraSessionManager interface {
	getSession() (*gocql.Session, error)
}

type cassandraSessionManager struct {
	config  *gocql.ClusterConfig
	session *gocql.Session
	mu      sync.Mutex
}

func (csm *cassandraSessionManager) getSession() (*gocql.Session, error) {
	csm.mu.Lock()
	defer csm.mu.Unlock()
	if csm.session == nil {
		session, err := csm.config.CreateSession()
		if err != nil {
			log.Println("ERR:", err)
			return nil, err
		}
		csm.session = session
	}
	return csm.session, nil
}

func connectToCassandra() (CassandraSessionManager, error) {
	host := util.GetEnvString(envCassandraServiceHost)
	port := util.GetEnvInt(envCassandraServicePort)
	keyspace := util.GetEnvString(envCassandraKeyspace)
	username := util.GetEnvString(envCassandraUsername)
	password := util.GetEnvString(envCassandraPassword)

	cluster := gocql.NewCluster(host)
	cluster.Port = port
	cluster.Keyspace = keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: username, Password: password}
	cluster.Consistency = gocql.Quorum

	csm := &cassandraSessionManager{config: cluster}
	log.Printf("INFO: Establishing connection to cassandra (hosts: %v, port: %v)\n", cluster.Hosts, cluster.Port)
	_, err := csm.getSession()
	if err != nil {
		log.Printf("ERROR: Failed to establish a connection to cassandra (hosts: %v, port: %v)\n", cluster.Hosts, cluster.Port)
		return nil, err
	}
	return csm, nil
}
