# vip_patroni
## Plugin for telegraf
The vip_patroni utility runs as a daemon, monitors the local state of Patroni via its API, and adds or removes a VIP on the local host as needed. In case of errors, it outputs them to stderr, and these errors and their descriptions are processed by the OS syslog and recorded in /var/log/syslog.  

vip_patroni has its own API for monitoring the status and viewing the current application configuration.  

vip_patroni can accept configuration through:
* Config file
* CLI arguments
* Environment variables


### **ARGs LIST:**
**--help**                - Show help.  
**--version**             - Show version of vip_patroni.  
**--config string**       - Location of the configuration file.  
**--ip string**           - Virtual IP address to configure.  
**--patroni-url string**  - URL of patroni API (default "http://localhost:8008")  
**--api-port string**     - Port of vip_patroni API (default "8010")  
**--netmask int**         - The netmask used for the IP address. Defaults to -1 which assigns ipv4 default mask. (default 32)  
**--interface string**    - Network interface to configure on.  
**--interval string**     - DCS scan interval in milliseconds. (default "1000")  
**--log-level string**    - Set log level (trace,debug,info,warning,error,fatal). (default "INFO")  


### **Examples of running vip_patroni**  
-  ip_patroni --config cmd/vip_patroni_config.yml"
### **Examples of vip_patroni configuration file**  
```ini
  ip: 80.73.64.7  
  netmask: 32  
  interface: eth0  
  interval: 1000  
  retry-after: 250  
  retry-num: 2  
  patroni-url: http://10.203.97.100:8008  
  api-port: 6010  
```
### **Examples of API requests to vip_patroni**  
http://10.203.97.100:8010/config  
http://10.203.97.100:8010/status  
  * desired - 1 if the IP should be present (server in master role)
  * state - 1 if the IP is actually present on the host








