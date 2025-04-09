# chrome.loginState

**Important:** This API works **only on ChromeOS**.

## Description

Use the `chrome.loginState` API to read and monitor the login state.

## Permissions

`loginState`

## Availability

Chrome 78+ ChromeOS only

## Types

### ProfileType

#### Enum

"SIGNIN\_PROFILE"  
Specifies that the extension is in the signin profile.

"USER\_PROFILE"  
Specifies that the extension is in the user profile.

### SessionState

#### Enum

"UNKNOWN"  
Specifies that the session state is unknown.

"IN\_OOBE\_SCREEN"  
Specifies that the user is in the out-of-box-experience screen.

"IN\_LOGIN\_SCREEN"  
Specifies that the user is in the login screen.

"IN\_SESSION"  
Specifies that the user is in the session.

"IN\_LOCK\_SCREEN"  
Specifies that the user is in the lock screen.

"IN\_RMA\_SCREEN"  
Specifies that the device is in RMA mode, finalizing repairs.

## Methods

### getProfileType()

Promise

```
chrome.loginState.getProfileType(
  callback?: function,
)
```

Gets the type of the profile the extension is in.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (result: ProfileType) => void
  ```
  
  - result
    
    [ProfileType](#type-ProfileType)

#### Returns

- Promise&lt;[ProfileType](#type-ProfileType)&gt;
  
  Chrome 96+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### getSessionState()

Promise

```
chrome.loginState.getSessionState(
  callback?: function,
)
```

Gets the current session state.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (result: SessionState) => void
  ```
  
  - result
    
    [SessionState](#type-SessionState)

#### Returns

- Promise&lt;[SessionState](#type-SessionState)&gt;
  
  Chrome 96+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

## Events

### onSessionStateChanged

```
chrome.loginState.onSessionStateChanged.addListener(
  callback: function,
)
```

Dispatched when the session state changes. `sessionState` is the new session state.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (sessionState: SessionState) => void
  ```
  
  - sessionState
    
    [SessionState](#type-SessionState)