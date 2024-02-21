package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-wordwrap"
	"golang.org/x/term"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(2)
	}

	// tolerant being set to true causes the program to quietly ignore any entries on STDIN that are not valid CIDRs
	tolerant := false

	var checkIP string

	for _, arg := range os.Args {
		if arg == "--help" || arg == "-h" {
			showUsage()
			os.Exit(2)
		}

		if arg == "--tolerant" {
			tolerant = true
		}

		if !strings.HasPrefix(arg, "-") {
			checkIP = arg
		}
	}

	if checkIP == "" {
		showUsage()
		os.Exit(2)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		checkCIDR := strings.Trim(line, "\n")

		cidrBaseIP, mask, err := splitMask(checkCIDR)
		if err != nil {
			if tolerant {
				continue
			} else {
				fatal(err)
			}
		}

		rawCIDRBaseIP, err := addressToInt(cidrBaseIP)
		if err != nil {
			if tolerant {
				continue
			} else {
				fatal(err)
			}
		}

		rawTestIP, err := addressToInt(checkIP)
		if err != nil {
			fatal(err)
		}

		ignoreLen := 32 - mask

		if rawCIDRBaseIP>>ignoreLen == rawTestIP>>ignoreLen {
			fmt.Println(checkCIDR)
			os.Exit(0)
		}
	}

	// end of input reached with no match found :(
	os.Exit(1)
}

func showUsage() {
	usageStr := "Usage: \033[1m" + os.Args[0] + "\033[0m \033[4mIP ADDRESS\033[0m\n\n" +
		"\t\033[4mIP ADDRESS\033[0m is an IPv4 address in the usual human-readable format, e.g. 192.168.0.1\n\n" +
		"\033[1m" + os.Args[0] + "\033[0m reads a list of CIDR-notation IP ranges, one per line, from STDIN and " +
		"checks to see if the IP address given by the first argument is within the range of any of them. When " +
		"a matching range is found, the matching CIDR-notation is printed back to STDOUT and the program exits " +
		"immediately. That is, it always outputs the first matching range in the input and then exits, ignoring any " +
		"additional matching ranges.\n\n" +
		"When a matching range was found, the program exits with code 0. When a match is not found, it exits with code " +
		"1. For any other errors, it exits with a code greater than 1."

	if term.IsTerminal(0) {
		width, _, err := term.GetSize(0)
		if err != nil {
			// if getting terminal size results in an error, just act as if we're not on a terminal
			fmt.Println(usageStr)
			return
		}
		fmt.Println(wordwrap.WrapString(usageStr, uint(width)))
	} else {
		fmt.Println(usageStr)
	}
}

func fatal(err error) {
	fmt.Println(err.Error())
	os.Exit(3)
}

func splitMask(cidr string) (ip string, mask int, err error) {
	parts := strings.Split(cidr, "/")
	if len(parts) != 2 {
		return "", 0, fmt.Errorf(cidr + " could not be parsed as a valid CIDR string")
	}
	mask, err = strconv.Atoi(parts[1])
	return parts[0], mask, err
}

func addressToInt(addressAsText string) (addressRaw uint32, err error) {
	parts := strings.Split(addressAsText, ".")

	part0, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf(addressAsText+" could not be parsed as an IP address: %w", err)
	}
	part1, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf(addressAsText+" could not be parsed as an IP address: %w", err)
	}
	part2, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, fmt.Errorf(addressAsText+" could not be parsed as an IP address: %w", err)
	}
	part3, err := strconv.Atoi(parts[3])
	if err != nil {
		return 0, fmt.Errorf(addressAsText+" could not be parsed as an IP address: %w", err)
	}

	return uint32(part0<<24) + uint32(part1<<16) + uint32(part2<<8) + uint32(part3), nil
}
