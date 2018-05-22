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

package gofal

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"os"
	"sync"
	"text/tabwriter"
)

/* ------------------------------------ Const -------------------------------- */

const (
	// NONE is the root layer constant
	NONE = int(0)
	// TREE is another 0 for eye ball sugar
	TREE = int(1)
	// SHA1 constant
	SHA1 = 1
	// SHA256 constant
	SHA256 = 2
	// DIR is directory
	DIR = 1
	// FILE is file
	FILE = 2
)

/* ------------------------------------ Types -------------------------------- */

// FD which can have instances of itself
type FD struct {
	Name  string      // Name of FD
	H1    []byte      // SHA1 hash
	H256  []byte      // SHA256 hash
	FDs   []*FD       // This is where things get interesting. If []F is empty, then it's a file, else directory
	IsDir bool        // True for directory, false for file
	Perm  os.FileMode // Permission on file/directory
	Dir   string      // Current location. This ultimately is an abstraction.
	Loc   string      // This is our full rooted path string
}

/* ------------------------------------ Methods -------------------------------- */

// BuildRoot function creates the root, does some init and returns the root FD{}
func BuildRoot(name string, perm os.FileMode) (*FD, error) {

	// A workaround for assignment annoyance
	var err error
	err = nil

	// Build root node
	fd := FD{}

	// Set Name
	fd.Name = name

	// Set Permissions
	fd.Perm = perm

	// Set IsDir to true, we are the root directory
	fd.IsDir = true

	// Set Dir
	fd.Dir, err = os.Getwd()
	fd.Loc = fd.Dir + "/" + fd.Name
	// fd.Loc = "./" + fd.Name
	if err != nil {
		return nil, err
	}

	// default return
	return &fd, err
}

// BuildNode creates an FD and populates it with stuffs
func BuildNode(root *FD, name string, perm os.FileMode, kind int) (*FD, error) {

	// A workaround for assignment annoyance
	var err error
	err = nil

	// Build root node
	fd := FD{}

	// Set Name
	fd.Name = name

	// Set Permissions
	fd.Perm = perm

	// if kind is DIR, then true, else false
	if kind == DIR {
		fd.IsDir = true
	}
	if kind == FILE {
		fd.IsDir = false
	}

	// Set Dir from the root parent
	fd.Dir = root.Loc
	fd.Loc = fd.Dir + "/" + fd.Name

	// Add node to parent
	root.FDs = append(root.FDs, &fd)

	// default return
	return &fd, err

}

// String is our human friendly checker
func (fd *FD) String() string {
	lvl := new(int)
	*lvl = 0
	return string(tstring(fd, lvl))
}

// TreeString is our human friendly checker
// It doesn't do things in order (file vs directory) so the tree looks strange
// TODO: Make this more like the app 'tree' at some point
// It's good enough today for the purposes
func (fd *FD) TreeString() string {

	var b []byte

	lvl := new(int)
	*lvl = 0

	buildStringRecursively(fd, &b, lvl)
	return string(b)
}

func buildStringRecursively(fd *FD, b *[]byte, lvl *int) {

	bs := tstring(fd, lvl)
	var tmp []byte
	tmp = *b
	tmp = append(tmp, bs...)
	*b = tmp

	if fd.IsDir == true {
		*lvl++
		fds := fd.FDs
		for _, f := range fds {
			buildStringRecursively(f, b, lvl)
		}
		*lvl--
	}
}

// TreeString is our human friendly checker
func tstring(fd *FD, layer *int) []byte {
	const padding = 3

	// Create our return data instance
	var returnData bytes.Buffer
	// Create a tabular writer for the return data
	w := tabwriter.NewWriter(&returnData, 0, 0, padding, ' ', 0)

	// Now we build an offset mechanism
	var offset []byte
	offset = make([]byte, *layer*3)
	for i := 0; i < *layer; i++ {
		copy(offset[i*3:], "---")
	}
	// Convert the offset to a string
	stroffset := string(offset)

	// Create each row, with adequate offset
	fmt.Fprintln(w, stroffset, "Root Name:\t", fd.Name)
	fmt.Fprintln(w, stroffset, "Is Directory:\t", fd.IsDir)
	fmt.Fprintln(w, stroffset, "Permissions:\t", fd.Perm)
	fmt.Fprintln(w, stroffset, "Directory:\t", fd.Dir)
	fmt.Fprintln(w, stroffset, "Full path:\t", fd.Loc)
	// If dir, we're not interested in dirs actual hashes
	if fd.IsDir == false {
		fmt.Fprintln(w, stroffset, "SHA1:\t", fmt.Sprintf("%x", fd.H1))
		fmt.Fprintln(w, stroffset, "SHA256:\t", fmt.Sprintf("%x", fd.H256))
	}
	w.Flush()

	// Return the string!
	return returnData.Bytes()
}

