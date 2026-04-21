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

// 创建支付意图
func CreatePaymentIntentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreatePaymentIntentReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := trade.NewCreatePaymentIntentLogic(r.Context(), svcCtx)
		resp, err := l.CreatePaymentIntent(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
