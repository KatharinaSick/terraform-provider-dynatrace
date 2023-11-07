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

package export

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

type ReplaceFunc func(s string, cnt int) string

func DefaultReplace(s string, cnt int) string {
	return fmt.Sprintf("%s(%d)", s, cnt)
}

func ResourceName(s string, cnt int) string {
	return fmt.Sprintf("%s_%d", s, cnt)
}

type UniqueNamer interface {
	Name(string) string
	Replace(ReplaceFunc) UniqueNamer
	BlockName(string)
	SetNameWritten(string) bool
}

func NewUniqueNamer() UniqueNamer {
	return &nameCounter{m: map[string]int{}, mFull: map[string]bool{}, mWritten: map[string]bool{}, mutex: new(sync.Mutex)}
}

type nameCounter struct {
	m        map[string]int
	mFull    map[string]bool
	mWritten map[string]bool
	replace  ReplaceFunc
	mutex    *sync.Mutex
}

func (me *nameCounter) Replace(replace ReplaceFunc) UniqueNamer {
	me.replace = replace
	return me
}

var cntreg = regexp.MustCompile(`^(.*)_(\d)*$`)

func (me *nameCounter) Name(name string) string {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	// If we're getting passed `Monitor_Name_1` here, we shorten it to `Monitor_Name`
	if matches := cntreg.FindStringSubmatch(name); matches != nil {
		if matches[1] != "" {
			name = matches[1]
		}
	}

	outName := ""

	for {
		cnt, found := me.m[strings.ToLower(name)]
		if !found {
			me.m[strings.ToLower(name)] = 0
			outName = name
			break
		} else {
			me.m[strings.ToLower(name)] = cnt + 1
		}

		if me.replace == nil {
			outName = DefaultReplace(name, me.m[strings.ToLower(name)])
		} else {
			outName = me.replace(name, me.m[strings.ToLower(name)])
		}

		_, foundFull := me.mFull[outName]
		if foundFull {
			continue
		}

		me.mFull[outName] = true
		break
	}

	return outName
}

func (me *nameCounter) BlockName(name string) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	me.mFull[name] = true

}

func (me *nameCounter) SetNameWritten(name string) bool {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	isWritten := me.mWritten[name]
	me.mWritten[name] = true

	return isWritten

}
