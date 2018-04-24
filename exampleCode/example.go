/*
Copyright 2018 David Gee

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

package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/arsonistgopher/glog"
	gf "github.com/arsonistgopher/gofal"
)

func checkerr(l glog.Logger, e error) {
	if e != nil {
		l.Error(e)
	}
}

func main() {

	/*
		Example of how to use the goFAL package.

		goFAL builds a tree of file type instances. This tree is then used
		to build a real file tree when the correct package functions are called.
		Each file has both SHA1 and SHA256 hashes calculated for easy use.

		This was created to make projects easier to build for packaging and scripting.

		This is an alpha release at best and comes without support.

		Author: David Gee
		Copyright: David Gee
		Date: 20th April 2018
		Contributors welcome!
	*/

	// De-coupled logging
	logger := glog.Logger{Name: "fds"}
	// Inject logrus
	logger.LoggingBase = logrus.New()

	logger.Info("--- Welcome to the FDS Demo ---")

	// Create root node called 'build'
	build, err := gf.BuildRoot("build", os.ModePerm)
	checkerr(logger, err)

	// Add a content directory called 'content' and set the permissions
	content, err := gf.BuildNode(build, "content", os.ModePerm, gf.DIR)
	checkerr(logger, err)

	// Add a file under the build directory called 'content1' and set perms
	file1, err := gf.BuildNode(build, "content1.txt", 0444, gf.FILE)
	checkerr(logger, err)

	// Add a file under the content directory called 'content2' and set perms
	file2, err := gf.BuildNode(content, "content2.txt", 0444, gf.FILE)
	checkerr(logger, err)

	// Generate file tree on disk
	err = gf.Generate(build)
	checkerr(logger, err)

	// Insert content
	file1Content := []byte("Hello from ArsonistGopher once.")
	file2Content := []byte("Hello from ArsonistGopher twice.")
	err = gf.FileWrite(file1, file1Content)
	checkerr(logger, err)
	err = gf.FileWrite(file2, file2Content)
	checkerr(logger, err)

	// Create hashes using the build root as an anchor
	err = gf.BuildHashes(build)
	checkerr(logger, err)

	// Set permissions (post writing) as per tree data
	err = gf.SetPerms(build)
	checkerr(logger, err)

	// Uncomment line below to print build info
	// fmt.Println(build.String())

	// Because we have a String() method, we can also call print directly. Uncomment line below.
	// fmt.Print(build)

	// This builds our stringfied object tree
	fmt.Println(build.TreeString())
}
