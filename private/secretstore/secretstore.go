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
	"fmt"
	"runtime"

	chelper "github.com/docker/docker-credential-helpers/client"

	"github.com/econbits/econkit/private/config"
	"github.com/git-lfs/go-netrc/netrc"
)

type Store string

const (
	OSX_KEYCHAIN         Store = "osxkeychain"
	D_BUS_SECRET_SERVICE Store = "secretservice"
	WINCRED              Store = "wincred"
	PASS                 Store = "pass"
	NETRC                Store = "netrc"
)

type Credentials struct {
	Domain   string
	Username string
	Secret   string
}

func Default() Store {
	if runtime.GOOS == "linux" {
		return D_BUS_SECRET_SERVICE
	}
	if runtime.GOOS == "darwin" {
		return OSX_KEYCHAIN
	}
	if runtime.GOOS == "windows" {
		return WINCRED
	}
	return NETRC
}

func Get(store Store, domain string) (*Credentials, error) {
	if store == NETRC {
		return getNetrc(domain)
	}
	return getFromHelper(store, domain)
}

func getFromHelper(store Store, domain string) (*Credentials, error) {
	storeid := fmt.Sprintf("docker-credential-%s", store)
	p := chelper.NewShellProgramFunc(storeid)

	creds, err := chelper.Get(p, domain)
	if err != nil {
		return nil, err
	}
	return &Credentials{Domain: domain, Username: creds.Username, Secret: creds.Secret}, nil
}

func getNetrc(domain string) (*Credentials, error) {
	rc, err := netrc.ParseFile(config.NETRC)
	if err != nil {
		return nil, err
	}
	machine := rc.FindMachine(domain)
	if machine == nil {
		return nil, fmt.Errorf("machine %s is not found in %s", domain, config.NETRC)
	}
	return &Credentials{Domain: domain, Username: machine.Login, Secret: machine.Password}, nil
}
