// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2015 Canonical Ltd
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

package tests

import (
	"fmt"
	"regexp"

	"launchpad.net/snappy/_integration-tests/testutils/cli"
	"launchpad.net/snappy/_integration-tests/testutils/common"
	"launchpad.net/snappy/_integration-tests/testutils/wait"

	"gopkg.in/check.v1"
)

var _ = check.Suite(&serviceSuite{})

type serviceSuite struct {
	common.SnappySuite
}

func (s *serviceSuite) TearDownTest(c *check.C) {
	if !common.NeedsReboot() && common.CheckRebootMark("") {
		common.RemoveSnap(c, "hello-dbus-fwk.canonical")
	}
	// run cleanup last
	s.SnappySuite.TearDownTest(c)
}

func isServiceRunning(c *check.C) bool {
	packageVersion := common.GetCurrentVersion(c, "hello-dbus-fwk")
	service := fmt.Sprintf("hello-dbus-fwk_srv_%s.service", packageVersion)

	err := wait.ForActiveService(c, service)
	c.Assert(err, check.IsNil)

	statusOutput := cli.ExecCommand(
		c, "systemctl", "status",
		service)

	expected := "(?ms)" +
		".* hello-dbus-fwk_srv_.*\\.service .*\n" +
		".*Loaded: loaded .*\n" +
		".*Active: active \\(running\\) .*\n" +
		".*"

	matched, err := regexp.MatchString(expected, statusOutput)
	c.Assert(err, check.IsNil)
	return matched
}

func (s *serviceSuite) TestInstalledServiceMustBeStarted(c *check.C) {
	common.InstallSnap(c, "hello-dbus-fwk.canonical")

	c.Assert(isServiceRunning(c), check.Equals, true)
}

func (s *serviceSuite) TestServiceMustBeStartedAfterReboot(c *check.C) {
	if common.BeforeReboot() {
		common.InstallSnap(c, "hello-dbus-fwk.canonical")
		common.Reboot(c)
	} else if common.AfterReboot(c) {
		common.RemoveRebootMark(c)
		c.Assert(isServiceRunning(c), check.Equals, true)
	}
}

func (s *serviceSuite) TestServiceMustBeStartedAfterUpdate(c *check.C) {
	if common.BeforeReboot() {
		common.InstallSnap(c, "hello-dbus-fwk.canonical")
		common.CallFakeUpdate(c)
		common.Reboot(c)
	} else if common.AfterReboot(c) {
		common.RemoveRebootMark(c)
		c.Assert(isServiceRunning(c), check.Equals, true)
	}
}