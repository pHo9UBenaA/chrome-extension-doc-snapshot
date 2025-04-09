# chrome.scripting

## Description

Use the `chrome.scripting` API to execute script in different contexts.

## Permissions

`scripting`

## Availability

Chrome 88+ MV3+

## Manifest

To use the `chrome.scripting` API, declare the `"scripting"` permission in the [manifest](/docs/extensions/reference/manifest) plus the host permissions for the pages to inject scripts into. Use the [`"host_permissions"`](/docs/extensions/develop/concepts/declare-permissions) key or the [`"activeTab"`](/docs/extensions/develop/concepts/activeTab) permission, which grants temporary host permissions. The following example uses the activeTab permission.

```
{
  "name": "Scripting Extension",
  "manifest_version": 3,
  "permissions": ["scripting", "activeTab"],
  ...
}
```

## Concepts and usage

You can use the `chrome.scripting` API to inject JavaScript and CSS into websites. This is similar to what you can do with [content scripts](/docs/extensions/develop/concepts/content-scripts). But by using the [`chrome.scripting`](/docs/extensions/reference/scripting) namespace, extensions can make decisions at runtime.

### Injection targets

You can use the `target` parameter to specify a target to inject JavaScript or CSS into.

The only required field is `tabId`. By default, an injection will run in the main frame of the specified tab.

```
function getTabId() { ... }

chrome.scripting
    .executeScript({
      target : {tabId : getTabId()},
      files : [ "script.js" ],
    })
    .then(() => console.log("script injected"));
```

To run in all frames of the specified tab, you can set the `allFrames` boolean to `true`.

```
function getTabId() { ... }

chrome.scripting
    .executeScript({
      target : {tabId : getTabId(), allFrames : true},
      files : [ "script.js" ],
    })
    .then(() => console.log("script injected in all frames"));
```

You can also inject into specific frames of a tab by specifying individual frame IDs. For more information on frame IDs, see the [`chrome.webNavigation` API](/docs/extensions/reference/api/webNavigation).

```
function getTabId() { ... }

chrome.scripting
    .executeScript({
      target : {tabId : getTabId(), frameIds : [ frameId1, frameId2 ]},
      files : [ "script.js" ],
    })
    .then(() => console.log("script injected on target frames"));
```

**Note:** You cannot specify both the `"frameIds"` and `"allFrames"` properties.

### Injected code

Extensions can specify the code to be injected either via an external file or a runtime variable.

#### Files

Files are specified as strings that are paths relative to the extension's root directory. The following code will inject the file `script.js` into the main frame of the tab.

```
function getTabId() { ... }

chrome.scripting
    .executeScript({
      target : {tabId : getTabId()},
      files : [ "script.js" ],
    })
    .then(() => console.log("injected script file"));
```

#### Runtime functions

When injecting JavaScript with `scripting.executeScript()`, you can specify a function to be executed instead of a file. This function should be a function variable available to the current extension context.

```
function getTabId() { ... }
function getTitle() { return document.title; }

chrome.scripting
    .executeScript({
      target : {tabId : getTabId()},
      func : getTitle,
    })
    .then(() => console.log("injected a function"));
```

```
function getTabId() { ... }
function getUserColor() { ... }

function changeBackgroundColor() {
  document.body.style.backgroundColor = getUserColor();
}

chrome.scripting
    .executeScript({
      target : {tabId : getTabId()},
      func : changeBackgroundColor,
    })
    .then(() => console.log("injected a function"));
```

You can work around this by using the `args` property:

```
function getTabId() { ... }
function getUserColor() { ... }
function changeBackgroundColor(backgroundColor) {
  document.body.style.backgroundColor = backgroundColor;
}

chrome.scripting
    .executeScript({
      target : {tabId : getTabId()},
      func : changeBackgroundColor,
      args : [ getUserColor() ],
    })
    .then(() => console.log("injected a function"));
```

#### Runtime strings

If injecting CSS within a page, you can also specify a string to be used in the `css` property. This option is only available for `scripting.insertCSS()`; you can't execute a string using `scripting.executeScript()`.

