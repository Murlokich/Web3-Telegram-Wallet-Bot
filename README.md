# Web3-Telegram-Wallet-Bot
Web3 Telegram Wallet Bot is a project that services as a custodial hot crypto wallet for ETH through telegram Bot.

## Functionality
* Create new or import existing ETH wallet
* Execute the transactions
* Provide comfortable UI 

## Technology stack
* Telegram API: [telebot](https://github.com/tucnak/telebot)
* PostgreSQL: [pgx](https://github.com/jackc/pgx)
* Ethereum API: [go-ethereum](https://github.com/ethereum/go-ethereum)

## Node
Due to the lack of free disk space both full (1TB disk space) and pruned (100GB disk space) nodes were not considered as an option,
so for this project infura API was used:
* Remote Node using [infura](https://www.infura.io/) - 0GB disk free space

As an alternative could have used [alchemy](https://www.alchemy.com/)

## Roadmap
1) Telegram API integration:
    * Bot responding with "hello" to slash command
2) Main functionality for ETH
   * Account creation/migration
   * Account management
   * Asset transfer
3) Visualisation of the account, transfers, statistics, coin prices

## Functionality:
* /register - creates a new wallet with mnemonic
* /migrate \<mnemonic\> - allows to import existing wallet via mnemonic
* /new_address - creates new address for existing account