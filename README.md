# netprobe

Reads in an IP or CIDR range and outputs any associated hostnames. 

## Installation

Assuming GO is installed...

```bash
go install github.com/cybercdh/netprobe@latest
```

## Usage

```bash
echo 1.2.3.4 | netprobe
```
or
```bash
echo 1.2.3.4/24 | netprobe -dns 8.8.8.8 -port 53 -v
```

## Options
```bash
Usage of netprobe:
  -c int
    	Set the concurrency level (default 20)
  -dns string
    	Custom DNS resolver address (ip only) (default "8.8.8.8")
  -port string
    	DNS server port (default "53")
  -v	See IP and Hostname as output
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)