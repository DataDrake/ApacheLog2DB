package core

import "regexp"

var APACHE_COMBINED = regexp.MustCompile("^(\\S*).(\\S*).(\\S*).\\[(.*)\\].\"([^\"]*)\".(\\d{3}).(\\d*).\"([^\"]*)\".\"([^\"]*)\"$")
var LOG2DB_TABLES = []string{"destinations", "sources", "transactions", "users", "user_agents"}
var APACHE_TIME_LAYOUT = "2006-01-02T15:04:05.000Z"

func SliceContains(vs []string, v string) bool {
	for _, curr := range vs {
		if curr == v {
			return true
		}
	}
	return false
}
