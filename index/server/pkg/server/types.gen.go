//
// Copyright 2023 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0. **DO NOT EDIT**
package server

import (
	"github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/registry-support/index/generator/schema"
)

// Optional list of processor architectures that the devfile supports, empty list suggests that the devfile can be used on any architecture
type ArchParam = []string

// Devfile defines model for Devfile.
type Devfile v1alpha2.Devfile

// Optional devfile icon, can be a URI or a relative path in the project
type IconParam = string

// IndexParams defines parameters for index endpoints.
type IndexParams struct {
	// Optional list of processor architectures that the devfile supports, empty list suggests that the devfile can be used on any architecture
	Archs *ArchParam `json:"arch,omitempty"`

	// Optional devfile icon, can be a URI or a relative path in the project
	IconType *IconParam `json:"icon,omitempty"`
}

// DevfileErrorResponse defines model for devfileErrorResponse.
type DevfileErrorResponse struct {
	Error  *string `json:"error,omitempty"`
	Status *string `json:"status,omitempty"`
}

// DevfileNotFoundResponse defines model for devfileNotFoundResponse.
type DevfileNotFoundResponse struct {
	Status *string `json:"status,omitempty"`
}

// DevfileResponse defines model for devfileResponse.
type DevfileResponse = Devfile

// HealthResponse defines model for healthResponse.
type HealthResponse struct {
	Message *string `json:"message,omitempty"`
}

// IndexResponse defines model for indexResponse.
type IndexResponse schema.Schema

// V2IndexResponse defines model for v2IndexResponse.
type V2IndexResponse schema.Schema

// ServeDevfileIndexV1Params defines parameters for ServeDevfileIndexV1.
type ServeDevfileIndexV1Params struct {
	// The target architecture filter
	Archs *ArchParam `form:"arch,omitempty" json:"arch,omitempty"`

	// The icon type filter
	IconType *IconParam `form:"icon,omitempty" json:"icon,omitempty"`
}

// ServeDevfileIndexV1WithTypeParams defines parameters for ServeDevfileIndexV1WithType.
type ServeDevfileIndexV1WithTypeParams struct {
	// The target architecture filter
	Archs *ArchParam `form:"arch,omitempty" json:"arch,omitempty"`

	// The icon type filter
	IconType *IconParam `form:"icon,omitempty" json:"icon,omitempty"`
}

// ServeDevfileIndexV2Params defines parameters for ServeDevfileIndexV2.
type ServeDevfileIndexV2Params struct {
	// The target architecture filter
	Archs *ArchParam `form:"arch,omitempty" json:"arch,omitempty"`

	// The icon type filter
	IconType *IconParam `form:"icon,omitempty" json:"icon,omitempty"`
}

// ServeDevfileIndexV2WithTypeParams defines parameters for ServeDevfileIndexV2WithType.
type ServeDevfileIndexV2WithTypeParams struct {
	// The target architecture filter
	Archs *ArchParam `form:"arch,omitempty" json:"arch,omitempty"`

	// The icon type filter
	IconType *IconParam `form:"icon,omitempty" json:"icon,omitempty"`
}