Design
===========

# main: #
1. get configuration
example configuration file:
ex-okex.json
<exchanger=okex>
    <server=>
</exchanger>
2. open websocket client to the server
    2.1 exURL
    2.2 userSeceret
3. subscribe to the ticker
    3.1 tickPair
    3.2 tickHandle
        3.2.1 init
        initialize the backend db
        3.2.2 validate
        correct form
        correct time
        3.2.3 store
        db tables:
        1 insert this receiver into the index table if not already exist.
        keys include: exchanger name/id(a exchanger table?), ticker pairs
        data: working table id
        2 create the working table for this receiver if the working table is not alread exist.
        keys include: 
        3.2.4 tables
        A: Trading
        1) Exchangers
        2) Traded Pairs
        3) Trade Book (market depth)
        4) Ticker
        6) Transactions (trades)
        #it Transactions are the fundemental. Should the database stores only the transactions and computes the trade book/ticker on demand? currently prepared in the table

        B: Account
        1) user  

# APIs #
1. query support (APIs)

# Related link #

## OKex ##
API Documents: https://www.okex.com/docs/en/
SDK: https://github.com/okcoin-okex/open-api-v3-sdk

## Binance ##
API Documents: https://github.com/binance-exchange/binance-official-api-docs
*By testing with curl and ping, the server api.binance.com has been blocked

## Huobi ##
API Documents: https://github.com/huobiapi/API_Docs

## Poloniex ##
API Documents: https://poloniex.com/support/api/