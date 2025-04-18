# chrome.enterprise.deviceAttributes

This API is only for [extensions installed by a policy](https://support.google.com/chrome/a/answer/1375694).

**Important:** This API works **only on ChromeOS**.

## Description

Use the `chrome.enterprise.deviceAttributes` API to read device attributes. Note: This API is only available to extensions force-installed by enterprise policy.

## Permissions

`enterprise.deviceAttributes`

## Availability

Chrome 46+ ChromeOS only [Requires policy](https://support.google.com/chrome/a/answer/9296680)

## Methods

### getDeviceAnnotatedLocation()

Promise Chrome 66+

```
chrome.enterprise.deviceAttributes.getDeviceAnnotatedLocation(
  callback?: function,
)
```

Fetches the administrator-annotated Location. If the current user is not affiliated or no Annotated Location has been set by the administrator, returns an empty string.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (annotatedLocation: string) => void
  ```
  
  - annotatedLocation
    
    string

#### Returns

- Promise&lt;string&gt;
  
  Chrome 96+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### getDeviceAssetId()

Promise Chrome 66+

```
chrome.enterprise.deviceAttributes.getDeviceAssetId(
  callback?: function,
)
```

Fetches the administrator-annotated Asset Id. If the current user is not affiliated or no Asset Id has been set by the administrator, returns an empty string.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (assetId: string) => void
  ```
  
  - assetId
    
    string

#### Returns

- Promise&lt;string&gt;
  
  Chrome 96+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### getDeviceHostname()

Promise Chrome 82+

```
chrome.enterprise.deviceAttributes.getDeviceHostname(
  callback?: function,
)
```

Fetches the device's hostname as set by DeviceHostnameTemplate policy. If the current user is not affiliated or no hostname has been set by the enterprise policy, returns an empty string.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (hostname: string) => void
  ```
  
  - hostname
    
    string

#### Returns

- Promise&lt;string&gt;
  
  Chrome 96+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### getDeviceSerialNumber()

Promise Chrome 66+

```
chrome.enterprise.deviceAttributes.getDeviceSerialNumber(
  callback?: function,
)
```

Fetches the device's serial number. Please note the purpose of this API is to administrate the device (e.g. generating Certificate Sign Requests for device-wide certificates). This API may not be used for tracking devices without the consent of the device's administrator. If the current user is not affiliated, returns an empty string.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (serialNumber: string) => void
  ```
  
  - serialNumber
    
    string

#### Returns

- Promise&lt;string&gt;
  
  Chrome 96+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### getDirectoryDeviceId()

Promise

```
chrome.enterprise.deviceAttributes.getDirectoryDeviceId(
  callback?: function,
)
```

Fetches the value of [the device identifier of the directory API](https://developers.google.com/admin-sdk/directory/v1/guides/manage-chrome-devices), that is generated by the server and identifies the cloud record of the device for querying in the cloud directory API. If the current user is not affiliated, returns an empty string.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (deviceId: string) => void
  ```
  
  - deviceId
    
    string

#### Returns

- Promise&lt;string&gt;
  
  Chrome 96+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.