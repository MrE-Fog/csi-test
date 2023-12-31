/*
Copyright 2018 The Kubernetes Authors.

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

package apitest2

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/kubernetes-csi/csi-test/v5/pkg/sanity"
)

// TestMyDriverWithCustomTargetPaths verifies that CreateTargetDir and
// CreateStagingDir are called a specific number of times.
func TestMyDriverWithCustomTargetPaths(t *testing.T) {
	var createTargetDirCalls, createStagingDirCalls,
		removeTargetDirCalls, removeStagingDirCalls int

	wantCreateTargetCalls := 3
	wantCreateStagingCalls := 3
	wantRemoveTargetCalls := 3
	wantRemoveStagingCalls := 3

	// tmpPath could be a CO specific directory under which all the target dirs
	// are created. For k8s, it could be /var/lib/kubelet/pods under which the
	// mount directories could be created.
	tmpPath := path.Join(os.TempDir(), "csi")
	config := sanity.NewTestConfig()
	config.TargetPath = "foo/target/mount"
	config.StagingPath = "foo/staging/mount"
	config.Address = "/tmp/e2e-csi-sanity.sock"
	config.CreateTargetDir = func(targetPath string) (string, error) {
		createTargetDirCalls++
		targetPath = path.Join(tmpPath, targetPath)
		return targetPath, createTargetDir(targetPath)
	}
	config.CreateStagingDir = func(targetPath string) (string, error) {
		createStagingDirCalls++
		targetPath = path.Join(tmpPath, targetPath)
		return targetPath, createTargetDir(targetPath)
	}
	config.RemoveTargetPath = func(targetPath string) error {
		removeTargetDirCalls++
		return os.RemoveAll(targetPath)
	}
	config.RemoveStagingPath = func(targetPath string) error {
		removeStagingDirCalls++
		return os.RemoveAll(targetPath)
	}

	sanity.Test(t, config)

	if createTargetDirCalls != wantCreateTargetCalls {
		t.Errorf("unexpected number of CreateTargetDir calls:\n(WNT) %d\n(GOT) %d", wantCreateTargetCalls, createTargetDirCalls)
	}

	if createStagingDirCalls != wantCreateStagingCalls {
		t.Errorf("unexpected number of CreateStagingDir calls:\n(WNT) %d\n(GOT) %d", wantCreateStagingCalls, createStagingDirCalls)
	}

	if removeTargetDirCalls != wantRemoveTargetCalls {
		t.Errorf("unexpected number of RemoveTargetDir calls:\n(WNT) %d\n(GOT) %d", wantRemoveTargetCalls, removeTargetDirCalls)
	}

	if removeStagingDirCalls != wantRemoveStagingCalls {
		t.Errorf("unexpected number of RemoveStagingDir calls:\n(WNT) %d\n(GOT) %d", wantRemoveStagingCalls, removeStagingDirCalls)
	}
}

func createTargetDir(targetPath string) error {
	fileInfo, err := os.Stat(targetPath)
	if err != nil && os.IsNotExist(err) {
		return os.MkdirAll(targetPath, 0755)
	} else if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("Target location %s is not a directory", targetPath)
	}

	return nil
}
