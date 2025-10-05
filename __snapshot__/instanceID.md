# chrome.instanceID

## Description

Use `chrome.instanceID` to access the Instance ID service.

## Permissions

`gcm`

## Availability

Chrome 44+

## Methods

### deleteID()

```
chrome.instanceID.deleteID(): Promise<void>
```

Resets the app instance identifier and revokes all tokens associated with it.

#### Returns

- Promise&lt;void&gt;
  
  Chrome 96+

### deleteToken()

```
chrome.instanceID.deleteToken(
  deleteTokenParams: object,
): Promise<void>
```

Revokes a granted token.

#### Parameters

- deleteTokenParams
  
  object
  
  Parameters for deleteToken.
  
  - authorizedEntity
    
    string
    
    Chrome 46+
    
    The authorized entity that is used to obtain the token.
  - scope
    
    string
    
    Chrome 46+
    
    The scope that is used to obtain the token.

#### Returns

- Promise&lt;void&gt;
  
  Chrome 96+

### getCreationTime()

```
chrome.instanceID.getCreationTime(): Promise<number>
```

Retrieves the time when the InstanceID has been generated. The creation time will be returned by the `callback`.

#### Returns

- Promise&lt;number&gt;
  
  Chrome 96+

### getID()

```
chrome.instanceID.getID(): Promise<string>
```

Retrieves an identifier for the app instance. The instance ID will be returned by the `callback`. The same ID will be returned as long as the application identity has not been revoked or expired.

#### Returns

- Promise&lt;string&gt;
  
  Chrome 96+

### getToken()

```
chrome.instanceID.getToken(
  getTokenParams: object,
): Promise<string>
```

Return a token that allows the authorized entity to access the service defined by scope.

#### Parameters

- getTokenParams
  
  object
  
  Parameters for getToken.
  
  - authorizedEntity
    
    string
    
    Chrome 46+
    
    Identifies the entity that is authorized to access resources associated with this Instance ID. It can be a project ID from [Google developer console](https://code.google.com/apis/console).
  - options
    
    object optional
    
    Chrome 46+ Deprecated since Chrome 89
    
    options are deprecated and will be ignored.
    
    Allows including a small number of string key/value pairs that will be associated with the token and may be used in processing the request.
  - scope
    
    string
    
    Chrome 46+
    
    Identifies authorized actions that the authorized entity can take. E.g. for sending GCM messages, `GCM` scope should be used.

#### Returns

- Promise&lt;string&gt;
  
  Chrome 96+

## Events

### onTokenRefresh

```
chrome.instanceID.onTokenRefresh.addListener(
  callback: function,
)
```

Fired when all the granted tokens need to be refreshed.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```