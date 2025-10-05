# chrome.system.memory

## Description

The `chrome.system.memory` API.

## Permissions

`system.memory`

## Types

### MemoryInfo

#### Properties

- availableCapacity
  
  number
  
  The amount of available capacity, in bytes.
- capacity
  
  number
  
  The total amount of physical memory capacity, in bytes.

## Methods

### getInfo()

```
chrome.system.memory.getInfo(): Promise<MemoryInfo>
```

Get physical memory information.

#### Returns

- Promise&lt;[MemoryInfo](#type-MemoryInfo)&gt;
  
  Chrome 91+