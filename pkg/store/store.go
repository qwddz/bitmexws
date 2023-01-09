package store

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"math/rand"
)

//pool of connections to master and slave
type connectionsPool struct {
	master *sqlx.DB
	slave  []*sqlx.DB
}

type Store struct {
	config     Config
	connection connectionsPool
}

//New connect to mysql db
func New(config Config) (*Store, error) {
	store := Store{
		config: config,
		connection: connectionsPool{
			master: nil,
			slave:  nil,
		},
	}

	var sconnections []*sqlx.DB

	for _, slave := range config.Host.Slave {
		slc, err := store.openConnection(slave)
		if err != nil {
			continue
		}

		if err := slc.Ping(); err != nil {
			continue
		}

		sconnections = append(sconnections, slc)
	}

	mc, err := store.openConnection(config.Host.Master)
	if err != nil {
		return nil, err
	}

	store.connection.master = mc
	store.connection.slave = sconnections

	return &store, nil
}

//Close opened connections
func (s *Store) Close() error {
	if err := s.connection.master.Close(); err != nil {
		return err
	}

	for _, slave := range s.connection.slave {
		if err := slave.Close(); err != nil {
			return err
		}
	}

	return nil
}

// SlaveConnection return opened slave connection or master connection if slave pool is empty
func (s *Store) SlaveConnection() *sqlx.DB {
	ls := len(s.connection.slave)

	if ls == 0 {
		return s.ForceMasterConnection()
	}

	rI := rand.Intn(ls)

	return s.connection.slave[rI].Unsafe()
}

// ForceMasterConnection return opened master connection
func (s *Store) ForceMasterConnection() *sqlx.DB {
	return s.connection.master.Unsafe()
}

//connect to host
func (s *Store) openConnection(host string) (*sqlx.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", s.config.User, s.config.Password, host, s.config.Name)

	c, err := sqlx.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	return c, nil
}
