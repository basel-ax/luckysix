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

3. Set up your PostgreSQL database and set the `DATABASE_URL` environment variable:
   ```bash
   export DATABASE_URL="postgres://user:password@localhost/dbname?sslmode=disable"
   ```

## Usage

### Generate LuckyTwo Combinations

To generate all possible two-word combinations from the BIP39 word list:

```bash
go run main.go luckytwo generate
```

This command will:
- Connect to the PostgreSQL database specified by `DATABASE_URL`.
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

The application uses the following environment variables:

- `DATABASE_URL`: PostgreSQL connection string (required)

## Development

To run the application in development:

1. Ensure PostgreSQL is running.
2. Set `DATABASE_URL`.
3. Run the command as shown above.

## License

[Add license information here]