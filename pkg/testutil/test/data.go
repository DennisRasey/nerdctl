/*
   Copyright The containerd Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package test

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// Contains the implementation of the Data interface

type data struct {
	config map[ConfigKey]ConfigValue

	system map[SystemKey]SystemValue

	labels  map[string]string
	testID  string
	tempDir string
}

func (dt *data) WithConfig(key ConfigKey, value ConfigValue) Data {
	if dt.config == nil {
		dt.config = make(map[ConfigKey]ConfigValue)
	}
	dt.config[key] = value
	return dt
}

func (dt *data) ReadConfig(key ConfigKey) ConfigValue {
	if dt.config == nil {
		dt.config = make(map[ConfigKey]ConfigValue)
	}
	if val, ok := dt.config[key]; ok {
		return val
	}
	return ""
}

func (dt *data) Get(key string) string {
	if dt.labels == nil {
		dt.labels = map[string]string{}
	}
	return dt.labels[key]
}

func (dt *data) Set(key string, value string) Data {
	if dt.labels == nil {
		dt.labels = map[string]string{}
	}
	dt.labels[key] = value
	return dt
}

func (dt *data) Identifier() string {
	return dt.testID
}

func (dt *data) TempDir() string {
	return dt.tempDir
}

func (dt *data) adopt(parent Data) {
	for k, v := range parent.getLabels() {
		// Only copy keys that are not set already
		if _, ok := dt.labels[k]; !ok {
			dt.Set(k, v)
		}
	}
	for k, v := range parent.getConfig() {
		// Only copy keys that are not set already
		if _, ok := dt.config[k]; !ok {
			dt.WithConfig(k, v)
		}
	}
}

func (dt *data) Sink(key SystemKey, value SystemValue) {
	if _, ok := dt.system[key]; !ok {
		dt.system[key] = value
	} else {
		// XXX should we really panic?
		panic(fmt.Sprintf("Unable to set system key %s multiple times", key))
	}
}

func (dt *data) Surface(key SystemKey) SystemValue {
	if v, ok := dt.system[key]; ok {
		return v
	}
	// XXX should we really panic?
	panic(fmt.Sprintf("Unable to retrieve system key %s", key))
}

func (dt *data) getLabels() map[string]string {
	return dt.labels
}

func (dt *data) getConfig() map[ConfigKey]ConfigValue {
	return dt.config
}

func defaultIdentifierHashing(names ...string) string {
	// Notes: identifier MAY be used for namespaces, image names, etc.
	// So, the rules are stringent on what it can contain.
	replaceWith := []byte("-")
	name := strings.ToLower(strings.Join(names, string(replaceWith)))
	// Ensure we have a unique identifier despite characters replacements (well, as unique as the names collection being passed)
	signature := fmt.Sprintf("%x", sha256.Sum256([]byte(name)))[0:8]
	// Make sure we do not use any unsafe characters
	safeName := regexp.MustCompile(`[^a-z0-9-]+`)
	// And we avoid repeats of the separator
	noRepeat := regexp.MustCompile(fmt.Sprintf(`[%s]{2,}`, replaceWith))
	sn := safeName.ReplaceAll([]byte(name), replaceWith)
	sn = noRepeat.ReplaceAll(sn, replaceWith)
	// Do not allow trailing or leading dash (as that may stutter)
	name = strings.Trim(string(sn), string(replaceWith))
	// Ensure we will never go above 76 characters in length (with signature)
	if len(name) > 67 {
		name = name[0:67]
	}
	return name + "-" + signature
}

// TODO: allow to pass custom hashing methods?
func configureData(t *testing.T, seedData Data, parent Data) Data {
	if seedData == nil {
		seedData = &data{}
	}
	dat := &data{
		config:  seedData.getConfig(),
		labels:  seedData.getLabels(),
		tempDir: t.TempDir(),
		testID:  defaultIdentifierHashing(t.Name()),
	}
	if parent != nil {
		dat.adopt(parent)
	}
	return dat
}
