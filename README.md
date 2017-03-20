# Mabul
Mabul is a tool used to test DDoS vectors within controlled envrioments that you have both permission and control over
and may only be used in that capacity.

#Go Doc

[![GoDoc](https://godoc.org/github.com/levigross/mabul?status.svg)](https://godoc.org/github.com/levigross/mabul)

# License Exception

You may only use Mabul in a lawful way. Any use of Mabul in an unlawful (in the country where you are,
the country where the packets originate from, and the country where the packets are sent to) manner will automatically
revoke both the software license as well as the license to the source code.

# How To Use
```
$ bin/mabul
Mabul is a program designed as a test suite for DDoS mitigation programs

Usage:
  mabul [command]

Available Commands:
  h2          Conducts connection attacks using H2
  http        This is designed to execute layer 7 attacks
  ip          This will execute IP level attacks (both stateless and stateful)
  tcp         All attacks using TCP as a vector (both stateful and stateless)
  tls         A brief description of your command
  udp         Launches UDP style attacks (stateless)

Flags:
      --config string     config file (default is $HOME/.mabul.yaml)
      --logLevel string   The level of logging you wish to have (default "info")
  -t, --toggle            Help message for toggle

Use "mabul [command] --help" for more information about a command.
```