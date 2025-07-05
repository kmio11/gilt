# Gilt
A library for performing Golden Tests in Go.

- Read and update golden files
- Use any assertion function
- Automatically clean unused golden files
- Customizable file handling with interfaces

## Basic Usage
> [!TIP]
> See the [golden_examples_test.go](./golden_examples_test.go) for more examples.

### Reading and writing Golden files as strings

Test code:
```go
func hello(to string) string {
	return "Hello, " + to + "!"
}

func TestNewStringGolden(t *testing.T) {
	golden := gilt.NewStringGolden(t.Name())

	tests := []struct {
		name string
	}{
		{"world"},
		{"gopher"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := hello(tt.name)
			golden.Assert(t, actual, tt.name, func(t *testing.T, actual string, expected string) {
				if actual != expected {
					t.Errorf("expected: %s\n, but was: %s", expected, actual)
				}
			})
		})
	}
}
```

### Updating and managing golden files

You can update the golden file by running the test with the `-update` flag.
```bash
$ go test -update
```

You can also clean unused golden files with the `-clean` flag:
```bash
$ go test -clean              # Remove unused golden files
$ go test -clean-dry-run      # Show what would be cleaned
$ go test -clean-backup       # Clean with backup creation
```

The following golden file will be generated:
```
testdata
└── testdata/TestNewStringGolden
    └── testdata/TestNewStringGolden/golden
        ├── testdata/TestNewStringGolden/golden/gopher.golden
        └── testdata/TestNewStringGolden/golden/world.golden
```

gopher.golden:
```
Hello, gopher!
```

Then you can run the test.
```bash
$ go test
```

### Reading and updating Golden files using json.Marshal/Unmarshal

Test code:
```go
type HelloMessage struct {
	Message string `json:"message"`
}

func helloMessage(to string) HelloMessage {
	return HelloMessage{"Hello, " + to + "!"}
}

func TestNewJSONGolden(t *testing.T) {
	golden := gilt.NewJSONGolden[HelloMessage, HelloMessage](t.Name())

	tests := []struct {
		name string
	}{
		{"world"},
		{"gopher"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := helloMessage(tt.name)
			golden.Assert(t, actual, tt.name, func(t *testing.T, actual HelloMessage, expected HelloMessage) {
				if actual.Message != expected.Message {
					t.Errorf("expected: %v\n, but was: %v", expected, actual)
				}
			})
		})
	}
}
```

By running `go test -update`, the following golden file will be generated:
```
testdata
└── testdata/TestNewJSONGolden
    └── testdata/TestNewJSONGolden/golden
        ├── testdata/TestNewJSONGolden/golden/gopher.golden
        └── testdata/TestNewJSONGolden/golden/world.golden
```

gopher.golden:
```
{
  "message": "Hello, gopher!"
}
```

When running `go test`, the golden file will be unmarshaled into the `HelloMessage` type and passed to the assertion function.  
You can use any assertion function you like.
```go
	golden.Assert(t, actual, tt.name, func(t *testing.T, actual HelloMessage, expected HelloMessage) {
		if actual.Message != expected.Message {
			t.Errorf("expected: %v\n, but was: %v", expected, actual)
		}
	})
```

## Advanced Usage

### Automatic cleanup for table-driven tests

For table-driven tests, enable cleanup at the test function level to avoid timing issues:

```go
func TestMyFeature(t *testing.T) {
    golden := gilt.NewStringGolden("my-feature")
    
    // Setup cleanup at the test function level
    gilt.SetupCleanup(t, golden, gilt.WithDryRun())
    
    testCases := []struct{
        name string
        input string
    }{
        {"case1", "test1"},
        {"case2", "test2"},
    }
    
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            golden.Assert(t, tc.input, tc.name, func(t *testing.T, actual, expected string) {
                if actual != expected {
                    t.Errorf("Expected %q, got %q", expected, actual)
                }
            })
        })
    }
    // Cleanup happens once at the end of TestMyFeature
}
```

### Customizing the reading and updating of Golden files
You can customize the reading and updating of Golden files by implementing the each of the following interfaces:

```mermaid
%%{init:{'theme':'neutral'}}%%
graph TB
  subgraph "Golden[TActual,TExpected]"
    Assert(("Assert"))
    GoldenFiles[("golden files")]
    IsUpdater["IsUpdater"]

    Assert ==> IsUpdater
    
    IsUpdater ==> Update
    IsUpdater ==> Load
    
    Assert --"actual TActual"--> Saver
    subgraph "Update"
      Saver["Saver"]
      Writer["Writer"]
      Saver ==> Writer
    end
    Writer --> GoldenFiles

    subgraph "Load"
      Reader["Reader"]
      Loader["Loader"]
      Reader -.-> Loader
      
    end
    GoldenFiles -.-> Reader
  
  end
  
  assertion[["assertion"]]
  Loader -."expected TExpected".-> assertion
  Assert --"actual TActual"--> assertion
```

For example, you can specify a `gilt.SaveLines` function for `Saver` and a `gilt.LoadLines` function for `Loader` as follows:
```go
golden := gilt.New(
	t.Name(),
	gilt.WithSaver[[]string, iter.Seq[string]](
		gilt.SaverFunc[[]string](gilt.SaveLines[string]),
	),
	gilt.WithLoader[[]string](
		gilt.LoaderFunc[iter.Seq[string]](gilt.LoadLines),
	),
)
```

Of course, you can also specify a function that you have implemented yourself.

### Custom cleanup behavior

You can customize cleanup behavior with options:

```go
func TestBasicCleanup(t *testing.T) {
    golden := gilt.NewStringGolden("basic")
    
    // Setup cleanup - only activates if -clean flags are used
    gilt.SetupCleanup(t, golden)
    
    // Your test logic here
}

func TestCustomCleanup(t *testing.T) {
    golden := gilt.NewStringGolden("custom")
    
    // Custom cleanup options - forces cleanup to be enabled
    gilt.SetupCleanup(t, golden,
        gilt.WithDryRun(),                     // Show what would be cleaned
        gilt.WithBackup(),                     // Create backups
        gilt.WithPatterns("*.keep", "*.ref"),  // Exclude important files
    )
    
    // Your test logic here
}

func TestEnvironmentBasedCleanup(t *testing.T) {
    golden := gilt.NewStringGolden("env")
    
    // Different behavior based on environment
    if os.Getenv("CI") == "true" {
        // Aggressive cleanup in CI
        gilt.SetupCleanup(t, golden, gilt.WithBackup())
    } else {
        // Safe cleanup in development
        gilt.SetupCleanup(t, golden, gilt.WithDryRun())
    }
    
    // Your test logic here
}
```