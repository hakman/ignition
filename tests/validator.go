// Copyright 2017 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package blackbox

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"

	"github.com/flatcar-linux/ignition/internal/exec/util"
	"github.com/flatcar-linux/ignition/tests/types"
)

func regexpSearch(itemName, pattern string, data []byte) (string, error) {
	re := regexp.MustCompile(pattern)
	match := re.FindSubmatch(data)
	if len(match) < 2 {
		return "", fmt.Errorf("couldn't find %s", itemName)
	}
	return string(match[1]), nil
}

func getPartitionSet(device string) (map[int]struct{}, error) {
	sgdiskOverview, err := exec.Command("sgdisk", "-p", device).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("sgdisk -p %s failed: %v", device, err)
	}

	//What this regex means:       num      start    end    size,code,name
	re := regexp.MustCompile("\n\\W+(\\d+)\\W+\\d+\\W+\\d+\\W+\\d+.*")
	ret := map[int]struct{}{}
	for _, match := range re.FindAllStringSubmatch(string(sgdiskOverview), -1) {
		if len(match) == 0 {
			continue
		}
		if len(match) != 2 {
			return nil, fmt.Errorf("Invalid regex result from parsing sgdisk")
		}
		num, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}
		ret[num] = struct{}{}
	}
	return ret, nil
}

func validateDisk(t *testing.T, d types.Disk) error {
	partitionSet, err := getPartitionSet(d.Device)
	if err != nil {
		return err
	}

	for _, e := range d.Partitions {
		if e.TypeCode == "blank" {
			continue
		}

		if _, ok := partitionSet[e.Number]; !ok {
			t.Errorf("Partition %d is missing", e.Number)
		}
		delete(partitionSet, e.Number)

		sgdiskInfo, err := exec.Command(
			"sgdisk", "-i", strconv.Itoa(e.Number),
			d.Device).CombinedOutput()
		if err != nil {
			t.Error("sgdisk -i", strconv.Itoa(e.Number), err)
			return nil
		}

		actualGUID, err := regexpSearch("GUID", "Partition unique GUID: (?P<partition_guid>[\\d\\w-]+)", sgdiskInfo)
		if err != nil {
			return err
		}
		actualTypeGUID, err := regexpSearch("type GUID", "Partition GUID code: (?P<partition_code>[\\d\\w-]+)", sgdiskInfo)
		if err != nil {
			return err
		}
		actualSectors, err := regexpSearch("partition size", "Partition size: (?P<sectors>\\d+) sectors", sgdiskInfo)
		if err != nil {
			return err
		}
		actualLabel, err := regexpSearch("partition name", "Partition name: '(?P<name>[\\d\\w-_]+)'", sgdiskInfo)
		if err != nil {
			return err
		}

		// have to align the size to the nearest sector alignment boundary first
		expectedSectors := types.Align(e.Length, d.Alignment)

		if e.TypeGUID != "" && formatUUID(e.TypeGUID) != formatUUID(actualTypeGUID) {
			t.Error("TypeGUID does not match!", e.TypeGUID, actualTypeGUID)
		}
		if e.GUID != "" && formatUUID(e.GUID) != formatUUID(actualGUID) {
			t.Error("GUID does not match!", e.GUID, actualGUID)
		}
		if e.Label != actualLabel {
			t.Error("Label does not match!", e.Label, actualLabel)
		}
		if strconv.Itoa(expectedSectors) != actualSectors {
			t.Error(
				"Sectors does not match!", expectedSectors, actualSectors)
		}
	}

	if len(partitionSet) != 0 {
		t.Error("Disk had extra partitions", partitionSet)
	}

	// TODO: inspect the disk without triggering partition rescans so we don't need to settle here
	if _, err := runWithoutContext("udevadm", "settle"); err != nil {
		t.Log(err)
	}
	return nil
}

func formatUUID(s string) string {
	return strings.ToUpper(strings.Replace(s, "-", "", -1))
}

func validateFilesystems(t *testing.T, expected []*types.Partition) error {
	for _, e := range expected {
		if e.FilesystemType != "" {
			filesystemType, err := util.FilesystemType(e.Device)
			if err != nil {
				return fmt.Errorf("couldn't determine filesystem type: %v", err)
			}
			if filesystemType != e.FilesystemType {
				t.Errorf("FilesystemType does not match, expected:%q actual:%q",
					e.FilesystemType, filesystemType)
			}
		}
		if e.FilesystemUUID != "" {
			filesystemUUID, err := util.FilesystemUUID(e.Device)
			if err != nil {
				return fmt.Errorf("couldn't determine filesystem uuid: %v", err)
			}
			if formatUUID(filesystemUUID) != formatUUID(e.FilesystemUUID) {
				t.Errorf("FilesystemUUID does not match, expected:%q actual:%q",
					e.FilesystemUUID, filesystemUUID)
			}
		}
		if e.FilesystemLabel != "" {
			filesystemLabel, err := util.FilesystemLabel(e.Device)
			if err != nil {
				return fmt.Errorf("couldn't determine filesystem label: %v", err)
			}
			if filesystemLabel != e.FilesystemLabel {
				t.Errorf("FilesystemLabel does not match, expected:%q actual:%q",
					e.FilesystemLabel, filesystemLabel)
			}
		}
	}
	return nil
}