```
function getTabId() { ... }
const css = "body { background-color: red; }";

chrome.scripting
    .insertCSS({
      target : {tabId : getTabId()},
      css : css,
    })
    .then(() => console.log("CSS injected"));
```

### Handle the results

The results of executing JavaScript are passed to the extension. A single result is included per-frame. The main frame is guaranteed to be the first index in the resulting array; all other frames are in a non-deterministic order.

```
function getTabId() { ... }
function getTitle() { return document.title; }

chrome.scripting
    .executeScript({
      target : {tabId : getTabId(), allFrames : true},
      func : getTitle,
    })
    .then(injectionResults => {
      for (const {frameId, result} of injectionResults) {
        console.log(`Frame ${frameId} result:`, result);
      }
    });
```

`scripting.insertCSS()` does not return any results.

#### Promises

If the resulting value of the script execution is a promise, Chrome will wait for the promise to settle and return the resulting value.

```
function getTabId() { ... }
async function addIframe() {
  const iframe = document.createElement("iframe");
  const loadComplete =
      new Promise(resolve => iframe.addEventListener("load", resolve));
  iframe.src = "https://example.com";
  document.body.appendChild(iframe);
  await loadComplete;
  return iframe.contentWindow.document.title;
}

chrome.scripting
    .executeScript({
      target : {tabId : getTabId(), allFrames : true},
      func : addIframe,
    })
    .then(injectionResults => {
      for (const frameResult of injectionResults) {
        const {frameId, result} = frameResult;
        console.log(`Frame ${frameId} result:`, result);
      }
    });
```

## Examples

### Unregister all dynamic content scripts

The following snippet contains a function that unregisters all dynamic content scripts the extension has previously registered.

```
async function unregisterAllDynamicContentScripts() {
  try {
    const scripts = await chrome.scripting.getRegisteredContentScripts();
    const scriptIds = scripts.map(script => script.id);
    return chrome.scripting.unregisterContentScripts({ ids: scriptIds });
  } catch (error) {
    const message = [
      "An unexpected error occurred while",
      "unregistering dynamic content scripts.",
    ].join(" ");
    throw new Error(message, {cause : error});
  }
}
```

**Key point:** Unregistering content scripts will not remove scripts or styles that have already been injected.

