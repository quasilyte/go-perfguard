## Statistics

### Without tests

Import stats.

```
   7799 "fmt"
   4702 "context"
   4450 "strings"
   3628 "time"
   2519 "os"
   2172 "io"
   1862 "sync"
   1768 "bytes"
   1651 "strconv"
   1536 "errors"
   1401 "net/http"
   1323 "path/filepath"
   1261 "sort"
   1109 "encoding/json"
    933 "math"
    878 "net"
    762 "io/ioutil"
    738 "reflect"
    718 "unsafe"
    685 "net/url"
    674 "regexp"
    665 "log"
    627 "runtime"
    609 "syscall"
    568 "github.com/pkg/errors"
    484 "path"
    477 "testing"
    441 "github.com/spf13/cobra"
    419 "bufio"
    374 "sync/atomic"
    355 "go/token"
    345 "flag"
    327 "github.com/sirupsen/logrus"
    321 "github.com/onsi/ginkgo"
    316 "golang.org/x/tools/go/analysis"
    315 "unicode/utf8"
    308 "math/rand"
    297 "os/exec"
    293 "go/ast"
    279 "encoding/binary"
    276 "gonum.org/v1/gonum/blas"
    249 "go/types"
    246 "unicode"
    227 "golang.org/x/exp/rand"
    225 "crypto/tls"
    206 "go.uber.org/zap"
    206 "github.com/spf13/pflag"
    184 "encoding/base64"
```

`fmt` package stats.
```
  29267 Errorf
  19352 Sprintf
   5135 Fprintf
   1347 Printf
    952 Fprintln
    740 Println
    704 Fprint
    330 Sprint
    103 Print
     47 Sscanf
     32 Sprintln
      7 Sscan
      6 Fscanf
      6 Formatter
      5 Scanln
      3 Scanf
      1 Fscan
```

`strings` package stats.
```
   3299 Replace
   3201 Join
   2447 HasPrefix
   1786 Split
   1601 Contains
   1586 TrimSpace
   1339 ToLower
    944 HasSuffix
    742 TrimPrefix
    480 TrimSuffix
    463 Index
    363 EqualFold
    340 SplitN
    259 ReplaceAll
    257 ToUpper
    212 Fields
    201 Repeat
    191 LastIndex
    188 Trim
    172 NewReader
    124 Count
    118 IndexByte
    117 TrimRight
    105 NewReplacer
     85 ContainsAny
     81 TrimLeft
     70 Title
     66 Cut
     47 ContainsRune
     42 IndexRune
     41 Compare
     33 LastIndexByte
     22 FieldsFunc
     20 IndexFunc
     20 IndexAny
     17 SplitAfter
     17 Map
     15 TrimLeftFunc
     13 TrimRightFunc
      9 LastIndexAny
      4 ToTitle
      3 SplitAfterN
      3 LastIndexFunc
      2 TrimFunc
      1 ToValidUTF8
```

`bytes` package stats.
```
    519 Compare
    475 Equal
    441 NewReader
    250 HasPrefix
    225 NewBuffer
    173 IndexByte
    112 Index
    104 Contains
    101 TrimSpace
     48 Split
     47 NewBufferString
     38 HasSuffix
     27 IndexAny
     25 Replace
     25 Join
     23 TrimRight
     22 LastIndexByte
     20 Count
     15 TrimPrefix
     15 LastIndex
     14 Get
     11 SplitN
     10 Fields
     10 Cut
      8 ReplaceAll
      8 IndexFunc
      7 ToLower
      6 Trim
      5 TrimLeftFunc
      5 Repeat
      5 EqualFold
      5 ContainsAny
      4 IndexRune
      3 TrimSuffix
      3 TrimLeft
      3 ToUpper
      3 SplitAfter
      3 Runes
      3 ContainsRune
      2 readBase128Int
      1 Value
      1 TrimRightFunc
      1 SplitAfterN
      1 Set
      1 Map
      1 Close
```

`time` package stats.
```
   2086 Now
   1330 Duration
    649 Sleep
    530 Since
    359 After
    243 Unix
    225 NewTicker
    212 Date
    170 NewTimer
    142 ParseDuration
    131 Parse
     68 AfterFunc
     44 Until
     18 ParseInLocation
     10 Time
      9 Tick
      9 Month
      9 LoadLocation
      7 FixedZone
      3 UnixNano
      2 IsZero
      2 Format
```

