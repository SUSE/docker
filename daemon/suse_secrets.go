/*
 * suse-secrets: patch for Docker to implement SUSE secrets
 * Copyright (C) 2017-2021 SUSE LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package daemon

import (
	"strings"

	"github.com/moby/moby/v2/daemon/container"

	swarmtypes "github.com/moby/moby/api/types/swarm"

	"github.com/sirupsen/logrus"
)

// clearSuseSecrets removes any SecretReferences which were added by us
// explicitly (this is detected by checking that the prefix has a 'suse_'
// prefix, which is a prefix that cannot exist for normal swarm secrets). See
// bsc#1057743 and bsc#1244035.
func (daemon *Daemon) clearSuseSecrets(c *container.Container) {
	var without []*swarmtypes.SecretReference
	for _, secret := range c.SecretReferences {
		if strings.HasPrefix(secret.SecretID, "suse_") {
			logrus.Debugf("SUSE:secrets :: removing 'old' suse secret %q from container %q", secret.SecretID, c.ID)
			continue
		}
		without = append(without, secret)
	}
	c.SecretReferences = without
}
