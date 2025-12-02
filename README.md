# LuckySix

LuckySix is a comprehensive Go application designed to systematically generate, process, and monitor a vast number of blockchain wallets derived from BIP39 mnemonic phrases. The project generates every possible two-word combination, builds upon them to create 12-word seeds.

## LuckySix
This is a part of the project. All results you can get by 3 projects.    
 - [LuckySix](https://github.com/basel-ax/luckysix)
 - [LuckyEth](https://github.com/basel-ax/lucky-eth)
 - [LuckyCosmos](https://github.com/basel-ax/lucky-cosmos)

## Features

- Combinatorial generation of BIP39 word combinations (LuckyTwo, LuckyFive, LuckySix).
- Generation of 12-word BIP39 wallet mnemonics from LuckySix combinations.
- Resumable generation from the last saved combination.
- PostgreSQL database storage using GORM.

## Prerequisites

- Go 1.25.3 or later
- PostgreSQL database

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/basel-ax/luckysix.git
   cd luckysix
   ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Copy the example environment file and configure your database connection:
    ```bash
    cp example.env .env
    # Edit .env with your database credentials
    ```

4. Set up your PostgreSQL database and configure the database connection in `.env` using either:
    - `DATABASE_URL` environment variable:
      ```bash
      export DATABASE_URL="postgres://user:password@localhost/dbname?sslmode=disable"
      ```
    - Or individual environment variables:
      ```bash
      export DB_HOST=localhost
      export DB_PORT=5432
      export DB_USER=user
      export DB_PASSWORD=password
      export DB_NAME=dbname
      export DB_SSLMODE=disable
      ```

## Usage

### Generate LuckyTwo Combinations

To generate all possible two-word combinations from the BIP39 word list:

```bash
go run main.go luckytwo generate
```

This command will:
- Connect to the PostgreSQL database using the configured environment variables.
- Automatically migrate the database schema.
- Generate and save all 2048 x 2048 = 4,194,304 possible two-word combinations.
- Resume from the last generated combination if the process was interrupted previously.
- Display progress logs for each completed word index.

The combinations are stored in the `luckytwos` table with `word_one` and `word_two` as indices into the BIP39 word list.

### Generate LuckyFive Combinations

To generate random five-word combinations from the BIP39 word list:

```bash
go run main.go luckyfive generate
```

To generate ALL possible five-word combinations (2048^5 combinations - this will take a very long time!):

```bash
go run main.go luckyfive generate --all
# or
go run main.go luckyfive generate -a
```

This command will:
- Generate 100 random, non-repeating five-word combinations per run (default behavior).
- When using `--all` flag: Generate ALL possible combinations with distinct words.
- Store each unique combination in the `lucky_fives` table.
- Skip duplicate combinations automatically.

**Note:** The `--all` flag will generate a massive number of combinations (approximately 34 trillion) and will take an extremely long time to complete. Use with caution!

### Generate LuckySix Combinations

To generate six-word combinations by combining LuckyFive and LuckyTwo:

```bash
go run main.go luckysix generate
```

This command will:
- Combine each LuckyFive (5 words) with available words from LuckyTwo entries.
- Ensure no word repetition in the final 6-word combination.
- Resume from the last generated combination if interrupted.
- Store combinations in the `lucky_sixes` table with references to source LuckyFive and LuckyTwo.

**Note:** You must generate LuckyTwo and LuckyFive data before running LuckySix generation.

### Generate Random LuckySix Combinations

To generate six-word combinations using random LuckyFive entries:

```bash
go run main.go luckysix generate-random
```

This command will:
- Select random LuckyFive entries from the database
- Check if each LuckyFive has already been used (duplicate checking via `lucky_five_id`)
- Combine them with random LuckyTwo entries to create LuckySix combinations
- Generate 1000 random combinations per run
- Use shuffling algorithms to ensure randomness
- Store combinations in the `lucky_sixes` table with references to source LuckyFive and LuckyTwo

**Note:** This is useful when you have a large amount of LuckyFive entries and want to generate LuckySix combinations randomly rather than sequentially.

## Database Schema

- **luckytwos** table:
  - `id`: Primary key
  - `word_one`: Index of the first word (0-2047)
  - `word_two`: Index of the second word (0-2047)
  - `created_at`, `updated_at`: Timestamps

- **lucky_fives** table:
  - `id`: Primary key
  - `word_one`, `word_two`, `word_three`, `word_four`, `word_five`: Indices (0-2047)
  - `created_at`, `updated_at`: Timestamps

- **lucky_sixes** table:
  - `id`: Primary key
  - `lucky_five_id`: Foreign key reference to lucky_fives
  - `lucky_two_id`: Foreign key reference to luckytwos
  - `word_one`, `word_two`, `word_three`, `word_four`, `word_five`, `word_six`: Indices (0-2047)
  - `created_at`, `updated_at`: Timestamps

- **wallet_balances** table:
  - `id`: Primary key
  - `lucky_six_id`: Foreign key reference to lucky_sixes
  - `mnemonic`: Complete 12-word BIP39 mnemonic phrase
  - `address`: Wallet address (optional)
  - `cosmos_address`: Cosmos wallet address (optional)
  - `balance`: Wallet balance (stored as string to handle large numbers)
  - `balance_updated_at`: Timestamp of last balance update
  - `is_notified`: Flag indicating if balance notifications have been sent
  - `created_at`, `updated_at`: Timestamps

## Configuration

The application uses the following environment variables for database connection:

- `DATABASE_URL`: PostgreSQL connection string (optional, if not set, individual vars are used)
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database user (required if DATABASE_URL not set)
- `DB_PASSWORD`: Database password (required if DATABASE_URL not set)
- `DB_NAME`: Database name (required if DATABASE_URL not set)
- `DB_SSLMODE`: SSL mode (default: disable)

## Development

To run the application in development:

1. Ensure PostgreSQL is running.
2. Configure database connection in `.env` file.
3. Run the commands as shown in the Usage section.

### Generate Wallet Mnemonics

To generate 12-word BIP39 wallet mnemonics from LuckySix combinations:

```bash
go run main.go wallet generate
```

This command will:
- Process LuckySix combinations and generate valid 12-word BIP39 mnemonics
- Use the first 6 words from each LuckySix combination as the base
- Generate the remaining 6 words to create a valid BIP39 mnemonic
- Store wallets in the `wallet_balances` table with the complete mnemonic
- Resume from the last processed LuckySix ID if interrupted
- Skip already processed LuckySix combinations automatically

### Generation Order

For best results, generate combinations in this order:
1. `luckytwo generate` - Generate all two-word combinations first
2. `luckyfive generate` - Generate five-word combinations (use `--all` for complete generation)
3. `luckysix generate` - Generate six-word combinations from LuckyFive + LuckyTwo (sequential)
4. `luckysix generate-random` - Generate six-word combinations from random LuckyFive entries
5. `wallet generate` - Generate 12-word wallet mnemonics from LuckySix combinations

**Note:** Using `luckyfive generate --all` will generate all possible combinations but will take a very long time to complete.
