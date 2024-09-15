# gopingstorm
gopingstorm is a multi-host ICMP monitoring tool. It sends/collects a bunch of packets and related datapoints.

It uses some kind of terminal based UI.

# Build and Test
go test -v

# Running gopingstorm
It'll look something like this... obviously I have no idea what I am doing.

go run gopingstorm.go CIDR [options...]

Options include:
  -a arp_timeout
  -r randomize_hosts
  -w icmp_wait

