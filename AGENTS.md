# Guidelines for AI Coding Agents

## Project Overview

LuckySix is a comprehensive Go application designed to systematically generate, process, and monitor a vast number of blockchain wallets derived from BIP39 mnemonic phrases. The project generates every possible two-word combination, builds upon them to create 12-word seeds.

Key features include:
- Combinatorial generation of BIP39 word combinations (LuckyTwo, LuckyFive, LuckySix).

The application uses PostgreSQL for data storage. Project uses GORM for database interactions and supports configuration via environment variables.
