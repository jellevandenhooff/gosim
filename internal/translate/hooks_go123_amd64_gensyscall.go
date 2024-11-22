// Code generated by gensyscall. DO NOT EDIT.
package translate

func init() {
	hooksGensyscallGo123ByArch["amd64"] = map[packageSelector]packageSelector{
		{Pkg: "golang.org/x/sys/unix", Selector: "accept4"}:     {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "bind"}:        {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Close"}:       {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "connect"}:     {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Fdatasync"}:   {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Flock"}:       {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Fstat"}:       {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Fstatat"}:     {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Fsync"}:       {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Ftruncate"}:   {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Getdents"}:    {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "getpeername"}: {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Getpid"}:      {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Getrandom"}:   {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "getsockname"}: {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "getsockopt"}:  {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Listen"}:      {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Madvise"}:     {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "mmap"}:        {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "munmap"}:      {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "openat"}:      {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "pread"}:       {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "pwrite"}:      {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "read"}:        {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Renameat"}:    {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "setsockopt"}:  {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "socket"}:      {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Uname"}:       {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "Unlinkat"}:    {Pkg: hooksGo123Package},
		{Pkg: "golang.org/x/sys/unix", Selector: "write"}:       {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "accept4"}:                   {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "bind"}:                      {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Close"}:                     {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "connect"}:                   {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "fcntl"}:                     {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Fdatasync"}:                 {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Flock"}:                     {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Fstat"}:                     {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "fstatat"}:                   {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Fsync"}:                     {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Ftruncate"}:                 {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Getdents"}:                  {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "getpeername"}:               {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Getpid"}:                    {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "getsockname"}:               {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "getsockopt"}:                {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Listen"}:                    {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Madvise"}:                   {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "mmap"}:                      {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "munmap"}:                    {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "openat"}:                    {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "pread"}:                     {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "pwrite"}:                    {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "read"}:                      {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Renameat"}:                  {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "setsockopt"}:                {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "socket"}:                    {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "Uname"}:                     {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "unlinkat"}:                  {Pkg: hooksGo123Package},
		{Pkg: "syscall", Selector: "write"}:                     {Pkg: hooksGo123Package},
	}
}
