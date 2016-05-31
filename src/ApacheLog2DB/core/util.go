package core

import "regexp"

var APACHE_COMBINED = regexp.MustCompile("^(\\S*).(\\S*).(\\S*).\\[(.*)\\].\"([^\"]*)\".(\\d{3}).(\\d*).\"([^\"]*)\".\"([^\"]*)\"$")
var LOG2DB_TABLES = []string{"destinations", "sources", "txns", "users", "user_agents"}
var APACHE_TIME_LAYOUT = "02/Jan/2006:15:04:05 -0700"

var APACHE_COMBINED_FORMAT = "%s %s %s [%s] \"%s %s %s\" %d %d \"%s\" \"%s\""
func SliceContains(vs []string, v string) bool {
	for _, curr := range vs {
		if curr == v {
			return true
		}
	}
	return false
}
