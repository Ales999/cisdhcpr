# Cisco DHCP Report

This script is designed to generate a report of DHCP from cisco devices.
It will parse the configuration files and extract information about DHCP configurations.

The report will include details such as:

- DHCP pool name
- Network address
- Subnet mask
- Start and stop IP addresses
- Gateway address
- Excluded IP addresses
- Reserved IP addresses
- And all options like 'option 42 ....'

Exampe run:

```bash
go run . ./example_host1.txt ./example_host2.txt
```

Example report with analyze example files `example_host1.txt` and `example_host2.txt`:

```text
--> Open file: example_host1.txt
-----
Scope:    lan_host1
Vlan:     Vlan1
Mask Bit: 26
Start IP: 192.168.1.1
Stop  IP: 192.168.1.62
Gateway:  192.168.1.62
	 option 42 ip 192.168.1.6
Excludes:
	 192.168.1.61
	 192.168.1.8
	 192.168.1.62
Reserved IP:
	 Host: policom_h1_de31 IP: 192.168.1.10 MAC: xxxx.xxxx.xxde.31 GW: 192.168.1.1
		  option 42 ip 192.168.1.1
		  option 4 ip 192.168.1.1
		  option 120 hex xxxx.xxxx.120.120.xxxx.xxx.xxxx
		  option 43 hex xxxx.example-data-for-43.xxxx
	 Host: policom_h1_fb4f IP: 192.168.1.12 MAC: xxxx.xxxxx.xxfb.4f GW: 192.168.1.1
	 Host: printer_fef5 IP: 192.168.1.14 MAC: xxxx.xxxx.xxxx.xx GW: 192.168.1.1
-----
Scope:    lan2_host1
Vlan:     Vlan2
Mask Bit: 24
Start IP: 192.168.100.1
Stop  IP: 192.168.100.254
Gateway:  192.168.100.1
Excludes:
	 192.168.100.1
Reserved IP:
	 Host: policom_h8_d132 IP: 192.168.100.10 MAC: xxxx.xxxx.xxd1.32 GW: 192.168.100.1
		  option 4 ip 192.168.100.1
-----------------------------
--> Open file: example_host2.txt
-----
Scope:    lan_host2
Vlan:     Vlan1
Mask Bit: 26
Start IP: 192.168.5.1
Stop  IP: 192.168.5.62
Gateway:  192.168.5.1
	 option 42 ip 192.168.5.1
Excludes:
	 192.168.5.1
	 192.168.5.3
	 192.168.5.4
	 192.168.5.5
	 192.168.5.6
	 192.168.5.12
Reserved IP:
	 Host: policom_h2_de3f IP: 192.168.5.10 MAC: xxx.xxx.xxde.3f GW: 192.168.5.1
		  option 42 ip 192.168.5.1
		  option 4 ip 192.168.5.1
		  option 120 hex xxxx.xxxx.xxx.xxx.xxxx.xxx.xxxx
	 Host: policom_h2_fb4f IP: 192.168.5.15 MAC: xxx.xxxx.xxxx.xx GW: 192.168.5.1
	 Host: printer_fef5 IP: 192.168.5.3 MAC: xxxx.xxxx.xxxx.xx GW: 192.168.5.1
-----------------------------
```
