// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/snapcore/snapd/tests/lib/fakestore/refresh"
	"github.com/snapcore/snapd/tests/lib/fakestore/store"
)

var (
	start           = flag.Bool("start", false, "Start the store service")
	topDir          = flag.String("dir", "", "Directory to be used by the store to keep and serve snaps, <dir>/asserts is used for assertions")
	makeRefreshable = flag.String("make-refreshable", "", "List of snaps with new versions separated by commas")
	addr            = flag.String("addr", "locahost:11028", "Store address")
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	if *start {
		return runServer(*topDir, *addr)
	}

	if *makeRefreshable != "" {
		return runManage(*topDir, *makeRefreshable)
	}

	return fmt.Errorf("please specify either start or make-refreshable")
}

func runServer(topDir, addr string) error {
	st := store.NewStore(topDir, addr)

	if err := st.Start(); err != nil {
		return err
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	return st.Stop()
}

func runManage(topDir, snaps string) error {
	// setup fake new revisions of snaps for refresh
	snapList := strings.Split(snaps, ",")
	return refresh.MakeFakeRefreshForSnaps(snapList, topDir)
}