`os` package stats.
```
    771 Exit
    751 Getenv
    749 Stat
    615 IsNotExist
    439 Open
    393 MkdirAll
    333 Remove
    288 RemoveAll
    261 Create
    237 ReadFile
    206 OpenFile
    150 WriteFile
    138 Setenv
    121 Environ
     99 Lstat
     98 Rename
     88 FileMode
     80 Getwd
     73 NewSyscallError
     73 Getpid
     59 Chmod
     57 MkdirTemp
     56 LookupEnv
     52 ReadDir
     48 TempDir
     42 Mkdir
     41 IsPathSeparator
     41 IsExist
     40 Hostname
     40 CreateTemp
     35 ExpandEnv
     31 Readlink
     31 Pipe
     28 Symlink
     24 NewFile
     21 Chown
     18 Executable
     17 Chtimes
     14 SameFile
     14 IsPermission
     14 Chdir
     10 NewComputeV2
     10 FindProcess
      9 Unsetenv
      9 Getuid
      8 getVolume
      7 Link
      5 Sweep
      5 Root
      4 volumeService
      4 UserHomeDir
      4 SetTransitionState
      4 Lchown
      4 LastRuleName
      4 Getppid
      4 GetOpts
      4 Done
      3 StartProcess
      3 Getgid
      3 Expand
      2 UserConfigDir
      2 UserCacheDir
      2 Next
      2 NewNetworkV2
      2 NewBlockStorageV3
      2 NewBlockStorageV2
      2 NewBlockStorageV1
      2 IsBetter
      2 ID
      2 Geteuid
      2 DiskIsAttached
      2 Add
      1 WithLogger
      1 Var
      1 Truncate
      1 shouldRemoveRemoteObject
      1 NewLoadBalancerV2
      1 New
      1 listCode
      1 Len
      1 isMemberFullyOptimized
      1 Getpagesize
      1 GetNodeNameByID
      1 Getegid
      1 getDevicePathFromInstanceMetadata
      1 GetDevicePathBySerialID
      1 DisksAreAttached
      1 diskIsUsed
      1 Difference
      1 DeepCopy
      1 CreatedAt
      1 Clearenv
```

`io` package stats.
```
    401 Copy
    277 ReadFull
    196 WriteString
    164 ReadAll
     95 Pipe
     81 LimitReader
     57 NewSectionReader
     39 NopCloser
     37 MultiWriter
     35 MultiReader
     34 CopyN
     16 TeeReader
     16 ReadDir
     11 CopyBuffer
      8 Writer
      8 ReadFile
      5 Reader
      5 ReadAtLeast
      3 WriteFile
      3 IsInconsistentReadError
      3 EvalSymlinks
      2 Readlink
      2 Lstat
      1 Close
      1 Cancel
```

`reflect` package stats.
```
   1216 TypeOf
    999 ValueOf
    440 DeepEqual
    240 Indirect
    183 New
     91 Zero
     53 MakeSlice
     51 Append
     25 PtrTo
     18 SliceOf
     14 MakeMap
     13 StructTag
     12 PointerTo
     10 Copy
      4 MapOf
      3 MakeMapWithSize
      1 StructOf
      1 Select
      1 NewAt
      1 MakeFunc
      1 ChanDir
      1 AppendSlice
```

`strconv` package stats.
```
    969 Itoa
    838 Atoi
    453 FormatInt
    443 ParseInt
    233 ParseUint
    184 ParseBool
    131 FormatBool
    124 ParseFloat
    111 FormatUint
    105 Unquote
    101 Quote
     60 AppendUint
     54 AppendInt
     51 FormatFloat
     31 AppendFloat
     10 UnquoteChar
      6 IsPrint
      6 AppendBool
      3 QuoteToASCII
      3 CanBackquote
      3 AppendQuote
      2 AppendQuoteToASCII
      2 AppendQuoteRuneToASCII
      2 AppendQuoteRune
      1 QuoteRune
```

### All files (with _test.go files)

Import stats.

```
  11356 "fmt"
   9463 "testing"
   6883 "context"
   6848 "strings"
   5859 "time"
   4347 "os"
   3285 "bytes"
   3036 "io"
   2453 "reflect"
   2406 "sync"
   2331 "net/http"
   2326 "path/filepath"
   2185 "strconv"
   2036 "errors"
   1911 "github.com/cockroachdb/errors"
   1769 "sort"
   1595 "io/ioutil"
   1450 "github.com/stretchr/testify/require"
   1439 "encoding/json"
   1417 "math"
   1281 "net"
   1228 "github.com/stretchr/testify/assert"
   1165 "runtime"
   1023 "net/url"
    960 "regexp"
    908 "log"
    809 "unsafe"
    763 "syscall"
    746 "math/rand"
    646 "path"
    619 "github.com/pkg/errors"
    602 "sync/atomic"
    564 "bufio"
    556 "os/exec"
    530 "flag"
    475 "github.com/spf13/cobra"
    468 "go/token"
    420 "net/http/httptest"
    373 "github.com/onsi/ginkgo"
    370 "go/ast"
    366 "gotest.tools/v3/assert"
    357 "unicode/utf8"
    352 "golang.org/x/exp/rand"
    351 "github.com/sirupsen/logrus"
    337 "crypto/tls"
    328 "encoding/binary"
    321 "golang.org/x/tools/go/analysis"
    300 "go/types"
    300 "gonum.org/v1/gonum/blas"
```

