# NETGEAR Check Plugin

The NETGEAR Monitoring Plugin is a lightweight Go-based command-line tool designed to collect and report hardware and network metrics from NETGEAR devices via their API.
It provides output using go-check, allowing easy integration into Icinga 2 checks.
> [!WARNING]
>
> This plugin is intended **only for use with NETGEAR AV Line devices**. Use with other devices is not supported and may produce incorrect results.

### Features

- Fetch and report:
  - CPU usage  
  - RAM usage  
  - Fan speed  
  - Temperature sensors  
  - Port statistics (inbound and outbound)
  - PoE statistics (enabled state and current power)
- Icinga-compatible check results with perfdata  
- Configurable output via command-line flags  

## Documentation

### Installation

Ensure you have the recent Go version installed.
You can build and install check-netgear as follows:
```bash
go install github.com/icinga/check-netgear@latest
```

## Arguments

| Argument           | Description                                                          |
|--------------------|----------------------------------------------------------------------|
| `-u`, `--username` | **Required**. Username for API login                                 |
| `-p`, `--password` | **Required**. Password for API login                                 |
| `-H`, `--hostname` | **Optional**. Device hostname or IP (default: http://192.168.112.19) |
| `--mode`           | **Optional**. Modes to display: basic, ports, poe (default: basic)   |
| `--port`           | **Optional**. List of port numbers to check (default: 1–8)           |
| `--nocpu`          | **Optional**. Hide CPU info                                          |
| `--noram`          | **Optional**. Hide RAM info                                          |
| `--notemp`         | **Optional**. Hide temperature info                                  |
| `--nofans`         | **Optional**. Hide fans info                                         |
| `-h`, `--help`     | **Optional**. Show help message                                      |


## Example
```bash
check_netgear -u admin -p VerySecurePassword --mode basic
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

## Support

For questions, suggestions, or issues, please reach out through the Icinga community channels or open an issue in the project repository.

## Thanks

Special thanks to the NETWAYS team for the go-check library.

## Contributing

Contributions are welcome!
You can help by:
- Submitting bug fixes
- Testing on different NETGEAR devices
- Improving documentation or examples

## License

NETGEAR Icinga Check Plugin and its documentation are licensed under the terms of the [GNU General Public License v3.0](LICENSE) (GPL-3.0).


