//
// Copyright 2016-2017 Bryan T. Meyers <bmeyers@datadrake.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package core

import "regexp"

var APACHE_COMBINED = regexp.MustCompile("^(\\S*).(\\S*).(\\S*).\\[(.*)\\].\"(.*)\".(\\d{3}).(\\d*).\"(.*)\".\"(.*)\"$")
var APACHE_REQUEST = regexp.MustCompile("^(?:(\\S*) )?(.*?)(?: (\\S*))?$")
var LOG2DB_TABLES = []string{"destinations", "sources", "txns", "users", "user_agents"}
var APACHE_TIME_LAYOUT = "02/Jan/2006:15:04:05 -0700"

var APACHE_COMBINED_FORMAT = "%s %s %s [%s] \"%s\" %d %d \"%s\" \"%s\"\n"

func SliceContains(vs []string, v string) bool {
	for _, curr := range vs {
		if curr == v {
			return true
		}
	}
	return false
}