`fmt` package stats.
```
  33909 Errorf
  29740 Sprintf
   5994 Fprintf
   2550 Printf
   2145 Println
   1137 Fprintln
    811 Fprint
    767 Sprint
    160 Print
     79 Sscanf
     39 Sprintln
     13 Sscan
     12 Formatter
     10 Fscanf
      5 Scanln
      3 Scanf
      3 Fscanln
      1 Fscan
```

`strings` package stats.
```
   5618 Contains
   4027 Join
   3523 Replace
   3346 TrimSpace
   3136 HasPrefix
   2457 Split
   1515 ToLower
   1247 HasSuffix
   1097 NewReader
   1074 Repeat
    887 TrimPrefix
    545 Index
    544 TrimSuffix
    415 ReplaceAll
    396 EqualFold
    384 SplitN
    340 ToUpper
    293 Fields
    277 Trim
    235 Count
    220 LastIndex
    134 TrimRight
    133 IndexByte
    120 NewReplacer
    100 ContainsAny
     98 Cut
     97 TrimLeft
     82 Compare
     78 Title
     56 ContainsRune
     52 IndexRune
     40 LastIndexByte
     25 Map
     24 IndexAny
     24 FieldsFunc
     22 IndexFunc
     19 SplitAfter
     18 TrimLeftFunc
     14 TrimRightFunc
     13 LastIndexAny
      8 SplitAfterN
      7 ToTitle
      6 LastIndexFunc
      4 TrimFunc
      3 Clone
      2 ToLowerSpecial
      1 ToValidUTF8
      1 ToUpperSpecial
      1 ToTitleSpecial
```

`bytes` package stats.
```
   1727 Equal
   1547 NewReader
    659 NewBuffer
    606 Compare
    341 NewBufferString
    312 HasPrefix
    285 Contains
    187 Repeat
    178 IndexByte
    141 Index
    128 TrimSpace
     77 Split
     75 Join
     52 HasSuffix
     37 Replace
     29 IndexAny
     29 Count
     27 TrimPrefix
     26 TrimRight
     26 LastIndexByte
     18 ReplaceAll
     17 LastIndex
     17 Get
     16 SplitN
     13 Cut
     11 Fields
     10 SplitAfter
     10 IndexFunc
      9 Trim
      9 ToLower
      9 ContainsAny
      8 TrimSuffix
      8 TrimLeftFunc
      8 ContainsRune
      7 IndexRune
      7 EqualFold
      6 TrimLeft
      6 ToUpper
      4 TrimRightFunc
      4 TrimFunc
      4 Runes
      3 LastIndexFunc
      3 LastIndexAny
      2 ToTitle
      2 SplitAfterN
      2 readBase128Int
      2 Len
      1 Value
      1 ToUpperSpecial
      1 ToTitleSpecial
      1 ToLowerSpecial
      1 Title
      1 Set
      1 Map
      1 FieldsFunc
      1 Close
```

`time` package stats.
```
   4318 Now
   2013 Duration
   1900 Sleep
   1721 After
   1281 Date
    855 Unix
    704 Since
    263 NewTicker
    231 Parse
    218 NewTimer
    185 ParseDuration
    181 AfterFunc
    124 FixedZone
     66 Until
     47 LoadLocation
     35 Time
     31 Tick
     28 ParseInLocation
     16 Month
     10 ForceZipFileForTesting
      7 Format
      5 Add
      4 UnixNano
      4 ResetZoneinfoForTesting
      4 ResetLocalOnceForTest
      4 LoadLocationFromTZData
      3 Zone
      3 Minute
      3 Hour
      2 ZoneinfoForTesting
      2 Year
      2 Weekday
      2 Unmarshal
      2 Second
      2 Nanosecond
      2 IsZero
      2 Equal
      2 Day
      1 UnixMilli
      1 UnixMicro
      1 TzsetRule
      1 TzsetOffset
      1 TzsetName
      1 Tzset
      1 LoadTzinfo
      1 LoadFromEmbeddedTZData
      1 ForceUSPacificForTesting
```

