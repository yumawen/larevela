package mysqlconf

// Conf holds MySQL connection settings shared across microservices.
// If DataSource is non-empty, DSN() returns it as-is (full driver DSN).
// Otherwise User/Password/Host/Port/DBname are composed.
type Conf struct {
	DataSource string `json:"dataSource,optional"`
	User       string `json:"user,optional"`
	Password   string `json:"password,optional"`
	Host       string `json:"host,optional"`
	Port       int    `json:"port,optional"`
	DBname     string `json:"dbname,optional"`
}
