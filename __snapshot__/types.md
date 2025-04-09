# chrome.types

## Description

The `chrome.types` API contains type declarations for Chrome.

## Chrome settings

The [`ChromeSetting`](#type-ChromeSetting) type provides a common set of functions (`get()`, `set()`, and `clear()`) as well as an event publisher (`onChange`) for settings of the Chrome browser. The [proxy settings examples](/docs/extensions/reference/api/proxy#examples) demonstrate how these functions are intended to be used.

### Scope and lifecycle

Chrome distinguishes between three different scopes of browser settings:

`regular`

Settings set in the `regular` scope apply to regular browser windows and are inherited by incognito windows if they are not overwritten. These settings are stored to disk and remain in place until they are cleared by the governing extension, or the governing extension is disabled or uninstalled.

`incognito_persistent`

Settings set in the `incognito_persistent` scope apply only to incognito windows. For these, they override `regular` settings. These settings are stored to disk and remain in place until they are cleared by the governing extension, or the governing extension is disabled or uninstalled.

`incognito_session_only`

Settings set in the `incognito_session_only` scope apply only to incognito windows. For these, they override `regular` and `incognito_persistent` settings. These settings are not stored to disk and are cleared when the last incognito window is closed. They can only be set when at least one incognito window is open.

### Precedence

Chrome manages settings on different layers. The following list describes the layers that may influence the effective settings, in increasing order of precedence.

1. System settings provided by the operating system
2. Command-line parameters
3. Settings provided by extensions
4. Policies

As the list implies, policies might overrule any changes that you specify with your extension. You can use the `get()` function to determine whether your extension is capable of providing a setting or whether this setting would be overridden.

As discussed previously, Chrome allows using different settings for regular windows and incognito windows. The following example illustrates the behavior. Assume that no policy overrides the settings and that an extension can set settings for regular windows **(R)** and settings for incognito windows **(I)**.

- If only **(R)** is set, these settings are effective for both regular and incognito windows.
- If only **(I)** is set, these settings are effective for only incognito windows. Regular windows use the settings determined by the lower layers (command-line options and system settings).
- If both **(R)** and **(I)** are set, the respective settings are used for regular and incognito windows.

If two or more extensions want to set the same setting to different values, the extension installed most recently takes precedence over the other extensions. If the most recently installed extension sets only **(I)**, the settings of regular windows can be defined by previously installed extensions.

The *effective* value of a setting is the one that results from considering the precedence rules. It is used by Chrome.

## Types

### ChromeSetting

An interface that allows access to a Chrome browser setting. See [`accessibilityFeatures`](https://developer.chrome.com/docs/extensions/reference/accessibilityFeatures/) for an example.

#### Properties

- onChange
  
  Event&lt;functionvoidvoid&gt;
  
  Fired after the setting changes.
  
  The `onChange.addListener` function looks like:
  
  ```
  (callback: function) => {...}
  ```
  
  - callback
    
    function
    
    The `callback` parameter looks like:
    
    ```
    (details: object) => void
    ```
    
    - details
      
      object
      
      - incognitoSpecific
        
        boolean optional
        
        Whether the value that has changed is specific to the incognito session. This property will *only* be present if the user has enabled the extension in incognito mode.
      - levelOfControl
        
        [LevelOfControl](#type-LevelOfControl)
        
        The level of control of the setting.
      - value
        
        T
        
        The value of the setting after the change.
- clear
  
  void
  
  Promise
  
  Clears the setting, restoring any default value.
  
  The `clear` function looks like:
  
  ```
  (details: object, callback?: function) => {...}
  ```
  
  - details
    
    object
    
    Which setting to clear.
    
    - scope
      
      [ChromeSettingScope](#type-ChromeSettingScope) optional
      
      Where to clear the setting (default: regular).
  - callback
    
    function optional
    
    The `callback` parameter looks like:
    
    ```
    () => void
    ```
  
  <!--THE END-->
  
  - returns
    
    Promise&lt;void&gt;
    
    Chrome 96+
    
    Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.
- get
  
  void
  
  Promise
  
  Gets the value of a setting.
  
  The `get` function looks like:
  
  ```
  (details: object, callback?: function) => {...}
  ```
  
  - details
    
    object
    
    Which setting to consider.
    
    - incognito
      
      boolean optional
      
      Whether to return the value that applies to the incognito session (default false).
  - callback
    
    function optional
    
    The `callback` parameter looks like:
    
    ```
    (details: object) => void
    ```
    
    - details
      
      object
      
      Details of the currently effective value.
      
      - incognitoSpecific
        
        boolean optional
        
        Whether the effective value is specific to the incognito session. This property will *only* be present if the `incognito` property in the `details` parameter of `get()` was true.
      - levelOfControl
        
        [LevelOfControl](#type-LevelOfControl)
        
        The level of control of the setting.
      - value
        
        T
        
        The value of the setting.
  
  <!--THE END-->
  
  - returns
    
    Promise&lt;object&gt;
    
    Chrome 96+
    
    Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.
- set
  
  void
  
  Promise
  
  Sets the value of a setting.
  
  The `set` function looks like:
  
  ```
  (details: object, callback?: function) => {...}
  ```
  
  - details
    
    object
    
    Which setting to change.
    
    - scope
      
      [ChromeSettingScope](#type-ChromeSettingScope) optional
      
      Where to set the setting (default: regular).
    - value
      
      T
      
      The value of the setting. Note that every setting has a specific value type, which is described together with the setting. An extension should *not* set a value of a different type.
  - callback
    
    function optional
    
    The `callback` parameter looks like:
    
    ```
    () => void
    ```
  
  <!--THE END-->
  
  - returns
    
    Promise&lt;void&gt;
    
    Chrome 96+
    
    Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### ChromeSettingScope

Chrome 44+

The scope of the ChromeSetting. One of

- `regular`: setting for the regular profile (which is inherited by the incognito profile if not overridden elsewhere),
- `regular\_only`: setting for the regular profile only (not inherited by the incognito profile),
- `incognito\_persistent`: setting for the incognito profile that survives browser restarts (overrides regular preferences),
- `incognito\_session\_only`: setting for the incognito profile that can only be set during an incognito session and is deleted when the incognito session ends (overrides regular and incognito\_persistent preferences).

#### Enum

"regular"

"regular\_only"

"incognito\_persistent"

"incognito\_session\_only"

### LevelOfControl

Chrome 44+

One of

- `not\_controllable`: cannot be controlled by any extension
- `controlled\_by\_other\_extensions`: controlled by extensions with higher precedence
- `controllable\_by\_this\_extension`: can be controlled by this extension
- `controlled\_by\_this\_extension`: controlled by this extension

#### Enum

"not\_controllable"

"controlled\_by\_other\_extensions"

"controllable\_by\_this\_extension"

"controlled\_by\_this\_extension"