// calcHashes is our dirty goRoutine
func calcHashes(t int, r *[]byte, w *sync.WaitGroup, f string, e chan error) {
	// This is a closure...we are safe! As a func call, we have everything on the stack frame

	// In case this is no longer the case, uncomment this line and the other terminator below.
	// go func(t int, r *[]byte, w *sync.WaitGroup, f string) {
	go func(e chan error) {
		var h hash.Hash
		*r = make([]byte, 0)

		switch t {
		case SHA1:
			h = sha1.New()
		case SHA256:
			h = sha256.New()
		default:
		}

		fh, err := os.Open(f)
		if err != nil {
			e <- err
			os.Exit(1)
		}

		if _, err := io.Copy(h, fh); err != nil {
			e <- err
			os.Exit(1)
		}
		fh.Close()

		*r = h.Sum(nil)
		w.Done()
		//}(t, r, w, f)
		e <- nil
	}(e)
}

// BuildHashes creates SHA1 and SHA256 hashes of the file contents (providing they exist)
// This is based on recusion. *WARNING - untested. BE CAREFUL WITH LARGE FILE STRUCTURES*
func BuildHashes(fd *FD) error {
	// 1.	If we're a directory, get a slice of *FDs and call BuildHashes(fd)
	// 2.	If we're not a directory, calculate hashes

	// 1.
	if fd.IsDir == true {
		var fds []*FD
		fds = fd.FDs
		for _, f := range fds {
			// Let's go inception
			BuildHashes(f)
		}
	}

	// 2.
	if fd.IsDir == false {
		// We only need wg for this logicical branch
		var wg sync.WaitGroup
		// Split each job in to 2x Go routines (GR in function itself to avoid stupid gotcha)
		e := make(chan error, 2)

		wg.Add(1)
		calcHashes(SHA1, &fd.H1, &wg, fd.Loc, e)
		wg.Add(1)
		calcHashes(SHA256, &fd.H256, &wg, fd.Loc, e)
		// Wait for the routines to exit
		wg.Wait()

		// We should have two nils in our channel
		for i := 0; i < 1; i++ {
			tmp := <-e
			if tmp != nil {
				return tmp
			}
		}
	}

	return nil
}

// Generate builds the file structures
func Generate(fd *FD) error {
	// 1. If we're a directory, build it and set the permissions
	// then grap the slice of FDs and recurse
	//
	// 2. If we're a file, create it and set the permissions
	// do not grab slice of FDs

	// 1.
	if fd.IsDir == true {
		var fds []*FD
		fds = fd.FDs
		// Create dir and build it. Note permissions here are loose.
		// Perms will be changed by SetPerms.
		err := os.Mkdir(fd.Loc, 0777)
		// If it exists, meh
		if err != nil {
			return err
		}

		for _, f := range fds {
			// Let's go inception
			e := Generate(f)
			if e != nil {
				return e
			}
		}
	}

	// We're not a directory, build it!
	if fd.IsDir == false {
		_, err := os.Create(fd.Loc)
		// If it exists, meh
		if err != nil {
			return err
		}

		err = os.Chmod(fd.Loc, 0777)
		if err != nil {
			return err
		}
		// err = os.Chmod(fd.Loc, fd.Perm)
	}
	return nil
}

// SetPerms sets the permissions when we're done building stuff
// Generate builds the file structures
func SetPerms(fd *FD) error {
	// 1. If we're a directory, build it and set the permissions
	// then grab the slice of FDs and recurse
	//
	// 2. If we're a file, create it and set the permissions
	// do not grab slice of FDs

	// 1.
	if fd.IsDir == true {
		var fds []*FD
		fds = fd.FDs

		for _, f := range fds {
			// Let's go inception
			err := os.Chmod(f.Loc, f.Perm)
			if err != nil {
				return err
			}
			// Recurse
			SetPerms(f)

		}
	}

	// We're not a directory, build it!
	if fd.IsDir == false {

		err := os.Chmod(fd.Loc, fd.Perm)
		if err != nil {
			return err
		}

	}

	return nil
}

// FileWrite allows us to write to the file.
// If this gets more advanced, move it to an interface for io implementation (Read/Write etc)
func FileWrite(f *FD, c []byte) error {
	filehandler, err := os.OpenFile(f.Loc, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	filehandler.Write(c)
	filehandler.Close()
	return nil
}

// SignHashes signs each hash on every file. Simples. Beware, could cost many CPU cycles.
func SignHashes() {
	// TODO:
}
