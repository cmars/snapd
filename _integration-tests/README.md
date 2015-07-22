# Integration testing for snappy

## Requirements

 *  autopkgtest (>= 3.15.1)

    Get the latest autopkgtest deb from
    https://packages.debian.org/sid/all/autopkgtest/download

 *  Internet access in the test bed.

## Testing a virtual machine

You can execute the full integration suite in a local virtual machine with:

    go run _integration-tests/main.go

The test runner will create the snappy images with `ubuntu-device-flash`, so it
will ask for your password to run this command with `sudo`.

You can also especify more options to customize the image being created, including
the release, the channel and the revision to use. This parameters will be passed
to `ubuntu-device-flash`:

    go run _integration-tests/main.go -release 15.04 -channel stable -revision 3

The default values are suited to testing the most recent version, `rolling` for
release, `edge` for channel and an empty revision, which picks the latest
available.

## Testing snappy from a branch

With the --snappy-from-branch flag, the snappy CLI command will be compiled
from the current branch, copied to the test bed and used during the integration
tests:

    go run _integration-tests/main.go --snappy-from-branch

You can use this flag to test in a remote machine too.

## Filtering the tests to run

With the --filter flag you can select the tests to run. For instance you can
pass MyTestSuite, MyTestSuite.FirstCustomTest or MyTestSuite.*CustomTest:

    go run _integration-tests/main.go --filter MyTestSuite.FirstCustomTest

## Testing a remote machine

You can execute the integration suite in a remote snappy machine with:

    go run _integration-tests/main.go --ip {testbed-ip} --port {testbed-port} \
    --arch {testbed-arch}

The test runner will use `ssh-copy-id` to send your identity file to the
testbed, so it will ask for the password of the ubuntu user in the test bed.

When running in a remote machine, the test runner assumes the test bed is in
the latest rolling edge version, and it will skip all the tests that
require a different version. See the following section for instructions for
setting up a BeagleBone Black as the test bed.

## Testing a BeagleBone Black

First flash the latest rolling edge version into the sd card
(replacing /dev/sdX with the path to your card):

    sudo ubuntu-device-flash core rolling --channel edge --oem beagleblack \
    --developer-mode --enable-ssh -o ubuntu-rolling-edge-armhf-bbb.img

    sudo dd if=ubuntu-rolling-edge-armhf-bbb.img of=/dev/sdX bs=32M
    sync

Then boot the board with the sd card, make sure that it is connected to the
same network as the test runner host, and find the {beaglebone-ip}.

Run the tests with:

    go run _integration-tests/main.go --ip {beaglebone-ip} --arch arm

## Testing an update

You can flash an old image, update to the latest and then run the whole suite
on the updated system. The release, the channel and the revision flags specify
the image that will be flashed, and the target-release and target-channel flags
specify the values to be used in the update.

For example, to update from rolling edge -1 to the latest and then run the
integration tests:

    go run _integration-tests/main.go --snappy-from-branch \
    --revision=-1 --target-release=rolling --target-channel=edge

To update from 15.04 edge to rolling edge and then run the integration tests:

    go run _integration-tests/main.go --snappy-from-branch \
    --release=15.04 --channel=edge \
    --target-release=rolling --target-channel=edge