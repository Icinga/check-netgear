# NETGEAR Monitoring Plugin for Icinga

## Table of Contents
- [About](#about)  
- [License](#license)  
- [Documentation](#documentation)  
- [Support](#support)  
- [Requirements](#requirements)  
- [Thanks](#thanks)  
- [Contributing](#contributing)
---

## About
The NETGEAR Monitoring Plugin is a lightweight Go-based command-line tool designed to collect and report hardware and network metrics from NETGEAR devices via their API.
It provides output using go-check, allowing easy integration into Icinga 2 checks.

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
---

## License
This plugin is released under the MIT License.  
See the included LICENSE file for full details.

This plugin uses the following third-party components:

- go-check — MIT License  
- spf13/pflag — BSD License  

---

## Documentation

### Installation
1. Ensure you have Go installed (1.20 or newer).  
2. Clone this repository and build the plugin:
   ```bash
   git clone https://git.icinga.com/obarbashyn/netgear-icinga-plugin.git
   cd netgear-icinga-plugin
3. The binary to deploy is located in bin/check-netgear


## Required Flags
- `-u`, `--username` — Username for API login
- `-p`, `--password` — Password for API login

## Optional Flags
- `-H`, `--hostname` — Device hostname or IP (default: http://192.168.112.19)
- `--mode` — Modes to display: basic, ports, poe (default: basic)
- `--port` — List of port numbers to check (default: 1–8)
- `--nocpu` — Hide CPU info
- `--noram` — Hide RAM info
- `--notemp` — Hide temperature info
- `--nofans` — Hide fans info
- `-h`, `--help` — Show help message

## Example
```bash
check_netgear -u admin -p VerySecurePassword --mode basic
```

## Output Example
```yaml
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

## Requirements
Go 1.20 or higher
Icinga 2

## Thanks
Special thanks to the NETWAYS team for the go-check library.

# Contributing
Contributions are welcome!
You can help by:
- Submitting bug fixes
- Testing on different NETGEAR devices
- Improving documentation or examples
