# luckysix

## Start docker
```bash
docker-compose up -d
# if u need to rebuild all docker containers
docker-compose up -d --force-recreate
```
Output u can see in Docker Dashboard

## Init, update and install
```bash
docker-compose up -d
go mod init github.com/basel-ax/luckysix
go get -u ./... && go mod tidy 
# just install
go mod tidy
```

## Start project 
Localy, with enviroment in docker containers.  
```bash
go run ./cmd/app
```

## Project Info
Application Based on [Go-clean-template](https://github.com/evrone/go-clean-template)  
Sheduller Based on [jasonlvhit/gocron](https://github.com/jasonlvhit/gocron)  
Add to project - go get github.com/jasonlvhit/gocron  

Telegram API [Go Telegram Bot API](https://go-telegram-bot-api.dev/)  
Add to project - go get -u github.com/go-telegram-bot-api/telegram-bot-api/v5 

## ORM
[GORM](https://gorm.io/), top [from list](https://github.com/d-tsuji/awesome-go-orms) 
# LuckySix Project

LuckySix is a comprehensive Go application designed to systematically generate, process, and monitor a vast number of blockchain wallets derived from BIP39 mnemonic phrases. The project generates every possible two-word combination, builds upon them to create 12-word seeds, derives wallets for multiple blockchains (EVM and Cosmos), checks their balances, and sends Telegram notifications when a wallet with a non-zero balance is found.

## Features

-   **Combinatorial Generation**: Systematically creates all 4+ million two-word combinations from the BIP39 word list.
-   **Randomized Wallet Creation**: Generates `LuckyFive` (10 words) and `LuckySix` (12 words) combinations.
-   **Multi-Chain Wallet Derivation**: Derives wallet addresses for both Ethereum (EVM) and Cosmos chains from a single 12-word mnemonic.
-   **Automated Balance Checking**: Periodically checks wallet balances across multiple blockchains using public RPC endpoints (Ethereum, Arbitrum, Base, BNB Chain, Cosmos Hub).
-   **Telegram Notifications**: Sends instant alerts to a configured Telegram chat when a wallet with a balance greater than zero is discovered.
-   **Task-Based Architecture**: Uses a robust task scheduler (`gocron`) to run background processes for wallet generation, balance checking, and notifications.
-   **RPC Control**: Exposes key tasks via an AMQP (RabbitMQ) RPC interface for manual triggering and control.

## Services Overview

The application is built on a service-oriented architecture:

-   **Bip39 Service**: The core generation engine. It's responsible for creating the foundational `LuckyTwo` entities and then building `LuckyFive` and `LuckySix` combinations on top of them.
-   **Blockchain Service**: Handles all blockchain interactions. It processes `LuckySix` combinations to derive wallet addresses, queries public nodes for balances, and triggers notifications.
-   **Tasks Service**: The application's scheduler. It orchestrates the execution of background jobs like balance checking and sending notifications using `gocron`. It also exposes one-off tasks to be triggered via RPC.
-   **Telegram Service**: A dedicated service for sending formatted messages to a Telegram chat via the Telegram Bot API.

## Getting Started

### Prerequisites

-   Go (version 1.21 or higher)
-   PostgreSQL
-   RabbitMQ
-   Docker and Docker Compose (recommended for local setup)

### Installation

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/your-username/luckysix.git
    cd luckysix
    ```

2.  **Install dependencies:**
    ```sh
    go mod tidy
    ```

### Configuration

The application is configured via `config/config.yml` and environment variables.

1.  **Copy the example configuration:**
    ```sh
    cp config/config.example.yml config/config.yml
    ```

2.  **Edit `config/config.yml`** to provide the following:
    -   `PG_URL`: Your PostgreSQL connection string.
    -   `RMQ_URL`: Your RabbitMQ connection string.
    -   `TG_BOT_API`: Your Telegram bot token.
    -   `TG_CHAT_ID`: The ID of the Telegram chat for notifications.
    -   **Blockchain RPC Endpoints**: Public RPC URLs for `eth_mainnet_rpc`, `arbitrum_rpc`, `base_rpc`, `bnb_rpc`, and a REST endpoint for `cosmos_rpc`.

### Database Setup

You must create the required tables in your PostgreSQL database before running the application. Ensure your schema includes the following tables:
-   `luckytwos`
--   `luckyfives`
-   `lucky_sixes`
-   `wallet_balances` (with a unique constraint on `lucky_six_id`)

### Running the Application

```sh
go run cmd/app/main.go
```
The server will start, and the scheduled tasks for balance checking and notifications will be active.

## Usage & Workflow

The main workflow involves a sequence of generation and processing steps. These are designed as long-running tasks that you can trigger via the AMQP RPC interface.

**Note:** You will need an AMQP client to send these RPC messages. The `handler` is the routing key, and the `request` body can be `null` for these trigger-style tasks.

### Step 1: Generate Base Data (LuckyTwo)

This is a one-time, foundational step that populates the `luckytwos` table with all ~4.2 million combinations. This will take a significant amount of time.

-   **RPC Handler**: `generateLuckyTwo`

### Step 2: Generate LuckyFive Combinations

This task creates a batch of 10,000 random `LuckyFive` combinations. You can run this multiple times to populate the `luckyfives` table.

-   **RPC Handler**: `generateLuckyFive`

### Step 3: Generate LuckySix Combinations

This task creates a batch of 100,000 random `LuckySix` combinations. Run this multiple times to generate a large pool of potential wallets.

-   **RPC Handler**: `generateLuckySix`

### Step 4: Process Wallets

This task takes a batch of 100 unprocessed `LuckySix` records, derives their 12-word mnemonics and addresses (EVM & Cosmos), and saves them to the `wallet_balances` table.

-   **RPC Handler**: `processWallets`

### Step 5: Automated Tasks (No Trigger Needed)

Once the application is running, the following tasks are executed automatically by the scheduler:

-   **Check Wallet Balances**: Runs every 10 minutes to query the balances of wallets that haven't been checked yet.
-   **Notify Wallets**: Runs every 15 minutes to send Telegram notifications for any wallets found with a positive balance.

## API Endpoints

### HTTP API

The project includes a RESTful HTTP API with Swagger documentation. Once the server is running, you can access the Swagger UI at:
`http://localhost:8080/swagger/index.html`

### AMQP RPC Endpoints

The following handlers are available for remote procedure calls via RabbitMQ:

-   `checkRabbit`: A simple health check endpoint.
-   `generateLuckyTwo`: Starts the one-time generation of all `LuckyTwo` combinations.
-   `generateLuckyFive`: Starts a background task to generate a batch of `LuckyFive` combinations.
-   `generateLuckySix`: Starts a background task to generate a batch of `LuckySix` combinations.
-   `processWallets`: Starts a background task to process a batch of `LuckySix` records into wallets.
