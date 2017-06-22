/*-
 * Copyright © 2015,2017 Jörg Pernfuß <code+github@paranoidbsd.net>
 * All rights reserved.
 *
 * Use of this source code is governed by a 2-clause BSD license
 * that can be found in the LICENSE file.
 */

package main // import "github.com/mjolnir42/legacy/cmd/metricsocket-client"

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	var sockPath string
	flag.StringVar(&sockPath, `socket`, ``, `Full patch to the metrics socket`)
	flag.Parse()
	if sockPath == `` {
		fmt.Fprintln(os.Stderr, `Error: -socket argument empty`)
		os.Exit(1)
	}
	b := make([]byte, 65536)
	empty := `{"metrics":[]}`

	conn, err := net.Dial("unixpacket", sockPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", empty)
		os.Exit(0)
	}
	defer conn.Close()
	i, err := conn.Read(b)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s\n", empty)
		os.Exit(0)
	}

	fmt.Fprintf(os.Stdout, "%s\n", string(b[:i]))
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