To try the `chrome.scripting` API, install the [scripting sample](https://github.com/GoogleChrome/chrome-extensions-samples/tree/main/api-samples/scripting) from the [Chrome extension samples](https://github.com/GoogleChrome/chrome-extensions-samples) repository.

## Types

### ContentScriptFilter

Chrome 96+

#### Properties

- ids
  
  string\[] optional
  
  If specified, [`getRegisteredContentScripts`](#method-getRegisteredContentScripts) will only return scripts with an id specified in this list.

### CSSInjection

#### Properties

- css
  
  string optional
  
  A string containing the CSS to inject. Exactly one of `files` and `css` must be specified.
- files
  
  string\[] optional
  
  The path of the CSS files to inject, relative to the extension's root directory. Exactly one of `files` and `css` must be specified.
- origin
  
  [StyleOrigin](#type-StyleOrigin) optional
  
  The style origin for the injection. Defaults to `'AUTHOR'`.
- target
  
  [InjectionTarget](#type-InjectionTarget)
  
  Details specifying the target into which to insert the CSS.

### ExecutionWorld

Chrome 95+

The JavaScript world for a script to execute within.

#### Enum

"ISOLATED"  
Specifies the isolated world, which is the execution environment unique to this extension.

"MAIN"  
Specifies the main world of the DOM, which is the execution environment shared with the host page's JavaScript.

### InjectionResult

#### Properties

- documentId
  
  string
  
  Chrome 106+
  
  The document associated with the injection.
- frameId
  
  number
  
  Chrome 90+
  
  The frame associated with the injection.
- result
  
  any optional
  
  The result of the script execution.

### InjectionTarget

#### Properties

- allFrames
  
  boolean optional
  
  Whether the script should inject into all frames within the tab. Defaults to false. This must not be true if `frameIds` is specified.
- documentIds
  
  string\[] optional
  
  Chrome 106+
  
  The [IDs](https://developer.chrome.com/docs/extensions/reference/webNavigation/#document_ids) of specific documentIds to inject into. This must not be set if `frameIds` is set.
- frameIds
  
  number\[] optional
  
  The [IDs](https://developer.chrome.com/docs/extensions/reference/webNavigation/#frame_ids) of specific frames to inject into.
- tabId
  
  number
  
  The ID of the tab into which to inject.

### RegisteredContentScript

Chrome 96+

#### Properties

- allFrames
  
  boolean optional
  
  If specified true, it will inject into all frames, even if the frame is not the top-most frame in the tab. Each frame is checked independently for URL requirements; it will not inject into child frames if the URL requirements are not met. Defaults to false, meaning that only the top frame is matched.
- css
  
  string\[] optional
  
  The list of CSS files to be injected into matching pages. These are injected in the order they appear in this array, before any DOM is constructed or displayed for the page.
- excludeMatches
  
  string\[] optional
  
  Excludes pages that this content script would otherwise be injected into. See [Match Patterns](https://developer.chrome.com/extensions/develop/concepts/match-patterns) for more details on the syntax of these strings.
- id
  
  string
  
  The id of the content script, specified in the API call. Must not start with a '\_' as it's reserved as a prefix for generated script IDs.
- js
  
  string\[] optional
  
  The list of JavaScript files to be injected into matching pages. These are injected in the order they appear in this array.
- matchOriginAsFallback
  
  boolean optional
  
  Chrome 119+
  
  Indicates whether the script can be injected into frames where the URL contains an unsupported scheme; specifically: about:, data:, blob:, or filesystem:. In these cases, the URL's origin is checked to determine if the script should be injected. If the origin is `null` (as is the case for data: URLs) then the used origin is either the frame that created the current frame or the frame that initiated the navigation to this frame. Note that this may not be the parent frame.
- matches
  
  string\[] optional
  
  Specifies which pages this content script will be injected into. See [Match Patterns](https://developer.chrome.com/extensions/develop/concepts/match-patterns) for more details on the syntax of these strings. Must be specified for [`registerContentScripts`](#method-registerContentScripts).
- persistAcrossSessions
  
  boolean optional
  
  Specifies if this content script will persist into future sessions. The default is true.
- runAt
  
  [RunAt](https://developer.chrome.com/docs/extensions/reference/extensionTypes/#type-RunAt) optional
  
  Specifies when JavaScript files are injected into the web page. The preferred and default value is `document_idle`.
- world
  
  [ExecutionWorld](#type-ExecutionWorld) optional
  
  Chrome 102+
  
  The JavaScript "world" to run the script in. Defaults to `ISOLATED`.

### ScriptInjection

#### Properties

- args
  
  any\[] optional
  
  Chrome 92+
  
  The arguments to pass to the provided function. This is only valid if the `func` parameter is specified. These arguments must be JSON-serializable.
- files
  
  string\[] optional
  
  The path of the JS or CSS files to inject, relative to the extension's root directory. Exactly one of `files` or `func` must be specified.
- injectImmediately
  
  boolean optional
  
  Chrome 102+
  
  Whether the injection should be triggered in the target as soon as possible. Note that this is not a guarantee that injection will occur prior to page load, as the page may have already loaded by the time the script reaches the target.
- target
  
  [InjectionTarget](#type-InjectionTarget)
  
  Details specifying the target into which to inject the script.
- world
  
  [ExecutionWorld](#type-ExecutionWorld) optional
  
  Chrome 95+
  
  The JavaScript "world" to run the script in. Defaults to `ISOLATED`.
- func
  
  void optional
  
  Chrome 92+
  
  A JavaScript function to inject. This function will be serialized, and then deserialized for injection. This means that any bound parameters and execution context will be lost. Exactly one of `files` or `func` must be specified.
  
  The `func` function looks like:
  
  ```
  () => {...}
  ```

### StyleOrigin

The origin for a style change. See [style origins](https://developer.mozilla.org/en-US/docs/Glossary/Style_origin) for more info.

#### Enum

"AUTHOR"

"USER"

## Methods

### executeScript()

Promise

```
chrome.scripting.executeScript(
  injection: ScriptInjection,
  callback?: function,
)
```

Injects a script into a target context. By default, the script will be run at `document_idle`, or immediately if the page has already loaded. If the `injectImmediately` property is set, the script will inject without waiting, even if the page has not finished loading. If the script evaluates to a promise, the browser will wait for the promise to settle and return the resulting value.

#### Parameters

- injection
  
  [ScriptInjection](#type-ScriptInjection)
  
  The details of the script which to inject.
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (results: InjectionResult[]) => void
  ```
  
  - results
    
    [InjectionResult](#type-InjectionResult)\[]

#### Returns

- Promise&lt;[InjectionResult](#type-InjectionResult)\[]&gt;
  
  Chrome 90+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### getRegisteredContentScripts()

Promise Chrome 96+

```
chrome.scripting.getRegisteredContentScripts(
  filter?: ContentScriptFilter,
  callback?: function,
)
```

Returns all dynamically registered content scripts for this extension that match the given filter.

#### Parameters

- filter
  
  [ContentScriptFilter](#type-ContentScriptFilter) optional
  
  An object to filter the extension's dynamically registered scripts.
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (scripts: RegisteredContentScript[]) => void
  ```
  
  - scripts
    
    [RegisteredContentScript](#type-RegisteredContentScript)\[]

#### Returns

- Promise&lt;[RegisteredContentScript](#type-RegisteredContentScript)\[]&gt;
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### insertCSS()

Promise

```
chrome.scripting.insertCSS(
  injection: CSSInjection,
  callback?: function,
)
```

Inserts a CSS stylesheet into a target context. If multiple frames are specified, unsuccessful injections are ignored.

#### Parameters

- injection
  
  [CSSInjection](#type-CSSInjection)
  
  The details of the styles to insert.
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

#### Returns

- Promise&lt;void&gt;
  
  Chrome 90+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### registerContentScripts()

Promise Chrome 96+

```
chrome.scripting.registerContentScripts(
  scripts: RegisteredContentScript[],
  callback?: function,
)
```

Registers one or more content scripts for this extension.

#### Parameters

- scripts
  
  [RegisteredContentScript](#type-RegisteredContentScript)\[]
  
  Contains a list of scripts to be registered. If there are errors during script parsing/file validation, or if the IDs specified already exist, then no scripts are registered.
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

#### Returns

- Promise&lt;void&gt;
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### removeCSS()

Promise Chrome 90+

```
chrome.scripting.removeCSS(
  injection: CSSInjection,
  callback?: function,
)
```

Removes a CSS stylesheet that was previously inserted by this extension from a target context.

#### Parameters

- injection
  
  [CSSInjection](#type-CSSInjection)
  
  The details of the styles to remove. Note that the `css`, `files`, and `origin` properties must exactly match the stylesheet inserted through [`insertCSS`](#method-insertCSS). Attempting to remove a non-existent stylesheet is a no-op.
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

#### Returns

- Promise&lt;void&gt;
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### unregisterContentScripts()

Promise Chrome 96+

```
chrome.scripting.unregisterContentScripts(
  filter?: ContentScriptFilter,
  callback?: function,
)
```

Unregisters content scripts for this extension.

#### Parameters

- filter
  
  [ContentScriptFilter](#type-ContentScriptFilter) optional
  
  If specified, only unregisters dynamic content scripts which match the filter. Otherwise, all of the extension's dynamic content scripts are unregistered.
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

#### Returns

- Promise&lt;void&gt;
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### updateContentScripts()

Promise Chrome 96+

```
chrome.scripting.updateContentScripts(
  scripts: RegisteredContentScript[],
  callback?: function,
)
```

Updates one or more content scripts for this extension.

#### Parameters

- scripts
  
  [RegisteredContentScript](#type-RegisteredContentScript)\[]
  
  Contains a list of scripts to be updated. A property is only updated for the existing script if it is specified in this object. If there are errors during script parsing/file validation, or if the IDs specified do not correspond to a fully registered script, then no scripts are updated.
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  () => void
  ```

#### Returns

- Promise&lt;void&gt;
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.