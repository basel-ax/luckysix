# LuckySix

LuckySix is a comprehensive Go application designed to systematically generate, process, and monitor a vast number of blockchain wallets derived from BIP39 mnemonic phrases. The project generates every possible two-word combination, builds upon them to create 12-word seeds.

## Features

- Combinatorial generation of BIP39 word combinations (LuckyTwo, LuckyFive, LuckySix).
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

## Database Schema

- **luckytwos** table:
  - `id`: Primary key
  - `word_one`: Index of the first word (0-2047)
  - `word_two`: Index of the second word (0-2047)
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
2. Set `DATABASE_URL`.
3. Run the command as shown above.

## License

[Add license information here]