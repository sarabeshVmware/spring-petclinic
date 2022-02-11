// Copyright 2021 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package git

import (
	"fmt"
	"log"
	"strings"

	"gitlab.eng.vmware.com/tap/tap-packages/suite/pkg/utils/linux_util"
)

func GitClone(path string, repo string) error {
	log.Printf("cloning repo %s at %s", repo, path)

	// execute cmd
	cmd := fmt.Sprintf("cd %s; git clone %s", path, repo)
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while cloning repo %s at %s", repo, path)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("repo %s cloned at %s", repo, path)
		log.Printf("output: %s", output)
	}

	return err
}

func GitAdd(path string, files []string) error {
	log.Printf("adding files %s for repo at %s to git index", files, path)

	// execute cmd
	cmd := fmt.Sprintf("cd %s; git add", path)
	for _, file := range files {
		cmd += fmt.Sprintf(" %s", file)
	}
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while adding files %s for repo at %s to git index", files, path)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("files %s added for repo at %s to git index", files, path)
		log.Printf("output: %s", output)
	}

	return err
}

func GitCommit(path string, message string) error {
	log.Printf("committing repo index at %s with message %s", path, message)

	// execute cmd
	cmd := fmt.Sprintf(`cd %s; git commit -m "%s"`, path, message)
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf(`error while committing repo index at %s with message "%s"`, path, message)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf(`commited repo index at %s with message "%s"`, path, message)
		log.Printf("output: %s", output)
	}

	return err
}

func GitPush(path string, force bool) error {
	log.Printf("pushing commits for repo at %s", path)

	// execute cmd
	cmd := fmt.Sprintf("cd %s; git push", path)
	if force {
		cmd += " --force"
	}
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while pushing commits for repo at %s", path)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("pushing commits for repo at %s", path)
		log.Printf("output: %s", output)
	}

	return err
}

func GitResetFromHead(path string, count int) error {
	log.Printf("resetting current HEAD by %d for repo at %s", count, path)

	// execute cmd
	cmd := fmt.Sprintf("cd %s; git reset --hard HEAD~%d", path, count)
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while resetting current HEAD by %d for repo at %s", count, path)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("reset current HEAD by %d for repo at %s", count, path)
		log.Printf("output: %s", output)
	}

	return err
}

func GitConfig(username string, email string) error {
	log.Printf("setting git config (username %s, email %s)", username, email)

	// execute cmd
	cmd := fmt.Sprintf("git config --global user.name %s; git config --global user.email %s", username, email)
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while setting git config (username %s, email %s)", username, email)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("git config (username %s, email %s)", username, email)
		log.Printf("output: %s", output)
	}

	return err
}

func GitSetRemoteUrl(path string, access_token string, repository string) error {
	log.Printf("setting git remote url (%s) for repo at %s", repository, path)

	// execute cmd
	url := strings.Split(repository, "//")[1]
	cmd := fmt.Sprintf("cd %s; git remote set-url origin https://%s@%s", path, access_token, url)
	output, err := linux_util.ExecuteCmdInBashMode(cmd)
	if err != nil {
		log.Printf("error while setting git remote url (%s) for repo at %s", repository, path)
		log.Printf("error: %s", err)
		log.Printf("output: %s", output)
	} else {
		log.Printf("git remote url (%s) set for repo at %s", repository, path)
		log.Printf("output: %s", output)
	}

	return err
}
