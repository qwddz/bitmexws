package store

type Host struct {
	Master string
	Slave  []string
}

type Config struct {
	Host     Host
	Name     string
	User     string
	Password string
}
