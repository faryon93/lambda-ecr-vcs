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
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/valyala/fastjson"
	"strings"
)

// ---------------------------------------------------------------------------------------
//  constants
// ---------------------------------------------------------------------------------------

const (
	ManifestType = "application/vnd.docker.distribution.manifest.v1+json"
)

var (
	ErrMalformedResponse   = errors.New("malformed response")
	ErrHistoryEmpty        = errors.New("image history is empty")
	ErrV1ManifestMissing   = errors.New("missing v1Compatibility manifest")
	ErrConfigObjectMissing = errors.New("missing config object")
	ErrLabelsObjectMissing = errors.New("missing labels object")
	ErrNoEcrClient         = errors.New("ecr client not initialized")
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type VcsInfo struct {
	RepoUrl string `json:"repo_url"`
	Branch  string `json:"branch"`
	Author  string `json:"author"`
	Ref     string `json:"ref"`
	Message string `json:"message"`
}

// ---------------------------------------------------------------------------------------
//  public functions
// ---------------------------------------------------------------------------------------

func GetVcsInfo(client *ecr.ECR, repo string, tag string) (*VcsInfo, error) {
	if client == nil {
		return nil, ErrNoEcrClient
	}

	params := ecr.BatchGetImageInput{
		AcceptedMediaTypes: aws.StringSlice([]string{ManifestType}),
		RepositoryName:     aws.String(repo),
		ImageIds:           []*ecr.ImageIdentifier{{ImageTag: aws.String(tag)}},
	}
	resp, err := client.BatchGetImage(&params)
	if err != nil {
		return nil, err
	}

	// make sure the necessary fields exist
	if len(resp.Images) < 1 || resp.Images[0].ImageManifest == nil {
		return nil, ErrMalformedResponse
	}

	return parseVcsInfo(*resp.Images[0].ImageManifest)
}

// ---------------------------------------------------------------------------------------
//  private functions
// ---------------------------------------------------------------------------------------

func parseVcsInfo(manifest string) (*VcsInfo, error) {
	var p fastjson.Parser
	v, err := p.Parse(manifest)
	if err != nil {
		return nil, err
	}

	history := v.GetArray("history")
	if len(history) == 0 {
		return nil, ErrHistoryEmpty
	}

	rawManifestV1 := history[0].GetStringBytes("v1Compatibility")
	if len(rawManifestV1) == 0 {
		return nil, ErrV1ManifestMissing
	}

	manifestV1, err := p.ParseBytes(rawManifestV1)
	if err != nil {
		return nil, err
	}

	config := manifestV1.Get("config")
	if config == nil {
		return nil, ErrConfigObjectMissing
	}

	labels := config.Get("Labels")
	if labels == nil {
		return nil, ErrLabelsObjectMissing
	}

	vcs := VcsInfo{
		RepoUrl: string(labels.GetStringBytes("org.label-schema.vcs-url")),
		Branch:  string(labels.GetStringBytes("org.factory360.vcs-branch")),
		Author:  string(labels.GetStringBytes("org.factory360.vcs-author")),
		Ref:     string(labels.GetStringBytes("org.label-schema.vcs-ref")),
		Message: string(labels.GetStringBytes("org.factory360.vcs-message")),
	}

	// remove leading and trailing quotation marks
	vcs.Message = strings.TrimPrefix(vcs.Message, "\"")
	vcs.Message = strings.TrimSuffix(vcs.Message, "\"")

	return &vcs, nil
}
