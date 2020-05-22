package eventbridge

// playground
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
//  types
// ---------------------------------------------------------------------------------------

type EcrEvent struct {
	AccountId string `json:"account"`
	Region    string `json:"region"`

	DetailType string `json:"detail-type"`
	Detail     struct {
		Repository  string `json:"repository-name"`
		ImageDigest string `json:"image-digest"`
		ImageTag    string `json:"image-tag"`
	} `json:"detail"`
}

// ---------------------------------------------------------------------------------------
//  public members
// ---------------------------------------------------------------------------------------

func (e *EcrEvent) GetFullImage() string {
	return e.GetRegistryUrl() + "/" + e.GetImage()
}

func (e *EcrEvent) GetRegistryUrl() string {
	if e.AccountId == "" || e.Region == "" {
		return ""
	}

	return e.AccountId + "." + "dkr.ecr." + e.Region + ".amazonaws.com"
}

func (e *EcrEvent) GetImage() string {
	return e.Detail.Repository + ":" + e.Detail.ImageTag
}
