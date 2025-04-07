# 🌐 Web3 Telegram Wallet Bot

**Web3 Telegram Wallet Bot** is a lightweight, custodial hot wallet for Ethereum, accessible directly via Telegram. Users can seamlessly create, manage, and interact with their ETH wallet using simple Telegram commands.

---

## 🚀 Features

- 📲 **Create or Import Wallet** – Generate a new Ethereum wallet or import an existing one using a mnemonic phrase.
- 💸 **Send Transactions** – Transfer ETH with a simple slash command.
- 👁 **User-Friendly Interface** – Convenient interaction through Telegram UI.
- 🔁 **Address Management** – Generate and switch between multiple wallet addresses.
- 🧮 **Balance Tracking** – Instantly check wallet balance in ETH.

---

## 🧠 Commands

| Command               | Description                                                  |
|-----------------------|--------------------------------------------------------------|
| `/register`           | Creates a new ETH wallet and returns the mnemonic.           |
| `/migrate <mnemonic>` | Imports an existing wallet using your mnemonic phrase.       |
| `/new_address`        | Generates a new address for your wallet.                     |
| `/switch_address`     | Allows switching between your previously created addresses.  |
| `/get_balance`        | Returns the balance (in ETH) of the currently active address.|

---

## ⚙️ Technology Stack

- **Telegram Bot API** – [telebot](https://github.com/tucnak/telebot)
- **Database** – PostgreSQL via [pgx](https://github.com/jackc/pgx)
- **Ethereum Integration** – [go-ethereum](https://github.com/ethereum/go-ethereum)

---

## 🌍 Ethereum Node Setup

Due to limited disk space, running a local node (full/pruned) was not feasible. Instead, a remote node via Infura is used:

- ✅ [Infura](https://www.infura.io/) – No disk space required.
- 🔁 Alternative: [Alchemy](https://www.alchemy.com/)

---

## 🛠 Makefile Commands

The project includes a `Makefile` for easier setup and management:

- 📦 **First-time setup:**  
  ```bash
  make rebuild-up
  ```
- 🔄 **Subsequent runs:**  
  ```bash
  make up
  ```

- 📘 **View all available commands:**  
  ```bash
  make help
  ```

## 🗺 Roadmap

1. ✅ **Basic Telegram Bot Setup**
   - Responds with "hello" to test commands.

2. ✅ **Wallet Management**
   - Wallet creation (`/register`)
   - Wallet migration via mnemonic (`/migrate`)
   - New address generation (`/new_address`)
   - Address switching (`/switch_address`)
   - ETH balance lookup (`/get_balance`)

3. ⏳ **ETH Transfer Feature**
   - Implement ETH transfer command (e.g., `/send <address> <amount>`)

4. 🛠 **Enhanced UX**
   - Wallet visualization
   - Transaction history
   - Account statistics
   - Real-time ETH price display
   
5. 🧪 **Token Trading Support**
   - Enable trading of a custom cryptocurrency using a smart contract
