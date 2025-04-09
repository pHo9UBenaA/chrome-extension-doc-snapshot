# chrome.runtime

## Description

Use the `chrome.runtime` API to retrieve the service worker, return details about the manifest, and listen for and respond to events in the extension lifecycle. You can also use this API to convert the relative path of URLs to fully-qualified URLs.

Most members of this API do **not** require any permissions. This permission is needed for [`connectNative()`](#method-connectNative), [`sendNativeMessage()`](#method-sendNativeMessage) and [`onNativeConnect`](#event-onConnectNative).

The following example shows how to declare the `"nativeMessaging"` permission in the manifest:

manifest.json:

```
{
  "name": "My extension",
  ...
  "permissions": [
    "nativeMessaging"
  ],
  ...
}
```

## Concepts and usage

The Runtime API provides methods to support a number of areas that your extensions can use:

Message passing

Your extension can communicate with different contexts within your extension and also with other extensions using these methods and events: [`connect()`](#method-connect), [`onConnect`](#event-onConnect), [`onConnectExternal`](#event-onConnectExternal), [`sendMessage()`](#method-sendMessage), [`onMessage`](#event-onMessage) and [`onMessageExternal`](#event-onMessageExternal). In addition, your extension can pass messages to native applications on the user's device using [`connectNative()`](#method-connectNative) and [`sendNativeMessage()`](#method-sendNativeMessage).

**Note:** See [Message Passing](/docs/extensions/mv3/messaging) for an overview of the subject.

Accessing extension and platform metadata

These methods let you retrieve several specific pieces of metadata about the extension and the platform. Methods in this category include [`getManifest()`](#method-getManifest), and [`getPlatformInfo()`](#method-getPlatformInfo).

Managing extension lifecycle and options

These properties let you perform some meta-operations on the extension, and display the options page. Methods and events in this category include [`onInstalled`](#event-onInstalled), [`onStartup`](#event-onStartup), [`openOptionsPage()`](#method-openOptionsPage), [`reload()`](#method-reload), [`requestUpdateCheck()`](#method-requestUpdateCheck), and [`setUninstallURL()`](#method-setUninstallURL).

Helper utilities

These methods provide utility such as the conversion of internal resource representations to external formats. Methods in this category include [`getURL()`](#method-getURL).

Kiosk mode utilities

These methods are available only on ChromeOS, and exist mainly to support kiosk implementations. Methods in this category include [`restart()`](#method-restart) and [`restartAfterDelay()`\`](#method-restartAfterDelay).

### Unpacked extension behavior

When an [unpacked](/docs/extensions/get-started/tutorial/hello-world#load-unpacked) extension is [reloaded](/docs/extensions/get-started/tutorial/hello-world#reload), this is treated as an update. This means that the [`chrome.runtime.onInstalled`](#event-onInstalled) event will fire with the `"update"` reason. This includes when the extension is reloaded with [`chrome.runtime.reload()`](#method-reload).

## Use cases

### Add an image to a web page

For a web page to access an asset hosted on another domain, it must specify the resource's full URL (e.g. `<img src="https://example.com/logo.png">`). The same is true to include an extension asset on a web page. The two differences are that the extension's assets must be exposed as [web accessible resources](/docs/extensions/mv3/manifest/web_accessible_resources) and that typically content scripts are responsible for injecting extension assets.

In this example, the extension will add `logo.png` to the page that the [content script](/docs/extensions/develop/concepts/content-scripts) is being [injected](/docs/extensions/develop/concepts/content-scripts#functionality) into by using `runtime.getURL()` to create a fully-qualified URL. But first, the asset must be declared as a web accessible resource in the manifest.

manifest.json:

```
{
  ...
  "web_accessible_resources": [
    {
      "resources": [ "logo.png" ],
      "matches": [ "https://*/*" ]
    }
  ],
  ...
}
```

content.js:

```
{ // Block used to avoid setting global variables
  const img = document.createElement('img');
  img.src = chrome.runtime.getURL('logo.png');
  document.body.append(img);
}
```

### Send data from a content script to the service worker

Its common for an extension's content scripts to need data managed by another part of the extension, like the service worker. Much like two browser windows opened to the same web page, these two contexts cannot directly access each other's values. Instead, the extension can use [message passing](/docs/extensions/develop/concepts/messaging) to coordinate across these different contexts.

In this example, the content script needs some data from the extension's service worker to initialize its UI. To get this data, it passes the developer-defined `get-user-data` message to the service worker, and it responds with a copy of the user's information.

content.js:

```
// 1. Send a message to the service worker requesting the user's data
chrome.runtime.sendMessage('get-user-data', (response) => {
  // 3. Got an asynchronous response with the data from the service worker
  console.log('received user data', response);
  initializeUI(response);
});
```

service-worker.js:

```
// Example of a simple user data object
const user = {
  username: 'demo-user'
};

chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  // 2. A page requested user data, respond with a copy of `user`
  if (message === 'get-user-data') {
    sendResponse(user);
  }
});
```

### Gather feedback on uninstall

Many extensions use post-uninstall surveys to understand how the extension could better serve its users and improve retention. The following example shows how to add this functionality.

background.js:

```
chrome.runtime.onInstalled.addListener(details => {
  if (details.reason === chrome.runtime.OnInstalledReason.INSTALL) {
    chrome.runtime.setUninstallURL('https://example.com/extension-survey');
  }
});
```

## Examples

See the [Manifest V3 - Web Accessible Resources demo](https://github.com/GoogleChrome/chrome-extensions-samples/tree/main/api-samples/web-accessible-resources) for more Runtime API examples.

## Types

### ContextFilter

Chrome 114+

A filter to match against certain extension contexts. Matching contexts must match all specified filters; any filter that is not specified matches all available contexts. Thus, a filter of \`{}\` will match all available contexts.

#### Properties

- contextIds
  
  string\[] optional
- contextTypes
  
  [ContextType](#type-ContextType)\[] optional
- documentIds
  
  string\[] optional
- documentOrigins
  
  string\[] optional
- documentUrls
  
  string\[] optional
- frameIds
  
  number\[] optional
- incognito
  
  boolean optional
- tabIds
  
  number\[] optional
- windowIds
  
  number\[] optional

### ContextType

Chrome 114+

#### Enum

"TAB"  
Specifies the context type as a tab

"POPUP"  
Specifies the context type as an extension popup window

"BACKGROUND"  
Specifies the context type as a service worker.

"OFFSCREEN\_DOCUMENT"  
Specifies the context type as an offscreen document.

"SIDE\_PANEL"  
Specifies the context type as a side panel.

"DEVELOPER\_TOOLS"  
Specifies the context type as developer tools.

### ExtensionContext

Chrome 114+

A context hosting extension content.

#### Properties

- contextId
  
  string
  
  A unique identifier for this context
- contextType
  
  [ContextType](#type-ContextType)
  
  The type of context this corresponds to.
- documentId
  
  string optional
  
  A UUID for the document associated with this context, or undefined if this context is hosted not in a document.
- documentOrigin
  
  string optional
  
  The origin of the document associated with this context, or undefined if the context is not hosted in a document.
- documentUrl
  
  string optional
  
  The URL of the document associated with this context, or undefined if the context is not hosted in a document.
- frameId
  
  number
  
  The ID of the frame for this context, or -1 if this context is not hosted in a frame.
- incognito
  
  boolean
  
  Whether the context is associated with an incognito profile.
- tabId
  
  number
  
  The ID of the tab for this context, or -1 if this context is not hosted in a tab.
- windowId
  
  number
  
  The ID of the window for this context, or -1 if this context is not hosted in a window.

### MessageSender

An object containing information about the script context that sent a message or request.

#### Properties

- documentId
  
  string optional
  
  Chrome 106+
  
  A UUID of the document that opened the connection.
- documentLifecycle
  
  string optional
  
  Chrome 106+
  
  The lifecycle the document that opened the connection is in at the time the port was created. Note that the lifecycle state of the document may have changed since port creation.
- frameId
  
  number optional
  
  The [frame](https://developer.chrome.com/docs/extensions/reference/webNavigation/#frame_ids) that opened the connection. 0 for top-level frames, positive for child frames. This will only be set when `tab` is set.
- id
  
  string optional
  
  The ID of the extension that opened the connection, if any.
- nativeApplication
  
  string optional
  
  Chrome 74+
  
  The name of the native application that opened the connection, if any.
- origin
  
  string optional
  
  Chrome 80+
  
  The origin of the page or frame that opened the connection. It can vary from the url property (e.g., about:blank) or can be opaque (e.g., sandboxed iframes). This is useful for identifying if the origin can be trusted if we can't immediately tell from the URL.
- tab
  
  [Tab](https://developer.chrome.com/docs/extensions/reference/tabs/#type-Tab) optional
  
  The [`tabs.Tab`](https://developer.chrome.com/docs/extensions/reference/tabs/#type-Tab) which opened the connection, if any. This property will **only** be present when the connection was opened from a tab (including content scripts), and **only** if the receiver is an extension, not an app.
- tlsChannelId
  
  string optional
  
  The TLS channel ID of the page or frame that opened the connection, if requested by the extension, and if available.
- url
  
  string optional
  
  The URL of the page or frame that opened the connection. If the sender is in an iframe, it will be iframe's URL not the URL of the page which hosts it.

### OnInstalledReason

Chrome 44+

The reason that this event is being dispatched.

#### Enum

"install"  
Specifies the event reason as an installation.

"update"  
Specifies the event reason as an extension update.

"chrome\_update"  
Specifies the event reason as a Chrome update.

"shared\_module\_update"  
Specifies the event reason as an update to a shared module.

### OnRestartRequiredReason

Chrome 44+

The reason that the event is being dispatched. 'app\_update' is used when the restart is needed because the application is updated to a newer version. 'os\_update' is used when the restart is needed because the browser/OS is updated to a newer version. 'periodic' is used when the system runs for more than the permitted uptime set in the enterprise policy.

#### Enum

"app\_update"  
Specifies the event reason as an update to the app.

"os\_update"  
Specifies the event reason as an update to the operating system.

"periodic"  
Specifies the event reason as a periodic restart of the app.

### PlatformArch

Chrome 44+

The machine's processor architecture.

#### Enum

"arm"  
Specifies the processer architecture as arm.

"arm64"  
Specifies the processer architecture as arm64.

"x86-32"  
Specifies the processer architecture as x86-32.

"x86-64"  
Specifies the processer architecture as x86-64.

"mips"  
Specifies the processer architecture as mips.

"mips64"  
Specifies the processer architecture as mips64.

### PlatformInfo

An object containing information about the current platform.

#### Properties

- arch
  
  [PlatformArch](#type-PlatformArch)
  
  The machine's processor architecture.
- nacl\_arch
  
  [PlatformNaclArch](#type-PlatformNaclArch)
  
  The native client architecture. This may be different from arch on some platforms.
- os
  
  [PlatformOs](#type-PlatformOs)
  
  The operating system Chrome is running on.

### PlatformNaclArch

Chrome 44+

The native client architecture. This may be different from arch on some platforms.

#### Enum

"arm"  
Specifies the native client architecture as arm.

"x86-32"  
Specifies the native client architecture as x86-32.

"x86-64"  
Specifies the native client architecture as x86-64.

"mips"  
Specifies the native client architecture as mips.

"mips64"  
Specifies the native client architecture as mips64.

### PlatformOs

Chrome 44+

The operating system Chrome is running on.

#### Enum

"mac"  
Specifies the MacOS operating system.

"win"  
Specifies the Windows operating system.

"android"  
Specifies the Android operating system.

"cros"  
Specifies the Chrome operating system.

"linux"  
Specifies the Linux operating system.

"openbsd"  
Specifies the OpenBSD operating system.

"fuchsia"  
Specifies the Fuchsia operating system.

### Port

An object which allows two way communication with other pages. See [Long-lived connections](https://developer.chrome.com/docs/extensions/messaging#connect) for more information.

#### Properties

- name
  
  string
  
  The name of the port, as specified in the call to [`runtime.connect`](#method-connect).
- onDisconnect
  
  Event&lt;functionvoidvoid&gt;
  
  Fired when the port is disconnected from the other end(s). [`runtime.lastError`](#property-lastError) may be set if the port was disconnected by an error. If the port is closed via [disconnect](#method-Port-disconnect), then this event is *only* fired on the other end. This event is fired at most once (see also [Port lifetime](https://developer.chrome.com/docs/extensions/messaging#port-lifetime)).
  
  The `onDisconnect.addListener` function looks like:
  
  ```
  (callback: function) => {...}
  ```
  
  - callback
    
    function
    
    The `callback` parameter looks like:
    
    ```
    (port: Port) => void
    ```
    
    - port
      
      [Port](#type-Port)
- onMessage
  
  Event&lt;functionvoidvoid&gt;
  
  This event is fired when [postMessage](#method-Port-postMessage) is called by the other end of the port.
  
  The `onMessage.addListener` function looks like:
  
  ```
  (callback: function) => {...}
  ```
  
  - callback
    
    function
    
    The `callback` parameter looks like:
    
    ```
    (message: any, port: Port) => void
    ```
    
    - message
      
      any
    - port
      
      [Port](#type-Port)
- sender
  
  [MessageSender](#type-MessageSender) optional
  
  This property will **only** be present on ports passed to [onConnect](#event-onConnect) / [onConnectExternal](#event-onConnectExternal) / [onConnectNative](#event-onConnectExternal) listeners.
- disconnect
  
  void
  
  Immediately disconnect the port. Calling `disconnect()` on an already-disconnected port has no effect. When a port is disconnected, no new events will be dispatched to this port.
  
  The `disconnect` function looks like:
  
  ```
  () => {...}
  ```
- postMessage
  
  void
  
  Send a message to the other end of the port. If the port is disconnected, an error is thrown.
  
  The `postMessage` function looks like:
  
  ```
  (message: any) => {...}
  ```
  
  - message
    
    any
    
    Chrome 52+
    
    The message to send. This object should be JSON-ifiable.

### RequestUpdateCheckStatus

Chrome 44+

Result of the update check.

#### Enum

"throttled"  
Specifies that the status check has been throttled. This can occur after repeated checks within a short amount of time.

"no\_update"  
Specifies that there are no available updates to install.

"update\_available"  
Specifies that there is an available update to install.

## Properties

### id

The ID of the extension/app.

#### Type

string

### lastError

Populated with an error message if calling an API function fails; otherwise undefined. This is only defined within the scope of that function's callback. If an error is produced, but `runtime.lastError` is not accessed within the callback, a message is logged to the console listing the API function that produced the error. API functions that return promises do not set this property.

#### Type

object

#### Properties

- message
  
  string optional
  
  Details about the error which occurred.

## Methods

### connect()

```
chrome.runtime.connect(
  extensionId?: string,
  connectInfo?: object,
)
```

Attempts to connect listeners within an extension (such as the background page), or other extensions/apps. This is useful for content scripts connecting to their extension processes, inter-app/extension communication, and [web messaging](https://developer.chrome.com/docs/extensions/manifest/externally_connectable). Note that this does not connect to any listeners in a content script. Extensions may connect to content scripts embedded in tabs via [`tabs.connect`](https://developer.chrome.com/docs/extensions/reference/tabs/#method-connect).

#### Parameters

- extensionId
  
  string optional
  
  The ID of the extension to connect to. If omitted, a connection will be attempted with your own extension. Required if sending messages from a web page for [web messaging](https://developer.chrome.com/docs/extensions/reference/manifest/externally-connectable).
- connectInfo
  
  object optional
  
  - includeTlsChannelId
    
    boolean optional
    
    Whether the TLS channel ID will be passed into onConnectExternal for processes that are listening for the connection event.
  - name
    
    string optional
    
    Will be passed into onConnect for processes that are listening for the connection event.

#### Returns

- [Port](#type-Port)
  
  Port through which messages can be sent and received. The port's [onDisconnect](#type-Port) event is fired if the extension does not exist.

### connectNative()

```
chrome.runtime.connectNative(
  application: string,
)
```

Connects to a native application in the host machine. This method requires the `"nativeMessaging"` permission. See [Native Messaging](https://developer.chrome.com/extensions/develop/concepts/native-messaging) for more information.

#### Parameters

- application
  
  string
  
  The name of the registered application to connect to.

#### Returns

- [Port](#type-Port)
  
  Port through which messages can be sent and received with the application

### getBackgroundPage()

Promise Foreground only Deprecated since Chrome 133

```
chrome.runtime.getBackgroundPage(
  callback?: function,
)
```

Background pages do not exist in MV3 extensions.

Retrieves the JavaScript 'window' object for the background page running inside the current extension/app. If the background page is an event page, the system will ensure it is loaded before calling the callback. If there is no background page, an error is set.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (backgroundPage?: Window) => void
  ```
  
  - backgroundPage
    
    Window optional
    
    The JavaScript 'window' object for the background page.

#### Returns

- Promise&lt;Window | undefined&gt;
  
  Chrome 99+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### getContexts()

Promise Chrome 116+ MV3+

```
chrome.runtime.getContexts(
  filter: ContextFilter,
  callback?: function,
)
```

Fetches information about active contexts associated with this extension

#### Parameters

- filter
  
  [ContextFilter](#type-ContextFilter)
  
  A filter to find matching contexts. A context matches if it matches all specified fields in the filter. Any unspecified field in the filter matches all contexts.
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (contexts: ExtensionContext[]) => void
  ```
  
  - contexts
    
    [ExtensionContext](#type-ExtensionContext)\[]
    
    The matching contexts, if any.

#### Returns

- Promise&lt;[ExtensionContext](#type-ExtensionContext)\[]&gt;
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### getManifest()

```
chrome.runtime.getManifest()
```

Returns details about the app or extension from the manifest. The object returned is a serialization of the full [manifest file](https://developer.chrome.com/docs/extensions/reference/manifest).

#### Returns

- object
  
  The manifest details.

### getPackageDirectoryEntry()

Promise Foreground only

```
chrome.runtime.getPackageDirectoryEntry(
  callback?: function,
)
```

Returns a DirectoryEntry for the package directory.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (directoryEntry: DirectoryEntry) => void
  ```
  
  - directoryEntry
    
    DirectoryEntry

#### Returns

- Promise&lt;DirectoryEntry&gt;
  
  Chrome 122+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### getPlatformInfo()

Promise

```
chrome.runtime.getPlatformInfo(
  callback?: function,
)
```

Returns information about the current platform.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (platformInfo: PlatformInfo) => void
  ```
  
  - platformInfo
    
    [PlatformInfo](#type-PlatformInfo)

#### Returns

- Promise&lt;[PlatformInfo](#type-PlatformInfo)&gt;
  
  Chrome 99+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### getURL()

```
chrome.runtime.getURL(
  path: string,
)
```

Converts a relative path within an app/extension install directory to a fully-qualified URL.

#### Parameters

- path
  
  string
  
  A path to a resource within an app/extension expressed relative to its install directory.

#### Returns

- string
  
  The fully-qualified URL to the resource.

### openOptionsPage()

Promise

```
chrome.runtime.openOptionsPage(
  callback?: function,
)
```

Open your Extension's options page, if possible.

The precise behavior may depend on your manifest's [`options_ui`](https://developer.chrome.com/docs/extensions/develop/ui/options-page#embedded_options) or [`options_page`](https://developer.chrome.com/docs/extensions/develop/ui/options-page#full_page) key, or what Chrome happens to support at the time. For example, the page may be opened in a new tab, within chrome://extensions, within an App, or it may just focus an open options page. It will never cause the caller page to reload.

If your Extension does not declare an options page, or Chrome failed to create one for some other reason, the callback will set [`lastError`](#property-lastError).

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

#### Returns

- Promise&lt;void&gt;
  
  Chrome 99+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### reload()

```
chrome.runtime.reload()
```

Reloads the app or extension. This method is not supported in kiosk mode. For kiosk mode, use chrome.runtime.restart() method.

### requestUpdateCheck()

Promise

```
chrome.runtime.requestUpdateCheck(
  callback?: function,
)
```

Requests an immediate update check be done for this app/extension.

**Important**: Most extensions/apps should **not** use this method, since Chrome already does automatic checks every few hours, and you can listen for the [`runtime.onUpdateAvailable`](#event-onUpdateAvailable) event without needing to call requestUpdateCheck.

This method is only appropriate to call in very limited circumstances, such as if your extension talks to a backend service, and the backend service has determined that the client extension version is very far out of date and you'd like to prompt a user to update. Most other uses of requestUpdateCheck, such as calling it unconditionally based on a repeating timer, probably only serve to waste client, network, and server resources.

Note: When called with a callback, instead of returning an object this function will return the two properties as separate arguments passed to the callback.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (result: object) => void
  ```
  
  - result
    
    object
    
    Chrome 109+
    
    RequestUpdateCheckResult object that holds the status of the update check and any details of the result if there is an update available
    
    - status
      
      [RequestUpdateCheckStatus](#type-RequestUpdateCheckStatus)
      
      Result of the update check.
    - version
      
      string optional
      
      If an update is available, this contains the version of the available update.

#### Returns

- Promise&lt;object&gt;
  
  Chrome 109+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### restart()

```
chrome.runtime.restart()
```

Restart the ChromeOS device when the app runs in kiosk mode. Otherwise, it's no-op.

### restartAfterDelay()

Promise Chrome 53+

```
chrome.runtime.restartAfterDelay(
  seconds: number,
  callback?: function,
)
```

Restart the ChromeOS device when the app runs in kiosk mode after the given seconds. If called again before the time ends, the reboot will be delayed. If called with a value of -1, the reboot will be cancelled. It's a no-op in non-kiosk mode. It's only allowed to be called repeatedly by the first extension to invoke this API.

#### Parameters

- seconds
  
  number
  
  Time to wait in seconds before rebooting the device, or -1 to cancel a scheduled reboot.
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

#### Returns

- Promise&lt;void&gt;
  
  Chrome 99+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### sendMessage()

Promise

```
chrome.runtime.sendMessage(
  extensionId?: string,
  message: any,
  options?: object,
  callback?: function,
)
```

Sends a single message to event listeners within your extension or a different extension/app. Similar to [`runtime.connect`](#method-connect) but only sends a single message, with an optional response. If sending to your extension, the [`runtime.onMessage`](#event-onMessage) event will be fired in every frame of your extension (except for the sender's frame), or [`runtime.onMessageExternal`](#event-onMessageExternal), if a different extension. Note that extensions cannot send messages to content scripts using this method. To send messages to content scripts, use [`tabs.sendMessage`](https://developer.chrome.com/docs/extensions/reference/tabs/#method-sendMessage).

#### Parameters

- extensionId
  
  string optional
  
  The ID of the extension to send the message to. If omitted, the message will be sent to your own extension/app. Required if sending messages from a web page for [web messaging](https://developer.chrome.com/docs/extensions/manifest/externally_connectable).
- message
  
  any
  
  The message to send. This message should be a JSON-ifiable object.
- options
  
  object optional
  
  - includeTlsChannelId
    
    boolean optional
    
    Whether the TLS channel ID will be passed into onMessageExternal for processes that are listening for the connection event.
- callback
  
  function optional
  
  Chrome 99+
  
  The `callback` parameter looks like:
  
  ```
  (response: any) => void
  ```
  
  - response
    
    any
    
    The JSON response object sent by the handler of the message. If an error occurs while connecting to the extension, the callback will be called with no arguments and [`runtime.lastError`](#property-lastError) will be set to the error message.

#### Returns

- Promise&lt;any&gt;
  
  Chrome 99+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### sendNativeMessage()

Promise

```
chrome.runtime.sendNativeMessage(
  application: string,
  message: object,
  callback?: function,
)
```

Send a single message to a native application. This method requires the `"nativeMessaging"` permission.

#### Parameters

- application
  
  string
  
  The name of the native messaging host.
- message
  
  object
  
  The message that will be passed to the native messaging host.
- callback
  
  function optional
  
  Chrome 99+
  
  The `callback` parameter looks like:
  
  ```
  (response: any) => void
  ```
  
  - response
    
    any
    
    The response message sent by the native messaging host. If an error occurs while connecting to the native messaging host, the callback will be called with no arguments and [`runtime.lastError`](#property-lastError) will be set to the error message.

#### Returns

- Promise&lt;any&gt;
  
  Chrome 99+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### setUninstallURL()

Promise

```
chrome.runtime.setUninstallURL(
  url: string,
  callback?: function,
)
```

Sets the URL to be visited upon uninstallation. This may be used to clean up server-side data, do analytics, and implement surveys. Maximum 1023 characters.

#### Parameters

- url
  
  string
  
  URL to be opened after the extension is uninstalled. This URL must have an http: or https: scheme. Set an empty string to not open a new tab upon uninstallation.
- callback
  
  function optional
  
  Chrome 45+
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

#### Returns

- Promise&lt;void&gt;
  
  Chrome 99+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

## Events

### onBrowserUpdateAvailable

Deprecated

```
chrome.runtime.onBrowserUpdateAvailable.addListener(
  callback: function,
)
```

Please use [`runtime.onRestartRequired`](#event-onRestartRequired).

Fired when a Chrome update is available, but isn't installed immediately because a browser restart is required.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

### onConnect

```
chrome.runtime.onConnect.addListener(
  callback: function,
)
```

Fired when a connection is made from either an extension process or a content script (by [`runtime.connect`](#method-connect)).

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (port: Port) => void
  ```
  
  - port
    
    [Port](#type-Port)

### onConnectExternal

```
chrome.runtime.onConnectExternal.addListener(
  callback: function,
)
```

Fired when a connection is made from another extension (by [`runtime.connect`](#method-connect)), or from an externally connectable web site.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (port: Port) => void
  ```
  
  - port
    
    [Port](#type-Port)

### onConnectNative

Chrome 76+

```
chrome.runtime.onConnectNative.addListener(
  callback: function,
)
```

Fired when a connection is made from a native application. This event requires the `"nativeMessaging"` permission. It is only supported on Chrome OS.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (port: Port) => void
  ```
  
  - port
    
    [Port](#type-Port)

### onInstalled

```
chrome.runtime.onInstalled.addListener(
  callback: function,
)
```

Fired when the extension is first installed, when the extension is updated to a new version, and when Chrome is updated to a new version.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (details: object) => void
  ```
  
  - details
    
    object
    
    - id
      
      string optional
      
      Indicates the ID of the imported shared module extension which updated. This is present only if 'reason' is 'shared\_module\_update'.
    - previousVersion
      
      string optional
      
      Indicates the previous version of the extension, which has just been updated. This is present only if 'reason' is 'update'.
    - reason
      
      [OnInstalledReason](#type-OnInstalledReason)
      
      The reason that this event is being dispatched.

### onMessage

```
chrome.runtime.onMessage.addListener(
  callback: function,
)
```

Fired when a message is sent from either an extension process (by [`runtime.sendMessage`](#method-sendMessage)) or a content script (by [`tabs.sendMessage`](https://developer.chrome.com/docs/extensions/reference/tabs/#method-sendMessage)).

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (message: any, sender: MessageSender, sendResponse: function) => boolean | undefined
  ```
  
  - message
    
    any
  - sender
    
    [MessageSender](#type-MessageSender)
  - sendResponse
    
    function
    
    The `sendResponse` parameter looks like:
    
    ```
    () => void
    ```
  
  <!--THE END-->
  
  - returns
    
    boolean | undefined

### onMessageExternal

```
chrome.runtime.onMessageExternal.addListener(
  callback: function,
)
```

Fired when a message is sent from another extension (by [`runtime.sendMessage`](#method-sendMessage)). Cannot be used in a content script.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (message: any, sender: MessageSender, sendResponse: function) => boolean | undefined
  ```
  
  - message
    
    any
  - sender
    
    [MessageSender](#type-MessageSender)
  - sendResponse
    
    function
    
    The `sendResponse` parameter looks like:
    
    ```
    () => void
    ```
  
  <!--THE END-->
  
  - returns
    
    boolean | undefined

### onRestartRequired

```
chrome.runtime.onRestartRequired.addListener(
  callback: function,
)
```

Fired when an app or the device that it runs on needs to be restarted. The app should close all its windows at its earliest convenient time to let the restart to happen. If the app does nothing, a restart will be enforced after a 24-hour grace period has passed. Currently, this event is only fired for Chrome OS kiosk apps.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (reason: OnRestartRequiredReason) => void
  ```
  
  - reason
    
    [OnRestartRequiredReason](#type-OnRestartRequiredReason)

### onStartup

```
chrome.runtime.onStartup.addListener(
  callback: function,
)
```

Fired when a profile that has this extension installed first starts up. This event is not fired when an incognito profile is started, even if this extension is operating in 'split' incognito mode.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

### onSuspend

```
chrome.runtime.onSuspend.addListener(
  callback: function,
)
```

Sent to the event page just before it is unloaded. This gives the extension opportunity to do some clean up. Note that since the page is unloading, any asynchronous operations started while handling this event are not guaranteed to complete. If more activity for the event page occurs before it gets unloaded the onSuspendCanceled event will be sent and the page won't be unloaded.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

### onSuspendCanceled

```
chrome.runtime.onSuspendCanceled.addListener(
  callback: function,
)
```

Sent after onSuspend to indicate that the app won't be unloaded after all.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

### onUpdateAvailable

```
chrome.runtime.onUpdateAvailable.addListener(
  callback: function,
)
```

Fired when an update is available, but isn't installed immediately because the app is currently running. If you do nothing, the update will be installed the next time the background page gets unloaded, if you want it to be installed sooner you can explicitly call chrome.runtime.reload(). If your extension is using a persistent background page, the background page of course never gets unloaded, so unless you call chrome.runtime.reload() manually in response to this event the update will not get installed until the next time Chrome itself restarts. If no handlers are listening for this event, and your extension has a persistent background page, it behaves as if chrome.runtime.reload() is called in response to this event.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (details: object) => void
  ```
  
  - details
    
    object
    
    - version
      
      string
      
      The version number of the available update.

### onUserScriptConnect

Chrome 115+ MV3+

```
chrome.runtime.onUserScriptConnect.addListener(
  callback: function,
)
```

Fired when a connection is made from a user script from this extension.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (port: Port) => void
  ```
  
  - port
    
    [Port](#type-Port)

### onUserScriptMessage

Chrome 115+ MV3+

```
chrome.runtime.onUserScriptMessage.addListener(
  callback: function,
)
```

Fired when a message is sent from a user script associated with the same extension.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (message: any, sender: MessageSender, sendResponse: function) => boolean | undefined
  ```
  
  - message
    
    any
  - sender
    
    [MessageSender](#type-MessageSender)
  - sendResponse
    
    function
    
    The `sendResponse` parameter looks like:
    
    ```
    () => void
    ```
  
  <!--THE END-->
  
  - returns
    
    boolean | undefined