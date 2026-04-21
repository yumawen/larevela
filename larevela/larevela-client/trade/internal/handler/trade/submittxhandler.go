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

// 前端提交链上交易ID
func SubmitTxHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SubmitTxReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := trade.NewSubmitTxLogic(r.Context(), svcCtx)
		resp, err := l.SubmitTx(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
