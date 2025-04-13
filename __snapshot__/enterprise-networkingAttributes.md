# chrome.enterprise.networkingAttributes

This API is only for [extensions installed by a policy](https://support.google.com/chrome/a/answer/1375694).

**Important:** This API works **only on ChromeOS**.

## Description

Use the `chrome.enterprise.networkingAttributes` API to read information about your current network. Note: This API is only available to extensions force-installed by enterprise policy.

## Permissions

`enterprise.networkingAttributes`

## Availability

Chrome 85+ ChromeOS only [Requires policy](https://support.google.com/chrome/a/answer/9296680)

## Types

### NetworkDetails

#### Properties

- ipv4
  
  string optional
  
  The device's local IPv4 address (undefined if not configured).
- ipv6
  
  string optional
  
  The device's local IPv6 address (undefined if not configured).
- macAddress
  
  string
  
  The device's MAC address.

## Methods

### getNetworkDetails()

Promise

```
chrome.enterprise.networkingAttributes.getNetworkDetails(
  callback?: function,
)
```

Retrieves the network details of the device's default network. If the user is not affiliated or the device is not connected to a network, [`runtime.lastError`](https://developer.chrome.com/docs/extensions/reference/api/runtime/#property-lastError) will be set with a failure reason.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (networkAddresses: NetworkDetails) => void
  ```
  
  - networkAddresses
    
    [NetworkDetails](#type-NetworkDetails)

#### Returns

- Promise&lt;[NetworkDetails](#type-NetworkDetails)&gt;
  
  Chrome 96+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.