# vip_patroni
The vip_patroni utility runs as a daemon, monitors the local state of Patroni via its API, and adds or removes a VIP on the local host as needed. In case of errors, it outputs them to stderr, and these errors and their descriptions are processed by the OS syslog and recorded in /var/log/syslog.  

vip_patroni has its own API for monitoring the status and viewing the current application configuration.  

vip_patroni can accept configuration through:
- Config file
- CLI arguments
- Environment variables



### 🛠️ CLI Arguments

| Argument                 | Description 
|--------------------------|----------------------------------------------------------------------------------------------
| `--help`                 | Show help.
| `--version`              | Show version of vip_patroni.
| `--config string`        | Location of the configuration file.
| `--ip string`            | Virtual IP address to configure.
| `--patroni-url string`   | URL of Patroni API (default: `http://localhost:8008`) 
| `--api-port string`      | Port of vip_patroni API (default: `8010`) 
| `--netmask int`          | Netmask used for the IP. `-1` assigns default IPv4 mask (default: `32`) 
| `--interface string`     | Network interface to configure on. 
| `--interval string`      | DCS scan interval in milliseconds (default: `1000`) 
| --patroni-timeout-millis | HTTP timeout for Patroni API requests (default: 500)
| `--log-level string`     | Set log level (`trace`, `debug`, `info`, `warning`, `error`, `fatal`) (default: `INFO`) 



### Example: Running vip_patroni

```bash 
vip_patroni --config cmd/vip_patroni_config.yml
```

### Example Configuration File
```ini
ip: 80.73.64.7  
netmask: 32  
interface: eth0  
interval: 1000  
retry-after: 250  
retry-num: 2  
patroni-url: http://10.203.97.100:8008  
patroni-timeout-millis: 700
api-port: 6010  
log-level: INFO
```
### Example API Requests  
- GET http://10.203.97.100:8010/config – shows current configuration
- GET http://10.203.97.100:8010/status – shows status info:
```json
{
  "desired": 1,  // IP should be present (Patroni says this node is primary)
  "state": 1     // IP is actually present on the host
}
```







