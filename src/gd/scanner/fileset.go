/*
 * Copyright (C) 2023 The GDLang Team.
 *
 * This file is part of GDLang.
 *
 * GDLang is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * GDLang is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with GDLang.  If not, see <http://www.gnu.org/licenses/>.
 */

package scanner

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type FileSet struct {
	mutex sync.RWMutex         // protects the file set
	base  int                  // base offset for the next file
	files []*File              // list of files in the order added to the set
	last  atomic.Pointer[File] // cache of last file looked up
}

func (fs *FileSet) Base() int {
	fs.mutex.RLock()
	b := fs.base
	fs.mutex.RUnlock()
	return b
}

func (fs *FileSet) AddFile(filename string, base, size int) (*File, error) {
	// Allocate f outside the critical section.
	f := &File{name: filename, size: size, lines: []int{0}}

	fs.mutex.Lock()
	defer fs.mutex.Unlock()
	if base < 0 {
		base = fs.base
	}
	if base < fs.base {
		return nil, fmt.Errorf("invalid base %d (should be >= %d)", base, fs.base)
	}
	f.base = base
	if size < 0 {
		return nil, fmt.Errorf("invalid size %d (should be >= 0)", size)
	}
	// base >= s.base && size >= 0
	base += size + 1 // +1 because EOF also has a position
	if base < 0 {
		return nil, fmt.Errorf("offset overflow (> 2G of source code in file set)")
	}
	// add the file to the file set
	fs.base = base
	fs.files = append(fs.files, f)
	fs.last.Store(f)

	return f, nil
}

// Reset

func (fs *FileSet) Reset() {
	fs.mutex.Lock()
	fs.base = 1
	fs.files = nil
	fs.last.Store((*File)(nil))
	fs.mutex.Unlock()
}

func NFileSet() *FileSet {
	return &FileSet{
		base: 1, // 0 == NoPos
	}
}
