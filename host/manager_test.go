/*
 @Author: ingbyr
*/

package host

import (
	"github.com/ingbyr/gohost/myfs"
	"testing"
)

func TestManager_CreateRemoveNewHost(t *testing.T) {
	M.SetFs(myfs.NewMemFs())
	M.CreateNewHost("f3", []string{"g1", "g2", "g3"}, false)
	M.CreateNewHost("f2", []string{"g1", "g2"}, false)
	M.CreateNewHost("f1", []string{"g1", "g4"}, false)
	M.CreateNewHost("f5", []string{"g5"}, false)
	M.LoadHosts()
	M.DisplayGroups()

	M.DeleteGroups([]string{"g4", "g5", "g0"})
	M.LoadHosts()
	M.DisplayGroups()
	M.DisplayHosts()
	M.printNodes()
}
