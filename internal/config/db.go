package config

type Host struct {
	Master string
	Slave  []string
}

type DB struct {
	Host     Host
	Name     string
	User     string
	Password string
}