`os` package stats.
```
   2092 RemoveAll
   1388 Stat
   1254 Getenv
   1213 Exit
   1137 Remove
    883 IsNotExist
    824 MkdirAll
    729 Open
    583 Setenv
    555 ReadFile
    494 WriteFile
    487 MkdirTemp
    463 Create
    284 OpenFile
    238 Environ
    205 FileMode
    203 TempDir
    200 Getwd
    191 CreateTemp
    188 Mkdir
    159 Symlink
    155 Chdir
    152 Lstat
    138 Rename
    126 Unsetenv
    122 Chmod
    105 ReadDir
    102 Getpid
     97 Getuid
     80 Pipe
     76 NewSyscallError
     72 LookupEnv
     60 Readlink
     53 Hostname
     46 IsExist
     44 IsPathSeparator
     42 NewFile
     37 SameFile
     36 ExpandEnv
     32 Chown
     25 Executable
     24 FindProcess
     20 IsPermission
     18 Chtimes
     17 Link
     16 Getgid
     16 Geteuid
     16 DirFS
     13 Getegid
     11 NewComputeV2
      9 getVolume
      9 Getpagesize
      6 StartProcess
      5 Sweep
      5 Root
      5 IsTimeout
      5 Getppid
      5 Expand
      4 volumeService
      4 UserHomeDir
      4 Truncate
      4 SetTransitionState
      4 Lchown
      4 LastRuleName
      4 GetOpts
      4 Done
      4 Clearenv
      3 WithLogger
      3 NewNetworkV2
      3 New
      2 WithRecording
      2 UserConfigDir
      2 UserCacheDir
      2 Next
      2 NewBlockStorageV3
      2 NewBlockStorageV2
      2 NewBlockStorageV1
      2 IsBetter
      2 InstanceID
      2 ID
      2 DiskIsAttached
      2 CommandLineToArgv
      2 Add
      1 Zones
      1 Var
      1 shouldRemoveRemoteObject
      1 Routes
      1 NewLoadBalancerV2
      1 NewConsoleFile
      1 LoadBalancer
      1 listCode
      1 Len
      1 isMemberFullyOptimized
      1 GetNodeNameByID
      1 getDevicePathFromInstanceMetadata
      1 GetDevicePathBySerialID
      1 GetDevicePath
      1 FixLongPath
      1 ExpandVolume
      1 DisksAreAttached
```

`io` package stats.
```
    820 Copy
    550 ReadAll
    485 WriteString
    395 ReadFull
    271 NopCloser
    150 Pipe
    107 LimitReader
     65 NewSectionReader
     58 CopyN
     54 MultiWriter
     51 MultiReader
     26 CopyBuffer
     19 TeeReader
     16 ReadDir
     12 ReadAtLeast
     10 Writer
      8 ReadFile
      8 Reader
      3 WriteFile
      3 IsInconsistentReadError
      3 EvalSymlinks
      2 Readlink
      2 Lstat
      2 Closer
      1 Close
      1 Cancel
```

`reflect` package stats
```
   4747 DeepEqual
   1736 TypeOf
   1544 ValueOf
    266 Indirect
    262 New
     97 Zero
     56 MakeSlice
     52 Append
     25 PtrTo
     18 SliceOf
     17 StructTag
     17 MakeMap
     12 PointerTo
     10 Copy
      6 SetArgRegs
      4 MapOf
      4 MakeFunc
      3 StructOf
      3 MakeMapWithSize
      2 Select
      2 NewAt
      1 ChanDir
      1 AppendSlice
```

`strconv` package stats.
```
   1533 Itoa
   1082 Atoi
    565 FormatInt
    516 ParseInt
    270 ParseUint
    190 ParseBool
    180 ParseFloat
    149 FormatBool
    143 FormatUint
    118 Unquote
    104 Quote
     85 FormatFloat
     73 AppendUint
     59 AppendInt
     33 AppendFloat
     11 UnquoteChar
      8 IsPrint
      7 AppendBool
      5 QuoteToASCII
      5 CanBackquote
      4 QuoteRuneToGraphic
      4 AppendQuote
      3 QuoteToGraphic
      3 IsGraphic
      3 AppendQuoteToASCII
      3 AppendQuoteRuneToASCII
      3 AppendQuoteRune
      2 QuoteRune
      1 QuoteRuneToASCII
```
