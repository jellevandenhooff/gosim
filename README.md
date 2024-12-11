# Gosim: Simulation testing for Go
[![Go Reference](https://pkg.go.dev/badge/github.com/jellevandenhooff/gosim.svg)](https://pkg.go.dev/github.com/jellevandenhooff/gosim)

Gosim is a simulation testing framework that aims to make it easier to write
reliable distributed systems code in go. Gosim takes standard go programs and
tests and runs them with a simulated network, filesystem, multiple machines, and
more.

In Gosim's simulated environment you can run an entire distributed system in a
single go process, without having to set up any real servers. The simulation can
also introduce all kinds of problems, such as flaky networks, restarting
machines, or lost writes on a filesystem, to test that the program survive those
problems before they happen in real life.

An interesting and entertaining introduction to simulation testing is the talk
[Testing Distributed Systems w/ Deterministic Simulation from Strangeloop
2014](https://www.youtube.com/watch?v=4fFDFbi3toc). You can also listen to a
[NotebookLM podcast](https://notebooklm.google.com/notebook/7b466821-e836-4836-9c46-eb280f385fcd/audio) of this repo for entertainment purposes.

Gosim is an experimental project. Feedback, suggestions, and ideas are all very
welcome. The underlying design feels quite solid, but the implemented and
simulated APIs are still limited: files and directories work, but no permissions
or links; TCP works, but not UDP; IP addresses work, but not hostnames. The API
can and will change. Now that those warnings are out of the way, please take a
look at what gosim can do.

# Using Gosim

## From go test to gosim test
Gosim tests standard go code and is used just like another go package. To get
started with `gosim`, you need to import it in your module:
```
> mkdir example && cd example
> go mod init example
> go get github.com/jellevandenhooff/gosim
```
To test code with Gosim, write a small a test in a file `simple_test.go` and then run
it using `go test`:
```go
package examples_test

import (
	"math/rand"
	"testing"

	"github.com/jellevandenhooff/gosim"
)

func TestGosim(t *testing.T) {
	t.Logf("Are we in the Matrix? %v", gosim.IsSim())
	t.Logf("Random: %d", rand.Int())
}
```
```
> go test -v -count=1 -run TestGosim
=== RUN   TestGosim
    simple_test.go:11: Are we in the Matrix? false
    simple_test.go:12: Random: 1225418128470362734
--- PASS: TestGosim (0.00s)
PASS
ok  	example	0.216s
```
The test prints that Gosim is not enabled, and each time the test runs the random number is different.

To run this test with `gosim` instead, replace `go test` with
`go run github.com/jellevandenhooff/gosim/cmd/gosim test`:
```
> go run github.com/jellevandenhooff/gosim/cmd/gosim test -v -run TestGosim
=== RUN   TestGosim (seed 1)
    1 main/4     14:10:03.000 INF examples/simple_test.go:12 > Are we in the Matrix? true method=t.Logf
    2 main/4     14:10:03.000 INF examples/simple_test.go:13 > Random: 811966193383742320 method=t.Logf
--- PASS: TestGosim (0.00s simulated)
ok  	translated/example	0.204s
```
The `gosim test` command has flags similar to `go test`. The test output is
more involved than a normal `go test`. Every log line includes the simulated
machine and the goroutine that invoked the log to help debug tests running on
multiple machines.

The `=== RUN` line includes the test's seed. Each time this test runs with the
same seed it will output the same random number. To test with different seeds,
you can pass ranges of seeds to `gosim test`:
```
go run github.com/jellevandenhooff/gosim/cmd/gosim test -v -seeds=1-3 -run=TestGosim .
=== RUN   TestGosim (seed 1)
    1 main/4     14:10:03.000 INF examples/simple_test.go:12 > Are we in the Matrix? true method=t.Logf
    2 main/4     14:10:03.000 INF examples/simple_test.go:13 > Random: 811966193383742320 method=t.Logf
--- PASS: TestGosim (0.00s simulated)
=== RUN   TestGosim (seed 2)
    1 main/4     14:10:03.000 INF examples/simple_test.go:12 > Are we in the Matrix? true method=t.Logf
    2 main/4     14:10:03.000 INF examples/simple_test.go:13 > Random: 5374891573232646577 method=t.Logf
--- PASS: TestGosim (0.00s simulated)
=== RUN   TestGosim (seed 3)
    1 main/4     14:10:03.000 INF examples/simple_test.go:12 > Are we in the Matrix? true method=t.Logf
    2 main/4     14:10:03.000 INF examples/simple_test.go:13 > Random: 3226404213937589817 method=t.Logf
--- PASS: TestGosim (0.00s simulated)
ok  	translated/github.com/jellevandenhooff/gosim/examples	0.254s
```

If running `gosim` fails with errors about missing `go.sum` entries, run
`go mod tidy` to update your `go.sum` file.

## Simulation
Tests running in Gosim run inside Gosim's simulation environment. In the
simulation, tests can create simulated machines that can talk to each other
over a simulated network, crash and restart machine, introduce latency,
and more. The following longer example shows some of those features:
```go
package example_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/netip"
	"testing"
	"time"

	"github.com/jellevandenhooff/gosim"
)

var count = 0

func server() {
	log.Printf("starting server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		count++
		log.Printf("got a request from %s", r.RemoteAddr)
		fmt.Fprintf(w, "hello from the server! request: %d", count)
	})
	http.ListenAndServe("10.0.0.1:80", nil)
}

func request() {
	log.Println("making a request")
	resp, err := http.Get("http://10.0.0.1/")
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(string(body))
}

func TestMachines(t *testing.T) {
	// run the server
	serverMachine := gosim.NewMachine(gosim.MachineConfig{
		Label:    "server",
		Addr:     netip.MustParseAddr("10.0.0.1"),
		MainFunc: server,
	})

	// let the server start
	time.Sleep(time.Second)

	// make some requests to see them work
	request()
	request()

	// restart the server
	log.Println("restarting the server")
	serverMachine.Crash()
	serverMachine.Restart()
	time.Sleep(time.Second)

	// make a new request to see a reset count
	request()

	// add some latency
	log.Println("adding latency")
	gosim.SetDelay("10.0.0.1", "11.0.0.1", time.Second)

	// make another request to see the latency
	request()
}
```
This example contains a HTTP server running on its own gosim machine, and a
client that makes several requests to the server. Then the client
messes with the server, restarting it and adding latency. Let's see
what happens this test is run:
<!-- TODO: two machines and a bug? -->
```
> go run github.com/jellevandenhooff/gosim/cmd/gosim test -v -run TestMachines
=== RUN   TestMachines (seed 1)
    1 server/5   14:10:03.000 INF example/machines_test.go:17 > starting server
    2 main/4     14:10:04.000 INF example/machines_test.go:27 > making a request
    3 server/8   14:10:04.000 INF example/machines_test.go:20 > got a request from 11.0.0.1:10000
    4 main/4     14:10:04.000 INF example/machines_test.go:41 > hello from the server! request: 1
    5 main/4     14:10:04.000 INF example/machines_test.go:27 > making a request
    6 server/8   14:10:04.000 INF example/machines_test.go:20 > got a request from 11.0.0.1:10000
    7 main/4     14:10:04.000 INF example/machines_test.go:41 > hello from the server! request: 2
    8 main/4     14:10:04.000 INF example/machines_test.go:60 > restarting the server
    9 server/13  14:10:04.000 INF example/machines_test.go:17 > starting server
   10 main/4     14:10:05.000 INF example/machines_test.go:27 > making a request
   11 server/16  14:10:05.000 INF example/machines_test.go:20 > got a request from 11.0.0.1:10001
   12 main/4     14:10:05.000 INF example/machines_test.go:41 > hello from the server! request: 1
   13 main/4     14:10:05.000 INF example/machines_test.go:69 > adding latency
   14 main/4     14:10:05.000 INF example/machines_test.go:27 > making a request
   15 server/16  14:10:06.000 INF example/machines_test.go:20 > got a request from 11.0.0.1:10001
   16 main/4     14:10:07.000 INF example/machines_test.go:41 > hello from the server! request: 2
--- PASS: TestMachines (4.00s simulated)
ok  	translated/example	0.237s
```
In the log, each line is annotated with the machine and goroutine that wrote
the log. When debugging code running on multiple machines, it is helpful to
know where logs are coming from.

The first thing that happens is the server starting at step 1. Then the client
makes a request at step 2, the server receives it, responds, and the client
prints the response. When the client makes another request the server this
time responds with a higher count at steps 6 and 7.

Afterwards, the client restarts the server at step 8 and the server prints that
is starting again at step 9. When the client makes a request after the restart
the server responds once again with a count of 1 at steps 11 and 12.  How did
that happen? After the server is restarted, the counter resets to zero. Each
machine has its own copy of global variables which get re-initialized when a
machine crashes.

Notice that the client address (after "got a request from:" in the logs) also
changes after the restart, and the goroutine handling the requests on the server
changes. The simulated TCP connection breaks when the machine crashes and the
`net/http` client makes a new connection afterwards.

After the restart, the client calls a Gosim API to add latency to the simulated network
at step 13. The earlier requests all got sent and received at the same timestamp,
such as on steps 10, 11, and 12 all at 14:10:05. The request after adding
latency is sent at 14:10:05 on step 14 and received a second later at 14:10:06
on step 15. This is the added latency in action.

Even though the test takes over 4 simulated seconds, running the tests takes
less than a second. Simulated time can be sped up when Gosim knows that all
goroutines are waiting for time advance. With simulated time, tests can use
realistic latency and timeouts without development time becoming slow.

## Debugging
Gosim's simulation is deterministic, which means that running a test twice
will result in the exact same behavior, from randomly generated numbers
to goroutine concurrency interleavings. That is useful because it means that
if we see an interesting log after running the program but do not fully
understand why that log printed a certain value, we can re-run the program
and see exactly what happened at the moment of printing.

The numbers at the start of each log are Gosim's step numbers. Take step 11 from
the log above, from the line "got a request ...". We can debug the state of the
program at the time the log was printed with `gosim debug`:
```
> go run github.com/jellevandenhooff/gosim/cmd/gosim debug -package=. -test=TestMachines -step=11
...
    26:	func (w gosimSlogHandler) Handle(ctx context.Context, r slog.Record) error {
    27:		r.AddAttrs(slog.Int("goroutine", gosimruntime.GetGoroutine()))
=>  28:		r.AddAttrs(slog.Int("step", gosimruntime.Step()))
    29:		return w.inner.Handle(ctx, r)
    30:	}
...
```
The `gosim debug` command runs the requested test inside of the Delve go
debugger and stops the test at the requested step. The debugger here puts us in
Gosim's `log/slog` handler which calls `gosimruntime.Step()` and includes the
step number.

Now we can debug the program as a standard Go program. To see what was happening
in the HTTP handler that called print, we run backtrace
```
(dlv) bt
 0  0x000000010129e7ac in translated/github.com/jellevandenhooff/gosim/internal_/simulation.gosimSlogHandler.Handle
    at ./github.com/jellevandenhooff/gosim/internal_/simulation/userspace.go:28
 1  0x00000001012a5e24 in translated/github.com/jellevandenhooff/gosim/internal_/simulation.(*gosimSlogHandler).Handle
    at <autogenerated>:1
 2  0x0000000101283d68 in translated/log/slog.(*handlerWriter).Write
    at ./log/slog/logger.go:99
 3  0x000000010124a600 in translated/log.(*Logger).output
    at ./log/log.go:242
 4  0x000000010124ae38 in translated/log.Printf
    at ./log/log.go:394
 5  0x000000010160e03c in translated/example_test.server.func1
    at ./example/machines_test.go:20
...
```
and then select the HTTP handler's stack frame that called `log.Printf`
```
(dlv) frame 5
> translated/github.com/jellevandenhooff/gosim/internal_/simulation.gosimSlogHandler.Handle() ./github.com/jellevandenhooff/gosim/internal_/simulation/userspace.go:28 (PC: 0x1029ce45c)
Frame 5: ./example/machines_test.go:19 (PC: 102d3decc)
    14:	)
    15:
    16:	func server() {
    17:		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    18:			G().count++
=>  19:			log.Printf("got a request from %s", r.RemoteAddr)
    20:			fmt.Fprintf(w, "hello from the server! request: %d", G().count)
    21:		})
    22:		http.ListenAndServe("10.0.0.1:80", nil)
    23:	}
    24:
```
Now we can inspect the state of memory to see details that the log line might have missed
```
(dlv) print r.RemoteAddr
"11.0.0.1:10001"
```

## Tracing system calls
Gosim's simulation is at the system call level. To see why a program behaved
the way it did, Gosim can output detailed logs of all system calls with the
`-simtrace=syscall` flag:
```
> go run github.com/jellevandenhooff/gosim/cmd/gosim test -v -simtrace=syscall -run TestMachines
=== RUN   TestMachines (seed 1)
    1 server/5   14:10:03.000 INF examples/machines_test.go:19 > starting server
    2 server/5   14:10:03.000 INF simulation/os_linux_.go:88 > call SysSocket flags=SOCK_STREAM|SOCK_NONBLOCK|SOCK_CLOEXEC net=AF_INET proto=IPPROTO_IP traceKind=syscall
    3 server/5   14:10:03.000 INF simulation/os_linux_.go:93 > ret  SysSocket err=null fd=5 relatedStep=2 traceKind=syscall
    4 server/5   14:10:03.000 INF simulation/os_linux_.go:88 > call SysOpenat dirfd=AT_FDCWD flags=O_RDONLY|O_CLOEXEC mode=0o0 path=/proc/sys/net/core/somaxconn traceKind=syscall
    5 server/5   14:10:03.000 INF simulation/os_linux_.go:93 > ret  SysOpenat err="no such file or directory" fd=0 relatedStep=4 traceKind=syscall
    6 server/5   14:10:03.000 INF simulation/os_linux_.go:88 > call SysSetsockopt traceKind=syscall
    7 server/5   14:10:03.000 INF simulation/os_linux_.go:93 > ret  SysSetsockopt relatedStep=6 traceKind=syscall
    8 server/5   14:10:03.000 INF simulation/os_linux_.go:88 > call SysBind addr=10.0.0.1:80 fd=5 traceKind=syscall
    9 server/5   14:10:03.000 INF simulation/os_linux_.go:93 > ret  SysBind err=null relatedStep=8 traceKind=syscall
...
   60 main/10    14:10:04.000 INF simulation/os_linux_.go:88 > call SysWrite fd=5 p="GET / HTTP/1.1\r\nHost: 10.0.0.1\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n" traceKind=syscall
   61 main/9     14:10:04.000 INF simulation/os_linux_.go:93 > ret  SysRead err="resource temporarily unavailable" n=0 p= relatedStep=59 traceKind=syscall
   62 main/10    14:10:04.000 INF simulation/os_linux_.go:93 > ret  SysWrite err=null n=89 relatedStep=60 traceKind=syscall
   63 server/8   14:10:04.000 INF simulation/os_linux_.go:88 > call SysRead fd=6 len(p)=4096 traceKind=syscall
   64 server/8   14:10:04.000 INF simulation/os_linux_.go:93 > ret  SysRead err=null n=89 p="GET / HTTP/1.1\r\nHost: 10.0.0.1\r\nUser-Agent: Go-http-client/1.1\r\nAccept-Encoding: gzip\r\n\r\n" relatedStep=63 traceKind=syscall
   65 server/8   14:10:04.000 INF examples/machines_test.go:22 > got a request from 11.0.0.1:10000
   66 server/8   14:10:04.000 INF simulation/os_linux_.go:88 > call SysWrite fd=6 p="HTTP/1.1 200 OK\r\nDate: Wed, 15 Jan 2020 14:10:04 GMT\r\nContent-Length: 33\r\nContent-Type: text/plain; charset=utf-8\r\n\r\nhello from the server! request: 1" traceKind=syscall
   67 server/11  14:10:04.000 INF simulation/os_linux_.go:88 > call SysRead fd=6 len(p)=1 traceKind=syscall
   68 server/8   14:10:04.000 INF simulation/os_linux_.go:93 > ret  SysWrite err=null n=150 relatedStep=66 traceKind=syscall
   69 server/11  14:10:04.000 INF simulation/os_linux_.go:93 > ret  SysRead err="resource temporarily unavailable" n=0 p= relatedStep=67 traceKind=syscall
   70 main/9     14:10:04.000 INF simulation/os_linux_.go:88 > call SysRead fd=5 len(p)=4096 traceKind=syscall
   71 main/9     14:10:04.000 INF simulation/os_linux_.go:93 > ret  SysRead err=null n=150 p="HTTP/1.1 200 OK\r\nDate: Wed, 15 Jan 2020 14:10:04 GMT\r\nContent-Length: 33\r\nContent-Type: text/plain; charset=utf-8\r\n\r\nhello from the server! request: 1" relatedStep=70 traceKind=syscall
   72 server/8   14:10:04.000 INF simulation/os_linux_.go:88 > call SysRead fd=6 len(p)=4096 traceKind=syscall
   73 main/4     14:10:04.000 INF examples/machines_test.go:43 > hello from the server! request: 1
```
Here you can see on step 60 the HTTP GET request being written by the client, and then on step 64 being read by the server.

## Examples

Gosim can run some sizeable real-world code. In this repository
there are examples for:

- [examples/bolt](./examples/bolt/): reading and writing a [Bolt](https://github.com/etcd-io/bbolt) database

- [examples/etcd](./examples/etcd/): running a 3-node [Etcd](https://github.com/etcd-io/etcd) cluster and partitioning the
  nodes

# API and documentation

A description of Gosim's architecture and design decisions is in
[docs/design.md](docs/design.md).

Gosim's public packages are:

- [github.com/jellevandenhooff/gosim](./): the core API for creating machines
  and manipulating the simulation environment

- [github.com/jellevandenhooff/gosim/cmd/gosim](./cmd/gosim): the CLI for
  running Gosim tests

- [github.com/jellevandenhooff/gosim/metatesting](./metatesting/): a package
  for running Gosim tests inside of normal go test.

- [github.com/jellevandenhooff/gosim/nemesis](./nemesis/): a package, still
  sparse, to introduce chaos into simulations

# Development

Gosim and its tools are tested using the `./test.sh` script, which invokes the
tests defined in `./Taskfile.yml`. These test the Gosim tooling as well as the
behavior of simulated code. Ideally tests verify that the simulation behaves
like reality: The tests for the filesystem in
[github.com/jellevandenhooff/gosim/internal/tests/behavior/disk_test.go](./internal/tests/behavior/disk_test.go)
are run in both simulated and non-simulated builds.

Gosim simulates Linux, but should work on macOS, Linux, and Windows (not tested
on Windows) running on either arm64 or amd64. To test that code works on Linux
amd64 and arm64, there are scripts `.ci/crossarch-tests/test-amd64.sh` and
`.ci/crossarch-tests/test-arm64.sh` that run the tests using Docker. A Github
Action runs the tests on Linux amd64.

