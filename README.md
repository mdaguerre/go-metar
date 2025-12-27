# go-metar

A command-line tool for fetching METAR aviation weather reports.

## Installation

```bash
go install github.com/mdaguerre/go-metar@latest
```

Or build from source:

```bash
go build -o go-metar
```

## Usage

```bash
# Get decoded METAR for an airport
go-metar KJFK

# Get raw METAR string only
go-metar EGLL --raw

# Get both raw and decoded output
go-metar LFPG --all
```

## Options

| Flag | Short | Description |
|------|-------|-------------|
| `--raw` | `-r` | Show raw METAR string only |
| `--all` | `-a` | Show both raw and decoded output |

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

## Data Source

Weather data is fetched from [Aviation Weather Center](https://aviationweather.gov/).
