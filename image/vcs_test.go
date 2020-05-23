package image

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

func TestGetVcsInfo(t *testing.T) {
	testManifest := `{"history": [{"v1Compatibility": "{\"config\": {\"Labels\": {\"org.factory360.vcs-message\": \"org.factory360.vcs-message\", \"org.label-schema.vcs-url\": \"org.label-schema.vcs-url\", \"org.factory360.vcs-branch\": \"org.factory360.vcs-branch\", \"org.factory360.vcs-author\": \"org.factory360.vcs-author\", \"org.label-schema.vcs-ref\": \"org.label-schema.vcs-ref\"}}}"}]}`
	vcs, err := parseVcsInfo(testManifest)
	if err != nil {
		t.Error("parseVcsInfo should not return an error:", err.Error())
		return
	}

	eq(t, "Message", vcs.RepoUrl, "org.label-schema.vcs-url")
	eq(t, "Message", vcs.Branch, "org.factory360.vcs-branch")
	eq(t, "Message", vcs.Author, "org.factory360.vcs-author")
	eq(t, "Message", vcs.Ref, "org.label-schema.vcs-ref")
	eq(t, "Message", vcs.Message, "org.factory360.vcs-message")
}

func TestGetVcsInfo_MessageQuotes(t *testing.T) {
	testManifest := `{"history": [{"v1Compatibility": "{\"config\": {\"Labels\": {\"org.factory360.vcs-message\": \"\\\"org.factory360.vcs-message\\\"\"}}}"}]}`
	vcs, err := parseVcsInfo(testManifest)
	if err != nil {
		t.Error("parseVcsInfo should not return an error:", err.Error())
		return
	}

	if vcs.Message != "org.factory360.vcs-message" {
		t.Errorf("vcs.Message should trim leading/trailing quotes but retunred: \"%s\"",
			vcs.Message)
	}
}

// ---------------------------------------------------------------------------------------
//  helpers
// ---------------------------------------------------------------------------------------

func eq(t *testing.T, field string, f string, v interface{}) {
	if f != v {
		t.Errorf("field \"%s\" should be \"%s\", but is \"%s\"", field, v, f)
	}
}
