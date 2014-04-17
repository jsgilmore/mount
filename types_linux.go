//   Copyright 2014 Vastech SA (PTY) LTD
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

// +build ignore

package mount

//#include <linux/magic.h>
import "C"

const (
	TMPFS_MAGIC       = C.TMPFS_MAGIC
	HUGETLBFS_MAGIC   = C.HUGETLBFS_MAGIC
	BTRFS_SUPER_MAGIC = C.BTRFS_SUPER_MAGIC
	EXT4_SUPER_MAGIC  = C.EXT4_SUPER_MAGIC
)
