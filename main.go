package main

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
//  imports
// ---------------------------------------------------------------------------------------

import (
	"errors"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"

	"log"

	"playground/eventbridge"
	"playground/image"
)

// ---------------------------------------------------------------------------------------
//  types
// ---------------------------------------------------------------------------------------

type Response struct {
	image.VcsInfo
	Image string `json:"image"`
}

// ---------------------------------------------------------------------------------------
//  global variables
// ---------------------------------------------------------------------------------------

var (
	EcrClient *ecr.ECR
)

// ---------------------------------------------------------------------------------------
//  handler functions
// ---------------------------------------------------------------------------------------

func Handle(event eventbridge.EcrEvent) (*Response, error) {
	if EcrClient == nil {
		log.Println("ecr client not initialized")
		return nil, errors.New("ecr client not initialized")
	}

	vcs, err := image.GetVcsInfo(EcrClient, event.Detail.Repository, event.Detail.ImageTag)
	if err != nil {
		return nil, err
	}

	r := Response{
		VcsInfo: *vcs,
		Image:   event.GetFullImage(),
	}

	log.Println(r)

	return &r, nil
}

// ---------------------------------------------------------------------------------------
//  package initialization
// ---------------------------------------------------------------------------------------

func init() {
	sess, err := session.NewSession(&aws.Config{})
	if err != nil {
		log.Println("failed to create aws sess:", err.Error())
		return
	}

	EcrClient = ecr.New(sess)
}

// ---------------------------------------------------------------------------------------
//  application entry
// ---------------------------------------------------------------------------------------

func main() {
	lambda.Start(Handle)
}
