# gosuiteql

A command-line tool for executing SuiteQL queries against NetSuite's REST API.

## Features

- Execute SuiteQL queries via command line
- Support for query input via:
  - Command line argument
  - File input
  - Standard input (pipe)
- Colored JSON output for better readability
- Pagination support via limit/offset parameters

## Installation

```bash
go install github.com/ajbonner/gosuiteql@latest
```

## Usage

```bash
# Execute query directly
gosuiteql -query "SELECT * FROM transaction"

# Execute query from file
gosuiteql -file query.sql

# Pipe query via stdin
echo "SELECT * FROM transaction" | gosuiteql

# With pagination
gosuiteql -query "SELECT * FROM transaction" -limit 10 -offset 0
```

## Configuration

Create a `.env` file with your NetSuite credentials:

```env
NETSUITE_ACCOUNT_ID=your_account_id
NETSUITE_CONSUMER_KEY=your_consumer_key
NETSUITE_CONSUMER_SECRET=your_consumer_secret
NETSUITE_TOKEN=your_token
NETSUITE_TOKEN_SECRET=your_token_secret
```

## Disclaimer

This software is provided "as is", without warranty of any kind, express or implied. Use at your own risk.

## License

MIT License - see [LICENSE](LICENSE) for details. 
