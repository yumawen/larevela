// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package trade

import (
	"net/http"
	"strings"

	"trade/internal/logic/trade"
	"trade/internal/svc"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询支付状态
func QueryPaymentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryPaymentReq
		path := strings.Trim(r.URL.Path, "/")
		parts := strings.Split(path, "/")
		if len(parts) > 0 {
			req.PaymentNo = strings.TrimSpace(parts[len(parts)-1])
		}

		l := trade.NewQueryPaymentLogic(r.Context(), svcCtx)
		resp, err := l.QueryPayment(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
