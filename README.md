# grepcidr-reverse

**grepcidr-reverse** is the opposite of [**grepcidr**](https://www.pc-tools.net/unix/grepcidr/).

**grepcidr** can take a list of IP addresses and tell you which of them fall into one (or more) given IP ranges. By contrast, **grepcidr-reverse** takes a single IP address and a list of ranges and tells you which, if any, range the given address falls into.

## Usage

### Synopsis

```
grepcidr-reverse [--tolerant] IP_ADDRESS
```

- IP ADDRESS is an IPv4 address in the usual human-readable format, e.g. 192.168.0.1

### Description

**grepcidr-reverse** reads a list of CIDR-notation IP ranges, one per line, from STDIN and checks to see if the IP address given by the first argument is within the range of any of them. When a matching range is found, the matching CIDR-notation is printed back to STDOUT and the program exits immediately. That is, it always outputs the first matching range in the input and then exits, ignoring any additional matching ranges.

When a matching range was found, the program exits with code 0. When a match is not found, it exits with code 1. For any other errors, it exits with a code greater than 1.

### Options

- `--tolerant` Silently ignore any rows in the IP range input that are not valid IP ranges. May be useful if piping input from tools such as iproute2 that may contain words like "default" instead of a range.

## Examples

The contrived example: find which range an address falls into from a hard-coded list:

```shell
grepcidr-reverse 192.168.1.30 <<< "10.0.0.0/16
172.0.0.0/8
192.168.1.0/24
1.1.1.0/24"
```

Expected output: `192.168.1.0/24`

---

See which iproute2 route is used to reach the destination address 192.168.1.1

```shell
ip route | cut -f1 -d' ' | grepcidr-reverse --tolerant 192.168.1.1
```

Here we use the `--tolerant` flag to make grepcidr-reverse silently ignore any rows from `ip route` that don't start with a valid CIDR range, such as "default".

## Building

- Make sure you have the Go developer tools installed. See https://go.dev/doc/install.
- In the project directory, run `go get` to install dependecies.
- Run `make` to generate the `grepcidr-reverse` binary.

### Build the debian package

- Do the above
- Install dpkg-deb on your system
- Run `make debain`. **grepcidr-reverse.deb** is generated in the "debian" directory.