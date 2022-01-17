package netrc_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/leejones/netrc"
)

// netrc format docs:
//	* https://everything.curl.dev/usingcurl/netrc
//	* https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html
//	* https://www.ibm.com/docs/en/aix/7.1?topic=formats-netrc-file-format-tcpip
//
// The doc from IBM mentions, "The .netrc can contain the following entries
// (separated by spaces, tabs, or new lines)". "Entries" means "machine",
// "login", "password" etc.

func TestGetMachineFound(t *testing.T) {
	t.Parallel()
	netrcFile, err := netrc.NewFile(
		netrc.WithFile("testdata/netrc"),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := netrc.Credentials{
		Username: "Blue",
		Password: "Yellow",
	}
	got, err := netrcFile.Get("green.example.com")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetMachineNotFound(t *testing.T) {
	t.Parallel()
	netrcFile, err := netrc.NewFile(
		netrc.WithFile("testdata/netrc"),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = netrcFile.Get("not-found.example.com")
	if err == nil {
		t.Error("expected an error, but got none")
	}
}
func TestGetMultipleMachinesReturnsFirst(t *testing.T) {
	t.Parallel()
	f, err := netrc.NewFile(
		netrc.WithFile("testdata/netrc"),
	)
	if err != nil {
		t.Error(err)
	}
	want := netrc.Credentials{
		Username: "Red",
		Password: "Yellow",
	}
	got, err := f.Get("orange.example.com")
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetMachineOneLineFormat(t *testing.T) {
	t.Parallel()
	f, err := netrc.NewFile(
		netrc.WithFile("testdata/netrc"),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := netrc.Credentials{
		Username: "Blue",
		Password: "Red",
	}
	got, err := f.Get("purple.example.com")
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
