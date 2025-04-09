# chrome.enterprise.hardwarePlatform

This API is only for [extensions installed by a policy](https://cloud.google.com/blog/products/chrome-enterprise/how-to-manage-chrome-extensions-on-windows-and-mac). The [`EnterpriseHardwarePlatformAPIEnabled`](https://chromeenterprise.google/policies/?policy=EnterpriseHardwarePlatformAPIEnabled) key must also be set.

## Description

Use the `chrome.enterprise.hardwarePlatform` API to get the manufacturer and model of the hardware platform where the browser runs. Note: This API is only available to extensions installed by enterprise policy.

## Permissions

`enterprise.hardwarePlatform`

## Availability

Chrome 71+ [Requires policy](https://support.google.com/chrome/a/answer/9296680)

## Types

### HardwarePlatformInfo

#### Properties

- manufacturer
  
  string
- model
  
  string

## Methods

### getHardwarePlatformInfo()

Promise

```
chrome.enterprise.hardwarePlatform.getHardwarePlatformInfo(
  callback?: function,
)
```

Obtains the manufacturer and model for the hardware platform and, if the extension is authorized, returns it via `callback`.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (info: HardwarePlatformInfo) => void
  ```
  
  - info
    
    [HardwarePlatformInfo](#type-HardwarePlatformInfo)

#### Returns

- Promise&lt;[HardwarePlatformInfo](#type-HardwarePlatformInfo)&gt;
  
  Chrome 96+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.