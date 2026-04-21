build front
npm create vite@latest larevela-frontend -- --template react
npm install

build design 
(Currently only Solana is supported, with extensibility reserved for future needs. The chain and ledger modules need to be modified) #暂时只支持solana，预留了可拓展性，需修改chain以及ledger

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

build client

goctl api new trade        #api
goctl rpc new order       
goctl rpc new payment
goctl rpc new chain
goctl rpc new ledger

goctl rpc protoc order.proto \
  -I . \
  --go_out=. \
  --go-grpc_out=. \
  --zrpc_out=.