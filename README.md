# NETGEAR monitoring plugin

A lightweight Go-based CLI tool to monitor hardware metrics and port statistics from network devices via an API.

## Features

- Fetches and reports:
  - CPU usage
  - RAM usage
  - Fan speed
  - Temperature sensors
  - Port statistics (inbound, TODO outbound)
- Icinga-compatible output using [`go-check`](https://github.com/NETWAYS/go-check)

## Usage

```bash
go run main.go [flags]
```

### Required Flags

- `-u`, `--username` Username for API login  
- `-p`, `--password` Password for API login

### Optional Flags

- `-H`, `--hostname`    Device hostname or IP (default: `http://192.168.112.19`)
- `--mode`  Modes to display: `basic`, `ports` (default: `basic`)
- `--port`  List of port numbers to check (default: 1–8)
- `--nocpu` Hide CPU info
- `--noram` Hide RAM info
- `--notemp` Hide temperature info
- `--nofans` Hide fans info
- `-h`, `--help` Show help message

## Example

```bash
go run main.go -u admin -p VerySecurePassword --mode basic --mode ports --port 1 --port 2 --port 3
```

## Output

Icinga-style status with perfdata, e.g.:

```
[WARNING] Device Info: Uptime - 1 days, 0 hrs, 31 mins, 29 secs
\_ [OK] CPU Usage: 7.13%
\_ [OK] RAM Usage: 32.46%
\_ [OK] Temperature
    \_ [OK] sensor-System1: 44.0°C
    \_ [OK] sensor-MAC: 47.0°C
    \_ [OK] sensor-System2: 45.0°C
\_ [WARNING] Fans
    \_ [WARNING] FAN-1: 0 RPM
|CPU=7.13;;;0;100 RAM=32.46;;;0;100 sensor-System1=44;;;0 sensor-MAC=47;;;0 sensor-System2=45;;;0 'Fans speed'=0;;;0
```

## Dependencies

- [go-check](https://github.com/NETWAYS/go-check)
- [spf13/pflag](https://github.com/spf13/pflag)

## License

MIT License
