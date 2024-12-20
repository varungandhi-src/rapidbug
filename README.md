Reproduce issue with

go test ./... -rapid.seed=11647171348957545426 -rapid.nofailfile

Failure trace:

```
--- FAIL: TestCheckers (0.00s)
    --- FAIL: TestCheckers/property-based_tests (0.00s)
        checkers_test.go:40: [rapid] failed after 0 tests: (*T).FailNow() called
            To reproduce, specify -run="TestCheckers/property-based_tests" -rapid.seed=11647171348957545426
            Failed test output:
        checkers_test.go:44: 
            	Error Trace:	/Users/varun/Code/play/go/rapidbug/checkers_test.go:44
            	            				/Users/varun/go/pkg/mod/pgregory.net/rapid@v1.1.0/engine.go:368
            	            				/Users/varun/go/pkg/mod/pgregory.net/rapid@v1.1.0/engine.go:232
            	            				/Users/varun/go/pkg/mod/pgregory.net/rapid@v1.1.0/engine.go:118
            	            				/Users/varun/Code/play/go/rapidbug/checkers_test.go:40
            	Error:      	func (assert.PanicTestFunc)(0x1028501d0) should not panic
            	            		Panic value:	failed to find suitable value in 5 tries
            	            		Panic stack:	goroutine 36 [running]:
            	            	runtime/debug.Stack()
            	            		/Users/varun/.local/share/mise/installs/go/1.23.0/src/runtime/debug/stack.go:26 +0x64
            	            	github.com/stretchr/testify/assert.didPanic.func1()
            	            		/Users/varun/go/pkg/mod/github.com/stretchr/testify@v1.10.0/assert/assertions.go:1234 +0x74
            	            	panic({0x1029150e0?, 0x14000105680?})
            	            		/Users/varun/.local/share/mise/installs/go/1.23.0/src/runtime/panic.go:785 +0x124
            	            	pgregory.net/rapid.find[...](0x14000167588?, 0x1400036c180, 0x5)
            	            		/Users/varun/go/pkg/mod/pgregory.net/rapid@v1.1.0/combinators.go:118 +0x164
            	            	pgregory.net/rapid.(*customGen[...]).value(0x102850d94?, 0x33?)
            	            		/Users/varun/go/pkg/mod/pgregory.net/rapid@v1.1.0/combinators.go:36 +0x50
            	            	pgregory.net/rapid.(*Generator[...]).value(0x102672634, 0x1400036c180?)
            	            		/Users/varun/go/pkg/mod/pgregory.net/rapid@v1.1.0/generator.go:74 +0x74
            	            	pgregory.net/rapid.(*Generator[...]).Draw(0x102994c60?, 0x1400036c180, {0x102854bf9, 0x4})
            	            		/Users/varun/go/pkg/mod/pgregory.net/rapid@v1.1.0/generator.go:47 +0x70
            	            	github.com/varungandhi-src/rapidbug.checkerPBT(0x1400036c180, 0x1?)
            	            		/Users/varun/Code/play/go/rapidbug/checkers_test.go:52 +0x44
            	            	github.com/varungandhi-src/rapidbug.TestCheckers.func1.1.1()
            	            		/Users/varun/Code/play/go/rapidbug/checkers_test.go:45 +0x64
            	            	github.com/stretchr/testify/assert.didPanic(0x102bc5d90?)
            	            		/Users/varun/go/pkg/mod/github.com/stretchr/testify@v1.10.0/assert/assertions.go:1239 +0x74
            	            	github.com/stretchr/testify/assert.NotPanics({0x1497ba2e8, 0x1400036c180}, 0x14000105670, {0x0, 0x0, 0x0})
            	            		/Users/varun/go/pkg/mod/github.com/stretchr/testify@v1.10.0/assert/assertions.go:1310 +0x70
            	            	github.com/stretchr/testify/require.NotPanics({0x10298f7f8, 0x1400036c180}, 0x14000105670, {0x0, 0x0, 0x0})
            	            		/Users/varun/go/pkg/mod/github.com/stretchr/testify@v1.10.0/require/require.go:1669 +0xa4
            	            	github.com/varungandhi-src/rapidbug.TestCheckers.func1.1(0x1400036c180)
            	            		/Users/varun/Code/play/go/rapidbug/checkers_test.go:44 +0x74
            	            	pgregory.net/rapid.checkOnce(0x1400036c180, 0x10298b1f0?)
            	            		/Users/varun/go/pkg/mod/pgregory.net/rapid@v1.1.0/engine.go:368 +0x84
            	            	pgregory.net/rapid.checkTB({0x1029951d0, 0x1400011b1e0}, {0x1400011b040?, 0x14000122120?, 0x102bd6840?}, 0x10298b1f0)
            	            		/Users/varun/go/pkg/mod/pgregory.net/rapid@v1.1.0/engine.go:232 +0x908
            	            	pgregory.net/rapid.Check({0x102994bc8, 0x1400011b1e0}, 0x10298b1f0)
            	            		/Users/varun/go/pkg/mod/pgregory.net/rapid@v1.1.0/engine.go:118 +0x9c
            	            	github.com/varungandhi-src/rapidbug.TestCheckers.func1(0x1400011b1e0)
            	            		/Users/varun/Code/play/go/rapidbug/checkers_test.go:40 +0x38
            	            	testing.tRunner(0x1400011b1e0, 0x10298b1e8)
            	            		/Users/varun/.local/share/mise/installs/go/1.23.0/src/testing/testing.go:1690 +0xe4
            	            	created by testing.(*T).Run in goroutine 35
            	            		/Users/varun/.local/share/mise/installs/go/1.23.0/src/testing/testing.go:1743 +0x314
            	Test:       	TestCheckers/property-based_tests
FAIL
FAIL	github.com/varungandhi-src/rapidbug	0.325s
FAIL
```
