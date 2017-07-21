package main

import (
	"go/build"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"strings"

	"github.com/posener/complete"
)

// predictPackages completes packages in the directory pointed by a.Last
// and packages that are one level below that package.
func predictPackages(a complete.Args) (prediction []string) {
	// prediction = complete.PredictFilesSet(listPackages(a.Directory())).Predict(a)
	// if len(prediction) != 1 {
	// 	return
	// }
	// return complete.PredictFilesSet(listPackages(prediction[0])).Predict(a)

	gopathEntries := findGopath()
	// complete.Log("gopathEntries = %v", gopathEntries)
	ossep := string([]byte{os.PathSeparator})
	fragment := strings.Replace(a.Last, "/", ossep, -1)
	for _, gopath := range gopathEntries {
		gosrc := filepath.Join(gopath, "src")
		fragmentPath := filepath.Join(gosrc, fragment)
		// complete.Log("fragmentPath is %v", fragmentPath)
		dir := filepath.Dir(fragmentPath)
		// complete.Log("dir is %v", dir)
		infos, err := ioutil.ReadDir(dir)
		if err != nil {
			complete.Log("listing directory %s: %s", dir, err)
			continue
		}
		for _, info := range infos {
			// complete.Log("info.Name() is %v", info.Name())
			candidatePath := filepath.Join(dir, info.Name())
			if info.IsDir() && strings.HasPrefix(candidatePath, fragmentPath) {
				completionItem := candidatePath[len(gosrc):]
				completionItem = strings.Trim(completionItem, "/\\")
				if !strings.HasPrefix(completionItem, fragment) {
					complete.Log("ouch: %q doesn't start with %q", completionItem, fragment)
					continue
				}
				prediction = append(prediction, completionItem+"/", completionItem+"/...")
			}
		}
	}
	return
}

func findGopath() []string {
	gopath := os.Getenv("GOPATH")
	// complete.Log("gopath = %v", gopath)
	if gopath == "" {
		// By convention
		usr, err := user.Current()
		if err != nil {
			return nil
		}
		usrgo := filepath.Join(usr.HomeDir, "go")
		return []string{usrgo}
	}
	listsep := string([]byte{os.PathListSeparator})
	entries := strings.Split(gopath, listsep)
	return entries
}

// listPackages looks in current pointed dir and in all it's direct sub-packages
// and return a list of paths to go packages.
func listPackages(dir string) (directories []string) {
	// add subdirectories
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		complete.Log("failed reading directory %s: %s", dir, err)
		return
	}

	// build paths array
	paths := make([]string, 0, len(files)+1)
	for _, f := range files {
		if f.IsDir() {
			paths = append(paths, filepath.Join(dir, f.Name()))
		}
	}
	paths = append(paths, dir)

	// import packages according to given paths
	for _, p := range paths {
		pkg, err := build.Import("/home/valentin/go/src", p, 0)
		if err != nil {
			complete.Log("failed importing directory %s: %s", p, err)
			continue
		}
		directories = append(directories, pkg.Dir)
	}
	return
}
