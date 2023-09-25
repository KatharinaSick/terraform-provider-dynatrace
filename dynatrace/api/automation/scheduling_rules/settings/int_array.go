/**
* @license
* Copyright 2020 Dynatrace LLC
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
 */

package scheduling_rules

import (
	"encoding/json"
)

type IntArray []int

func (me *IntArray) UnmarshalJSON(data []byte) (err error) {
	var ints []int
	if err = json.Unmarshal(data, &ints); err != nil {
		if err.Error() != "json: cannot unmarshal number into Go value of type []int" {
			return err
		}
		var i int
		if err = json.Unmarshal(data, &i); err != nil {
			return err
		}
		ints = []int{i}
	}
	for _, i := range ints {
		*me = append(*me, i)
	}
	return nil
}
