package utils

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/Sirupsen/logrus"
)

var (
	IMG_SUFFIXS   = []string{".jpg", ".bmp", ".jpeg", ".png", ".gif"}
	VIDEO_SUFFIXS = []string{".avi", ".rmvb", ".rm", ".asf", ".divx", ".mpg", ".mpeg", ".mpe", ".wmv", ".mp4", ".mkv", ".vob", ".h264"}
	TXT_SUFFIXS = []string{".txt"}
)

// return all images under basepath
func LsImages(basepath string) []string {

	return listFile(basepath, IMG_SUFFIXS)
}

// return all videos under basepath
func LsVideos(basepath string) []string {

	return listFile(basepath, VIDEO_SUFFIXS)
}

// recursive function, return a slice of full path of files under basepath and all its sub directory.
func listFile(basepath string, suffixs []string) (result []string) {

	files, _ := ioutil.ReadDir(basepath)

	for _, file := range files {

		//generate full path
		filepath := basepath + "/" + file.Name()

		if file.IsDir() {
			result = append(result, listFile(filepath, suffixs)...)
		} else {
			if suffixs == nil || StrInSliceIgnoreCase(path.Ext(filepath), suffixs) {

				result = append(result, filepath)
			}
		}
	}

	return result
}

//list files that last modified time > modetime in UNIXTIME
func ListFileFromTime(basepath string, suffixs []string, modtime int64) (result []string) {

	files, _ := ioutil.ReadDir(basepath)

	for _, file := range files {

		//generate full path
		filepath := basepath + "/" + file.Name()

		if file.IsDir() {
			result = append(result, listFile(filepath, suffixs)...)
		} else {
			if suffixs == nil || StrInSliceIgnoreCase(path.Ext(filepath), suffixs) {

				info, _ := os.Stat(filepath)

				if info.ModTime().Unix() > modtime {
					result = append(result, filepath)
				}
			}
		}
	}

	return result
}

// similar as listFile, but return a silce of os.FileInfo
func listFileInfo(basepath string, suffixs []string) (result []os.FileInfo) {

	files, _ := ioutil.ReadDir(basepath)

	for _, file := range files {

		//generate full path
		filepath := basepath + "/" + file.Name()

		if file.IsDir() {
			listFile(filepath, suffixs)
		} else {
			if suffixs == nil || StrInSliceIgnoreCase(path.Ext(filepath), suffixs) {

				info, _ := os.Stat(filepath)

				result = append(result, info)
			}
		}
	}

	return result

}

//execute command
func ExecCommand(input string) (output string, errput string, err error) {

	return ExecCommandinDir("", input)
}

//execute command in a certain directory
func ExecCommandinDir(workDir string, input string) (output string, errput string, err error) {
	var retoutput string
	var reterrput string
	cmd := exec.Command("/bin/bash", "-c", input)

	//workDir is not ""
	if !strings.EqualFold(workDir, "") {
		cmd.Dir = workDir
	}

	logrus.Debugf("execute local command [%v]", cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Errorf("init stdout failed, error is %v", err)
		return "", "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		logrus.Errorf("init stderr failed, error is %v", err)
		return "", "", err
	}

	if err := cmd.Start(); err != nil {
		logrus.Errorf("start command failed, error is %v", err)
		return "", "", err
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		logrus.Errorf("read stderr failed, error is %v", err)
		return "", "", err
	}

	if len(bytesErr) != 0 {
		reterrput = strings.Trim(string(bytesErr), "\n")
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		logrus.Errorf("read stdout failed, error is %v", err)
		return "", reterrput, err
	}

	if len(bytes) != 0 {
		retoutput = strings.Trim(string(bytes), "\n")
	}

	if err := cmd.Wait(); err != nil {
		logrus.Errorf("wait command failed, error is %v", err)
		logrus.Errorf("reterrput is %s", reterrput)
		return retoutput, reterrput, err
	}

	logrus.Debugf("retouput is %s", retoutput)
	logrus.Debugf("reterrput is %s", reterrput)
	return retoutput, reterrput, err
}