func validatePartitionNodes(t *testing.T, ctx context.Context, partition *types.Partition) {
	if err := mountPartition(ctx, partition); err != nil {
		t.Errorf("failed to mount %s: %v", partition.Device, err)
	}
	defer func() {
		if err := umountPartition(partition); err != nil {
			// failing to unmount is not a validation failure
			t.Logf("Failed to unmount %s: %v", partition.MountPath, err)
		}
	}()
	for _, file := range partition.Files {
		validateFile(t, partition, file)
	}
	for _, dir := range partition.Directories {
		validateDirectory(t, partition, dir)
	}
	for _, link := range partition.Links {
		validateLink(t, partition, link)
	}
	for _, node := range partition.RemovedNodes {
		path := filepath.Join(partition.MountPath, node.Directory, node.Name)
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			t.Error("Node was expected to be removed and is present!", path)
		}
	}
}

func validateFilesDirectoriesAndLinks(t *testing.T, ctx context.Context, expected []*types.Partition) {
	for _, partition := range expected {
		if partition.TypeCode == "blank" || partition.Length == 0 || partition.FilesystemType == "" || partition.FilesystemType == "swap" {
			continue
		}
		validatePartitionNodes(t, ctx, partition)
	}
}

func validateFile(t *testing.T, partition *types.Partition, file types.File) {
	path := filepath.Join(partition.MountPath, file.Node.Directory, file.Node.Name)
	fileInfo, err := os.Stat(path)
	if err != nil {
		t.Errorf("Error stat'ing file %s: %v", path, err)
		return
	}
	if file.Contents != "" {
		dat, err := ioutil.ReadFile(path)
		if err != nil {
			t.Error("Error when reading file", path)
			return
		}

		actualContents := string(dat)
		if file.Contents != actualContents {
			t.Error("Contents of file", path, "do not match!",
				file.Contents, actualContents)
		}
	}

	validateMode(t, path, file.Mode)
	validateNode(t, fileInfo, file.Node)
}

func validateDirectory(t *testing.T, partition *types.Partition, dir types.Directory) {
	path := filepath.Join(partition.MountPath, dir.Node.Directory, dir.Node.Name)
	dirInfo, err := os.Stat(path)
	if err != nil {
		t.Errorf("Error stat'ing directory %s: %v", path, err)
		return
	}
	if !dirInfo.IsDir() {
		t.Errorf("Node at %s is not a directory!", path)
	}
	validateMode(t, path, dir.Mode)
	validateNode(t, dirInfo, dir.Node)
}

func validateLink(t *testing.T, partition *types.Partition, link types.Link) {
	linkPath := filepath.Join(partition.MountPath, link.Node.Directory, link.Node.Name)
	linkInfo, err := os.Lstat(linkPath)
	if err != nil {
		t.Error("Error stat'ing link \"" + linkPath + "\": " + err.Error())
		return
	}
	if link.Hard {
		targetPath := filepath.Join(partition.MountPath, link.Target)
		targetInfo, err := os.Stat(targetPath)
		if err != nil {
			t.Error("Error stat'ing target \"" + targetPath + "\": " + err.Error())
			return
		}
		if linkInfoSys, ok := linkInfo.Sys().(*syscall.Stat_t); ok {
			if targetInfoSys, ok := targetInfo.Sys().(*syscall.Stat_t); ok {
				if linkInfoSys.Ino != targetInfoSys.Ino {
					t.Error("Hard link and target don't have same inode value: " + linkPath + " " + targetPath)
					return
				}
			} else {
				t.Error("Stat type assertion failed, this will only work on Linux")
				return
			}
		} else {
			t.Error("Stat type assertion failed, this will only work on Linux")
			return
		}
	} else {
		if linkInfo.Mode()&os.ModeType != os.ModeSymlink {
			t.Errorf("Node at symlink path is not a symlink (it's a %s): %s", linkInfo.Mode().String(), linkPath)
			return
		}
		targetPath, err := os.Readlink(linkPath)
		if err != nil {
			t.Error("Error reading symbolic link: " + err.Error())
			return
		}
		if targetPath != link.Target {
			t.Errorf("Actual and expected symbolic link targets don't match. Expected %q, got %q", link.Target, targetPath)
			return
		}
	}
	validateNode(t, linkInfo, link.Node)
}

func validateMode(t *testing.T, path string, mode int) {
	if mode != 0 {
		fileInfo, err := os.Stat(path)
		if err != nil {
			t.Error("Error running stat on node", path, err)
			return
		}

		if fileInfo.Mode() != os.FileMode(mode) {
			t.Error("Node Mode does not match", path, mode, fileInfo.Mode())
		}
	}
}

func validateNode(t *testing.T, nodeInfo os.FileInfo, node types.Node) {
	if nodeInfoSys, ok := nodeInfo.Sys().(*syscall.Stat_t); ok {
		if nodeInfoSys.Uid != uint32(node.User) {
			t.Error("Node has the wrong owner", node.User, nodeInfoSys.Uid)
		}

		if nodeInfoSys.Gid != uint32(node.Group) {
			t.Error("Node has the wrong group owner", node.Group, nodeInfoSys.Gid)
		}
	} else {
		t.Error("Stat type assertion failed, this will only work on Linux")
		return
	}
}
