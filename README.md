# go-metar

A command-line tool for fetching METAR aviation weather reports.

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap mdaguerre/tap
brew install go-metar
```

### Go

```bash
go install github.com/mdaguerre/go-metar@latest
```

### From source

```bash
go build -o go-metar
```

## Usage

```bash
# Get decoded METAR for an airport
go-metar KJFK

# Get METARs for multiple airports
go-metar KJFK KLAX EGLL

# Get raw METAR string only
go-metar EGLL --raw

# Get both raw and decoded output
go-metar LFPG --all

# Multiple airports with flags
go-metar KJFK KLAX --raw

# Include TAF forecast
go-metar KJFK --taf

# Raw METAR and TAF
go-metar KJFK --raw --taf
```

## Options

| Flag | Short | Description |
|------|-------|-------------|
| `--raw` | `-r` | Show raw METAR string only |
| `--all` | `-a` | Show both raw and decoded output |
| `--taf` | `-t` | Include TAF forecast |

## Example Output

```
╭──────────────────────────────────────────────────╮
│ KJFK · John F Kennedy International              │
│ Time       02 Jan 2025 14:51 UTC                 │
│ Flight     VFR                                   │
│ Wind       350° at 8 kt                          │
│ Visibility 10+ SM                                │
│ Temp       7°C (Dewpoint: -1°C)                  │
│ Altimeter  30.21 inHg / 1023 hPa                 │
│ Clouds     Few @ 4500 ft, Scattered @ 25000 ft   │
╰──────────────────────────────────────────────────╯
```

## Testing

```bash
# Run unit tests
go test ./... -short

# Run all tests (includes integration tests that hit the API)
go test ./...
```

## Data Source

Weather data is fetched from [Aviation Weather Center](https://aviationweather.gov/).
