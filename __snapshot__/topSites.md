# chrome.topSites

## Description

Use the `chrome.topSites` API to access the top sites (i.e. most visited sites) that are displayed on the new tab page. These do not include shortcuts customized by the user.

## Permissions

`topSites`

You must declare the "topSites" permission in your [extension's manifest](/docs/extensions/reference/manifest) to use this API.

```
{
  "name": "My extension",
  ...
  "permissions": [
    "topSites",
  ],
  ...
}
```

## Examples

To try this API, install the [topSites API example](https://github.com/GoogleChrome/chrome-extensions-samples/tree/main/api-samples/topSites) from the [chrome-extension-samples](https://github.com/GoogleChrome/chrome-extensions-samples/tree/main/api-samples) repository.

## Types

### MostVisitedURL

An object encapsulating a most visited URL, such as the default shortcuts on the new tab page.

#### Properties

- title
  
  string
  
  The title of the page
- url
  
  string
  
  The most visited URL.

## Methods

### get()

Promise

```
chrome.topSites.get(
  callback?: function,
)
```

Gets a list of top sites.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (data: MostVisitedURL[]) => void
  ```
  
  - data
    
    [MostVisitedURL](#type-MostVisitedURL)\[]

#### Returns

- Promise&lt;[MostVisitedURL](#type-MostVisitedURL)\[]&gt;
  
  Chrome 96+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.