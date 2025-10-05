# chrome.pageCapture

## Description

Use the `chrome.pageCapture` API to save a tab as MHTML.

MHTML is a [standard format](https://tools.ietf.org/html/rfc2557) supported by most browsers. It encapsulates in a single file a page and all its resources (CSS files, images..).

Note that for security reasons a MHTML file can only be loaded from the file system and that it can only be loaded in the main frame.

## Permissions

`pageCapture`

You must declare the "pageCapture" permission in the [extension manifest](/docs/extensions/reference/manifest) to use the pageCapture API. For example:

```
{
  "name": "My extension",
  ...
  "permissions": [
    "pageCapture"
  ],
  ...
}
```

## Methods

### saveAsMHTML()

```
chrome.pageCapture.saveAsMHTML(
  details: object,
): Promise<Blob | undefined>
```

Saves the content of the tab with given id as MHTML.

#### Parameters

- details
  
  object
  
  - tabId
    
    number
    
    The id of the tab to save as MHTML.

#### Returns

- Promise&lt;Blob | undefined&gt;
  
  Chrome 116+