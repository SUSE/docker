/*
 *
 * Copyright 2018 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package envconfig contains grpc settings configured by environment variables.
package envconfig

import (
	"os"
	"strings"
)

const (
	prefix          = "GRPC_GO_"
	txtErrIgnoreStr = prefix + "IGNORE_TXT_ERRORS"
)

var (
	// TXTErrIgnore is set if TXT errors should be ignored ("GRPC_GO_IGNORE_TXT_ERRORS" is not "false").
	TXTErrIgnore = !strings.EqualFold(os.Getenv(txtErrIgnoreStr), "false")

	// DisableStrictPathChecking indicates whether strict path checking is
	// disabled. This feature can be disabled by setting the environment
	// variable GRPC_GO_EXPERIMENTAL_DISABLE_STRICT_PATH_CHECKING to "true".
	//
	// When strict path checking is enabled, gRPC will reject requests with
	// paths that do not conform to the gRPC over HTTP/2 specification found at
	// https://github.com/grpc/grpc/blob/master/doc/PROTOCOL-HTTP2.md.
	//
	// When disabled, gRPC will allow paths that do not contain a leading slash.
	// Enabling strict path checking is recommended for security reasons, as it
	// prevents potential path traversal vulnerabilities.
	//
	// A future release will remove this environment variable, enabling strict
	// path checking behavior unconditionally.
	DisableStrictPathChecking = boolFromEnv("GRPC_GO_EXPERIMENTAL_DISABLE_STRICT_PATH_CHECKING", false)

	// EnableReceiveBufferCompaction controls whether the transport
	// compacts small buffered DATA frame payloads in a stream's receive
	// backlog into a single larger buffer once their per-message tracking
	// overhead becomes excessive relative to actual payload size. This
	// defends against a peer fragmenting a stream's payload into many
	// tiny (e.g. 1-byte) DATA frames to cause heap memory amplification.
	// Can be disabled by setting the environment variable
	// GRPC_GO_EXPERIMENTAL_ENABLE_RECEIVE_BUFFER_COMPACTION to "false".
	EnableReceiveBufferCompaction = boolFromEnv("GRPC_GO_EXPERIMENTAL_ENABLE_RECEIVE_BUFFER_COMPACTION", true)
)

func boolFromEnv(envVar string, def bool) bool {
	if def {
		// The default is true; return true unless the variable is "false".
		return !strings.EqualFold(os.Getenv(envVar), "false")
	}
	// The default is false; return false unless the variable is "true".
	return strings.EqualFold(os.Getenv(envVar), "true")
}
