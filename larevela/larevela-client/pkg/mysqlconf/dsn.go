package mysqlconf

import (
	"fmt"
	"net"
	"strings"
)

// DSN returns a MySQL driver DSN for go-zero sqlx.
func (m Conf) DSN() string {
	if strings.TrimSpace(m.DataSource) != "" {
		return m.DataSource
	}
	if !m.hasStructuredFields() {
		return ""
	}
	user := strings.TrimSpace(m.User)
	host := strings.TrimSpace(m.Host)
	if host == "" {
		host = "localhost"
	}
	port := m.Port
	if port == 0 {
		port = 3306
	}
	db := strings.TrimSpace(m.DBname)
	addr := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		user, m.Password, addr, db)
}

func (m Conf) hasStructuredFields() bool {
	return strings.TrimSpace(m.User) != "" && strings.TrimSpace(m.DBname) != ""
}
