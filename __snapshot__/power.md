# chrome.power

## Description

Use the `chrome.power` API to override the system's power management features.

## Permissions

`power`

## Concepts and usage

By default, operating systems dim the screen when users are inactive and eventually suspend the system. With the power API, an app or extension can keep the system awake.

Using this API, you can specify the [Level](#type-Level) to which power management is disabled. The `"system"` level keeps the system active, but allows the screen to be dimmed or turned off. For example, a communication app can continue to receive messages while the screen is off. The `"display"` level keeps the screen and system active. E-book and presentation apps, for example, can keep the screen and system active while users read.

When a user has more than one app or extension active, each with its own power level, the highest-precedence level takes effect; `"display"` always takes precedence over `"system"`. For example, if app A asks for `"system"` power management, and app B asks for `"display"`, `"display"` is used until app B is unloaded or releases its request. If app A is still active, `"system"` is then used.

## Types

### Level

#### Enum

"system"  
Prevents the system from sleeping in response to user inactivity.

"display"  
Prevents the display from being turned off or dimmed, or the system from sleeping in response to user inactivity.

## Methods

### releaseKeepAwake()

```
chrome.power.releaseKeepAwake()
```

Releases a request previously made via requestKeepAwake().

### reportActivity()

Promise Chrome 113+ ChromeOS only

```
chrome.power.reportActivity(
  callback?: function,
)
```

Reports a user activity in order to awake the screen from a dimmed or turned off state or from a screensaver. Exits the screensaver if it is currently active.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

#### Returns

- Promise&lt;void&gt;
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### requestKeepAwake()

```
chrome.power.requestKeepAwake(
  level: Level,
)
```

Requests that power management be temporarily disabled. `level` describes the degree to which power management should be disabled. If a request previously made by the same app is still active, it will be replaced by the new request.

#### Parameters

- level
  
  [Level](#type-Level)