package trademodel

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type Model struct {
	conn sqlx.SqlConn
}

func NewModel(conn sqlx.SqlConn) *Model {
	return &Model{conn: conn}
}

func (m *Model) ensureReady() error {
	if m == nil || m.conn == nil {
		return fmt.Errorf("trademodel is not initialized")
	}
	return nil
}
