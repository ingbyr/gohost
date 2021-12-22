/*
 @Author: ingbyr
*/

package host

import (
	"testing"
)

type hosts struct {
	name   string
	groups []string
}

type hostTestOp struct {
	initHosts    []hosts
	addGroups    []string
	deleteGroups []string
	addHosts     []hosts
	deleteHosts  []string
	renameHosts  [][2]string
}

func TestManager_CreateRemoveNewHost(t *testing.T) {
	var tests = []hostTestOp{
		{
			initHosts: []hosts{
				{"f3", []string{"g1", "g2", "g3", "g4"}},
				{"f2", []string{"g1", "g2"}},
				{"f1", []string{"g1", "g4"}},
				{"f5", []string{"g5"}},
			},
			deleteGroups: []string{"g4", "g5", "g0"},
		},
	}
	for _, test := range tests {
		for _, initHost := range test.initHosts {
			M.CreateNewHost(initHost.name, initHost.groups, false)
		}
		M.LoadHosts()
		M.DeleteGroups(test.deleteGroups)
	}

	M.LoadHosts()
	M.printHosts()

	M.ChangeGroupName("g1", "gg")
	M.LoadHosts()
	M.printHosts()
}
