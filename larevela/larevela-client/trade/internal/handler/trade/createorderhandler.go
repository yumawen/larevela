// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package trade

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"trade/internal/logic/trade"
	"trade/internal/svc"
	"trade/internal/types"
)

// 创建业务订单
func CreateOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := trade.NewCreateOrderLogic(r.Context(), svcCtx)
		resp, err := l.CreateOrder(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
