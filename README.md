# goublu

Goublu is a [Go language](http://golang.org) front end that provides a better console interface to [Ublu](https://github.com/jwoehr/ublu) than the console support provided by Java.

![goublu_screenshot](https://user-images.githubusercontent.com/4604036/28322382-317d05fa-6b93-11e7-8457-b07eec2873af.png)

- Works on
  - OpenBSD
  - Linux
  - Mac OS X
- Works poorly on
  - Windows

Report bugs or make feature requests in the [Issue Tracker](https://github.com/jwoehr/goublu/issues)

## Building

### Prerequisites

- Go 1.17 or later
- Make (for Unix-like systems)
- Git (for automatic version detection)

### Quick Start

```bash
# Clone the repository
git clone https://github.com/jwoehr/goublu.git
cd goublu

# Build using Make (recommended)
make build

# Or install directly to $GOPATH/bin
make install
```

### Build Methods

#### Using Make (Recommended)

The project includes a Makefile with automatic version detection from git tags:

```bash
make build         # Build the binary
make install       # Install to $GOPATH/bin
make clean         # Remove build artifacts
make test          # Run tests
make test-coverage # Run tests with coverage report
make all           # Run fmt, vet, test, and build
make help          # Show all available targets
```

#### Using the Legacy Shell Script

For compatibility, the old build script is still available:

```bash
./make.sh
```

**Note:** The Makefile is preferred as it provides automatic version detection and additional development targets.

#### Manual Build

If you prefer to build manually without version information:

```bash
go build .
```

**Note:** Manual builds won't include version and compile date information.

### Alternative Installation Methods

#### Using go install

```bash
go install github.com/jwoehr/goublu@latest
```

#### From Release

Download a pre-built binary from the [releases page](https://github.com/jwoehr/goublu/releases).

## Invoking

- Invoke: `./goublu [-v] [-g "GoubluOpt1=SomeThing:GoubluOpt2=Other:..."] ublu_arg ublu_arg ...`
  - If the first argument to goublu is `-v` then Goublu prints a version message and exits 0.
  - If the first argument to goublu is `-g` then the next element in the command line is assumed
  to be a string of Goublu property-like options of the form Opt=Value, each option separated from
  the next by `:` . All remaining commandline arguments are passed to Ublu. The Goublu options and their
  values are case-sensitive and are as follows:
    - `UbluDir`
      - abs path to dir where ublu.jar resides, default `/opt/ublu`
    - `JavaOpt`
      - any option to the Java runtime, e.g, `JavaOpt=-Dsomething=other` (one option per JavaOpt line)
    - `SaveOutDir`
      - abs path to where pressing F4 saves the output text, default `/tmp`
    - `PropsFile`
      - abs path to a properties file containing these same `option=value` pairs
    - `BgColorIn`
      - Input background color, one of:
        - `ColorBlack`
        - `ColorRed`
        - `ColorGreen`
        - `ColorYellow`
        - `ColorBlue`
        - `ColorMagenta`
        - `ColorCyan`
        - `ColorWhite`
        - `ColorDefault` (default terminal colors)
    - `FgColorIn`
      - Input foreground color, as above
    - `BgColorOut`
      - Output background color, as above
    - `FgColorOut`
      - Output foreground color, as above
    - `Macro=name freeform string of Ublu commands`
      - Sets macro `name` to `freeform string of Ublu commands`
- Assumes in absence of property set as above that Ublu is found in `/opt/ublu/ublu.jar`

## Example setup

In `.bash_aliases` (or `.bashrc` or whatever):

```bash
alias gu='/home/jax/gopath/src/github.com/jwoehr/goublu/main/goublu -g PropsFile=/home/jax/.config/ublu/goublu.properties $*'
```

In the `PropsFile`:

```text
#BgColorOut=ColorBlack
FgColorOut=ColorRed
UbluDir=/opt/ublu
SaveOutDir=.
JavaOpt=-Djavax.net.ssl.trustStore=/opt/ublu/keystore/ublutruststore
JavaOpt=-Dublu.includepath=/opt/ublu/examples:/opt/ublu/extensions
JavaOpt=-Dublu.usage.linelength=100
Macro=sys1 as400 -to @sys1 SYS1.FOO.COM myusrprf
Macro=in include
Macro=jl joblist -as400
Macro=db db -to @myDb -dbtype as400 -connect
Macro=ublutest /QSYS.LIB/UBLUTEST.LIB/
Macro=spfl spoolflist -as400
Macro=ul userlist -as400
Macro=ref desktop -browse file:///opt/ublu/userdoc/ubluref.html#
```

## Working in Goublu

- Basic line editing
  - Ctl-a move to head of line
  - Ctl-b move one back.
  - Ctl-e move to end of line.
  - Ctl-d delete to end of current word.
  - Ctl-f move one ahead.
  - Ctl-k delete to end of line.
    - This doesn't work entirely right if line is longer than view width.
  - Alt-b back a word.
  - Alt-f forward a word.
  - These work as you would expect:
    - Home
    - End
    - Backspace
    - Left-arrow
    - Right-arrow
    - Insert
    - Delete
- History
  - Up-arrow previous command
  - Down-arrow next command
  - PgUp first command
  - PgDn last command
- F1 shows Goublu help.
- F2 shows entire session's output
- F3 offers a quick exit for when Ublu gets caught in a loop or network timeout
- F4 saves the entire session's output to a file `SaveOutDir/goublu.out.`_xxx..._
  - SaveOutDir set above as Goublu property, default is /tmp
  - Output announces the save file name
  - You can do this as many times as you like during a session, a new file is created each time.
- F5 expands last element on command line as macro you set in the properties file.
  - On empty line, F5 lists Goublu version, compile date, and all Goublu options and macros.
- F9 rotates through previous commands wrapping.
- Ctrl-Space at the end of a partial command name rotates through completions, if any.

## Testing

The project includes comprehensive unit tests for core functionality:

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run specific tests
go test -v -run TestStartStreamReader

# Run benchmarks
go test -bench=. -benchmem
```

### Test Coverage

Current test coverage focuses on:

- Stream reader functionality with context cancellation
- Panic recovery in stream readers
- Concurrent stream reader operations
- GUI initialization error handling
- Context cancellation behavior

The test suite includes:

- **Unit tests**: Testing individual functions in isolation
- **Integration tests**: Testing component interactions
- **Concurrency tests**: Verifying thread-safe operations
- **Benchmarks**: Performance testing of critical paths

## Notes

- The Ublu prompt appears on a line by itself in Goublu.
- Goublu "history" is input line recall and is separate from Ublu's own `history` command.
- Any Ublu application program output should include a newline as the Goublu output mechanism requires it.
- This document as displayed on the project page always reflects the current state of the tree and may be in
advance of the release version.

## Recent Code Improvements (2026-05-19)

The main function has been significantly refactored to improve code quality, maintainability, and robustness:

### Refactoring Changes

1. **Extracted Stream Reader Logic**
   - Eliminated 35+ lines of duplicated code by creating a reusable `startStreamReader()` function
   - Both stdout and stderr now use the same underlying stream reading logic
   - Added context-based cancellation for graceful goroutine shutdown

2. **Visual Distinction for Error Output**
   - Added new `Ubluerr()` method to `UbluManager` for handling stderr
   - Error output now displays in red color to distinguish it from normal output
   - Improves debugging and error identification during Ublu sessions

3. **Extracted GUI Initialization**
   - Created `initializeGUI()` function to separate GUI setup from main logic
   - Improves code organization and testability
   - Makes the main function more readable and focused

4. **Improved Control Flow**
   - Version flag now uses early return pattern (idiomatic Go)
   - Nil checks converted to early returns with explicit error messages
   - Reduced nesting depth from 3 levels to 1 level in main function

5. **Enhanced Application Lifecycle Management**
   - Created `runApplication()` function to manage GUI and Ublu execution
   - Better error handling with proper cleanup guarantees
   - Improved error propagation from both GUI and Ublu processes

6. **Specific Exit Codes**
   - Exit code 0: Normal exit or version display
   - Exit code 1: Failed to initialize Ublu
   - Exit code 2: Failed to initialize GUI
   - Exit code 3: Application runtime error

7. **Refactored CommandLineEditor**
   - Extracted 80+ line editor function into 10 focused methods
   - Added panic recovery for all editor operations
   - Improved error handling for file save operations
   - Organized key handling by category (character input, navigation, editing, function keys)
   - Better separation of concerns and maintainability

### Benefits

- **Maintainability**: Reduced code duplication and improved separation of concerns
- **Reliability**: Context-based cancellation prevents goroutine leaks, panic recovery in editor
- **Debuggability**: Visual distinction for errors and specific exit codes
- **Readability**: Main function reduced from ~82 lines to ~45 lines, editor organized into logical methods
- **Testability**: Extracted functions can be tested independently
- **Error Handling**: Graceful error recovery in editor operations, no more panics on file save errors

### Unit Tests Added

Comprehensive test suite covering:

- Stream reader with multiple scenarios (empty input, multiple lines, context cancellation)
- Panic recovery in stream readers
- Concurrent stream reader operations
- GUI initialization with error handling
- Context cancellation verification
- Performance benchmarks

Test execution:

```bash
make test           # Run all tests
make test-coverage  # Generate coverage report (coverage.html)
```

Current coverage: 4.9% of statements (focused on newly refactored functions)

## Bugs

- Serious
  - Ublu prompts for a password when an AS400 object is created with an invalid password and does not echo. However,
  Goublu **will indeed echo the password** even though Ublu's password prompt says the password will not be echoed.
- Trivial
  - Command lines longer that the view width of the input line behave erratically in response to edit commands.
  - On Mac OS X in Terminal, mouse actions fill the input line with escape sequences and do not otherwise work.

## The default branch has been renamed to main

master is now named main

If you have a local clone, you can update it by running:

```bash
git branch -m master main
git fetch origin
git branch -u origin/main main
git remote set-head origin -a
```

Jack Woehr 2026-05-19
