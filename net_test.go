/*
 * Copyright (c) 2019 Manabu Sonoda
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

package dtap_test

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mikiT/dtap"
)

func TestNet(t *testing.T) {
	n := dtap.Net{
		IP:           net.ParseIP("192.168.0.1"),
		PrefixLength: 32,
	}
	bs, err := n.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, string(bs), `"192.168.0.1/32"`)

	n = dtap.Net{}
	err = n.UnmarshalJSON(bs)
	assert.NoError(t, err)
	assert.Equal(t, n.String(), "192.168.0.1/32")

	n = dtap.Net{
		IP:           net.ParseIP("2001:db8::1"),
		PrefixLength: 48,
	}
	bs, err = n.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, string(bs), `"2001:db8::1/48"`)

	n = dtap.Net{}
	err = n.UnmarshalJSON(bs)
	assert.NoError(t, err)
	assert.Equal(t, n.String(), "2001:db8::1/48")

	n = dtap.Net{}
	bs, err = n.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, string(bs), `"<nil>"`)

	n = dtap.Net{}
	err = n.UnmarshalJSON(bs)
	assert.Nil(t, n.IP)

}
