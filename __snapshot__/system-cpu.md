# chrome.system.cpu

## Description

Use the `system.cpu` API to query CPU metadata.

## Permissions

`system.cpu`

## Types

### CpuInfo

#### Properties

- archName
  
  string
  
  The architecture name of the processors.
- features
  
  string\[]
  
  A set of feature codes indicating some of the processor's capabilities. The currently supported codes are "mmx", "sse", "sse2", "sse3", "ssse3", "sse4\_1", "sse4\_2", and "avx".
- modelName
  
  string
  
  The model name of the processors.
- numOfProcessors
  
  number
  
  The number of logical processors.
- processors
  
  [ProcessorInfo](#type-ProcessorInfo)\[]
  
  Information about each logical processor.
- temperatures
  
  number\[]
  
  Chrome 60+
  
  List of CPU temperature readings from each thermal zone of the CPU. Temperatures are in degrees Celsius.
  
  **Currently supported on Chrome OS only.**

### CpuTime

#### Properties

- idle
  
  number
  
  The cumulative time spent idle by this processor.
- kernel
  
  number
  
  The cumulative time used by kernel programs on this processor.
- total
  
  number
  
  The total cumulative time for this processor. This value is equal to user + kernel + idle.
- user
  
  number
  
  The cumulative time used by userspace programs on this processor.

### ProcessorInfo

#### Properties

- usage
  
  [CpuTime](#type-CpuTime)
  
  Cumulative usage info for this logical processor.

## Methods

### getInfo()

```
chrome.system.cpu.getInfo(): Promise<CpuInfo>
```

Queries basic CPU information of the system.

#### Returns

- Promise&lt;[CpuInfo](#type-CpuInfo)&gt;
  
  Chrome 91+