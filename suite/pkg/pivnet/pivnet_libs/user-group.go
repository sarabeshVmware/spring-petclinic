// add-user-group                    Add user group to release (aliases: aug)
// add-user-group-member             Add user group member to group (aliases: augm)
// create-user-group                 Create user group (aliases: cug)
// delete-user-group                 Delete user group (aliases: dug)
// remove-user-group                 Remove user group from release (aliases: rug)
// remove-user-group-member          Remove user group member from group (aliases: rugm)
// update-user-group                 Update user group (aliases: uug)
// user-group                        Show user group (aliases: ug)
// user-groups                       List user groups (aliases: ugs)

package pivnet_libs

import (
	"encoding/json"
	"fmt"
	"log"

	linux_util "gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

type ListUserGroupsOutput []struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func ListUserGroups(productSlug string) ListUserGroupsOutput {

	cmd := fmt.Sprintf("pivnet-cli user-groups --product-slug= %s --format json", productSlug)
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil {
		log.Println("something bad happened")
	}
	in := []byte(response)
	var raw ListUserGroupsOutput
	if err := json.Unmarshal(in, &raw); err != nil {
		panic(err)
	}
	return raw
}

func AddUserGroup(productSlug string, releaseVersion string, userGroupId int) {
	cmd := fmt.Sprintf("pivnet-cli add-user-group --product-slug=%s --release-version %s --user-group-id=%d --format json", productSlug, releaseVersion, userGroupId)
	response, err := linux_util.ExecuteCmd(cmd)
	if err != nil || response != "" {
		log.Println("something bad happened")
	}
}
