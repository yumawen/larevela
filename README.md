Chain ID ：103 #solana devnet
  sol : System Program (11111111111111111111111111111111)  
  usdc/usdt :  USDC Mint (devnet default): 4zMMC9srt5Ri5X14GAgXhaHii3GnPAEERYPJgZJDncDU
               Token Program（SPL）：TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA
               Associated Token Program：ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL  #USDC ATA exists first.

build front
npm create vite@latest larevela-frontend -- --template react
npm install
t
build design 
(Currently only Solana is supported, with extensibility reserved for future needs. The chain and payment need to be modified) #暂时只支持solana，预留了可拓展性，需修改chain以及ledger

Frontend -> trade-api -> order-rpc -> payment-rpc -> chain-rpc -> Wallet/RPC -> payment-rpc -> ledger-rpc + order-rpc -> Frontend
下单、生成支付意图、链上支付、链上确认、记账、订单完成

1.trade-api.api
receiving frontend requests
auth and validation
orchestrating RPC calls
returning the payment params and status the frontend actually needs

2.order-rpc
creating orders
storing product, user, and price information
managing order states such as unpaid, paid, or closed
linking orders with payment records

3.payment-rpc
creating payment intents
generating payment parameters
recording the submitted txHash
validating whether a payment meets the required conditions
managing the payment state machine such as created, confirming, and paid

4.chain-rpc
etching transaction details
fetching receipts
parsing  transfer logs or contract events
calculating confirmations
hiding differences across chains and RPC providers

5.ledger-rpc
recording fund flows after successful payment
distinguishing payment, refund, and adjustment accounting actions
serving as the basis for audit and reconciliation
ensuring there is a clear record between “order succeeded” and “funds were booked”

build client (http://localhost:8888)

goctl api new trade        #api
goctl rpc new order       
goctl rpc new payment
goctl rpc new chain
goctl rpc new ledger


  database build
  docker start mysql-server
  /home/ma/larevela/larevela/larevela-client/model/sql

orders (create)  #orderNo = "order-" + time.Now().UnixNano()
Purpose: create the business master order first; orderNo becomes the primary correlation key across the flow.

payment_intents (initial insert)
Purpose: persist a pending payment intent with amount/receiver/unsigned-tx context—defines what should be paid.

idempotency_records (create_intent)
Purpose: protect create-intent from duplicates (double-click/retries) and avoid duplicate side effects.

payment_transactions
Purpose: record submitted on-chain transaction facts (txId, etc.) as broadcast evidence.

payment_intents update (submitted)
Purpose: bind the intent to a real tx_id and advance status to submitted for confirmation flow.

idempotency_records (submit_tx)
Purpose: idempotent guard + processing snapshot for submit-tx (processing -> success/failed) for safe retry/recovery.

chain_scan
Purpose: persist chain scan progress (block/slot/tx cursor) for resume/reconciliation jobs.

payment_intents confirmation update
Purpose: write chain parse results back into payment state machine (confirming/paid/failed) as payment truth.

ledger_entries (success only)
Purpose: post accounting ledger entries for audit-grade financial traceability.

orders update to paid (success only)
Purpose: finalize business order state and trigger delivery/entitlement, completing business closure.

  run a local node
  
  npm run dev 

#(order 8081, payment 8082, chain 8083, ledger 8084; trade HTTP 8888)

  go run order.go -f etc/order.yaml
  go run payment.go -f etc/payment.yaml
  go run chain.go -f etc/chain.yaml
  go run ledger.go -f etc/ledger.yaml
  go run trade.go -f etc/trade-api.yaml # last boot
  