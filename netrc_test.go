package netrc_test

import (
	"testing"

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

func TestFetchReturnsCredentials(t *testing.T) {
	t.Parallel()
	f, err := netrc.NewFile(
		netrc.WithFile("testdata/netrc"),
	)
	if err != nil {
		t.Fatal(err)
	}
	want := "Perry"
	credentials := f.Fetch("perry.example.com")
	got := credentials.Login
	if want != got {
		t.Errorf("Want: %v, got: %v", want, got)
	}

	want = "WheresPerry?"
	got = credentials.Password
	if want != got {
		t.Errorf("Want: %v, got: %v", want, got)
	}
}

func TestFetchReturnsCredentialsWhenOneLineFormat(t *testing.T) {
	t.Parallel()
	f, err := netrc.NewFile(
		netrc.WithFile("testdata/netrc"),
	)
	if err != nil {
		t.Fatal(err)
	}
	credentials := f.Fetch("isabella.example.com")
	want := "Isabella"
	got := credentials.Login
	if want != got {
		t.Errorf("Want: %v, got: %v", want, got)
	}
}
