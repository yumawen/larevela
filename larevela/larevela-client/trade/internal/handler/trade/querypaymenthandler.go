// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package trade

import (
	"net/http"

	"trade/internal/logic/trade"
	"trade/internal/svc"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// 查询支付状态
func QueryPaymentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := trade.NewQueryPaymentLogic(r.Context(), svcCtx)
		resp, err := l.QueryPayment()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
