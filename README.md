## README

goFAL is a File Abstraction Layer written in Go.

__Purpose__

I've found myself scripting a fair bit in Go and needed this library. So built it.

__Idea__

Create a file structure in memory. This file structure contains directories and files. Each directory and file has permissions. Once you've built the structure, generate it with the `Generate()` function using the root directory as the anchor point.

Next, you can populate the files with data using helper functions.

Once happy, run the `BuildHashes()` function, again using the root directory as the anchor.

Other stages like `SetPerms()` will set the permissions to what you intend, else they remain at `0777` in octal which is basically a free for all.

I'll add `BuildSigs()` soon allowing you to generate all of the sigs or partial, depending on how many files there are.

__Example__

```go
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

	logger.Info("--- Welcome to the goAFL Demo ---")

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

```

The actual stdout from this demo looks like this:

```bash
cd exampleCode
go build
./exampleCode
INFO[0000] --- Welcome to the FDS Demo ---
 Root Name:       build
 Is Directory:    true
 Permissions:     -rwxrwxrwx
 Directory:       /Users/davidgee/Documents/GoDev/src/github.com/davidjohngee/goFAL/exampleCode
 Full path:       /Users/davidgee/Documents/GoDev/src/github.com/davidjohngee/goFAL/exampleCode/build
--- Root Name:       content
--- Is Directory:    true
--- Permissions:     -rwxrwxrwx
--- Directory:       /Users/davidgee/Documents/GoDev/src/github.com/davidjohngee/goFAL/exampleCode/build
--- Full path:       /Users/davidgee/Documents/GoDev/src/github.com/davidjohngee/goFAL/exampleCode/build/content
------ Root Name:       content2.txt
------ Is Directory:    false
------ Permissions:     -r--r--r--
------ Directory:       /Users/davidgee/Documents/GoDev/src/github.com/davidjohngee/goFAL/exampleCode/build/content
------ Full path:       /Users/davidgee/Documents/GoDev/src/github.com/davidjohngee/goFAL/exampleCode/build/content/content2.txt
------ SHA1:            f1a9c7040d7e12388ef8cff7c46a775f817e023a
------ SHA256:          5f99c18e0aacfcaba41a170bfcf72bed18ef922c3625c8db5f60c9662ff3b71a
--- Root Name:       content1.txt
--- Is Directory:    false
--- Permissions:     -r--r--r--
--- Directory:       /Users/davidgee/Documents/GoDev/src/github.com/davidjohngee/goFAL/exampleCode/build
--- Full path:       /Users/davidgee/Documents/GoDev/src/github.com/davidjohngee/goFAL/exampleCode/build/content1.txt
--- SHA1:            b06ffa9228f909b137f05b722134587e4cac50c7
--- SHA256:          b2eb802149ad810cc9db94b1648224633c575ee55b682f8e683330adb7e96b15
```

__Result__

Running some simple stat and tree commands.

```bash
.
├── build
│   ├── content
│   │   └── content2.txt
│   └── content1.txt
```

Ok, so that worked. What about the permissions?

```bash
stat -f '%A %N' *
# 755 build
stat -f '%A %N' build/*
# 777 build/content
# 444 build/content1.txt
stat -f '%A %N' build/content/*
# 444 build/content/content2.txt
```

__Useful?__

Hopefully for someone. If not, I'll continue to use it.

__Contributors__

Fork it, create a PR to main. Avoid branches.
