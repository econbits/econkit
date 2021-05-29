/*
Copyright © 2021 Germán Fuentes Capella

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package secretstore

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/econbits/econkit/private/config"
)

func TestDefaultStore(t *testing.T) {
	if runtime.GOOS == "linux" {
		if Default() != D_BUS_SECRET_SERVICE {
			t.Fatalf("D-BUS Secrete Service should be the default store in linux; found: %s", Default())
		}
	} else if runtime.GOOS == "darwin" {
		if Default() != OSX_KEYCHAIN {
			t.Fatalf("OSX Keychain should be the default store in macos; found: %s", Default())
		}
	} else if runtime.GOOS == "windows" {
		if Default() != WINCRED {
			t.Fatalf("Wincred should be the default store in Windows; found: %s", Default())
		}
	} else {
		if Default() != NETRC {
			t.Fatalf(".netrc should be the default store; found: %s", Default())
		}
	}
}

func TestNetrcWithValidDomain(t *testing.T) {
	_, thisfile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Error getting current file name")
	}
	mydir := filepath.Dir(thisfile)

	baseNetrc := config.NETRC
	defer func() { config.NETRC = baseNetrc }()

	config.NETRC = filepath.Join(mydir, "..", "..", "test", "data", "netrc.txt")

	cred, err := Get(NETRC, "econbits.com")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cred.Domain != "econbits.com" {
		t.Errorf("Unexpected domain: %s", cred.Domain)
	}
	if cred.Username != "username" {
		t.Errorf("Unexpected username: %s", cred.Username)
	}
	if cred.Secret != "secret" {
		t.Errorf("Unexpected secret: %s", cred.Secret)
	}
}

func TestNetrcWithInvalidDomain(t *testing.T) {
	_, thisfile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Error getting current file name")
	}
	mydir := filepath.Dir(thisfile)

	baseNetrc := config.NETRC
	defer func() { config.NETRC = baseNetrc }()

	config.NETRC = filepath.Join(mydir, "..", "..", "test", "data", "netrc.txt")

	_, err := Get(NETRC, "this-is-not-a-good-domain")
	if err == nil {
		t.Fatal("Expected error; none found")
	}
}
