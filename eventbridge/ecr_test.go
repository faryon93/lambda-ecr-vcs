package eventbridge

// lambda-ecr-push-vcs
// Copyright (C) 2020 Maximilian Pachl

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// ---------------------------------------------------------------------------------------
//  imports
// ---------------------------------------------------------------------------------------

import (
	"testing"
)

// ---------------------------------------------------------------------------------------
//  tests
// ---------------------------------------------------------------------------------------

func TestEcrEvent_GetRegistryUrl(t *testing.T) {
	tests := []struct {
		Name        string
		Region      string
		AccountId   string
		ExpectedUrl string
	}{
		{"default", "eu-central-1", "1234567", "1234567.dkr.ecr.eu-central-1.amazonaws.com"},
		{"no-account-id", "eu-central-1", "", ""},
		{"no-region", "", "1234567", ""},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			e := EcrEvent{
				Region:    test.Region,
				AccountId: test.AccountId,
			}

			url := e.GetRegistryUrl()
			if url != test.ExpectedUrl {
				t.Errorf("GetRegistryUrl should return \"%s\" but returned \"%s\"",
					test.ExpectedUrl, url)
			}
		})
	}
}
