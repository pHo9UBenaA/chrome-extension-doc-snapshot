# chrome.dns

## Description

Use the `chrome.dns` API for dns resolution.

## Permissions

`dns`

## Availability

Dev channel

To use this API, you must declare the `"dns"` permission in the [manifest](/docs/extensions/mv3/manifest).

```
{
  "name": "My extension",
  ...
  "permissions": [
    "dns"
  ],
  ...
}
```

**Note:** This API is only available in [Chrome Dev](https://www.google.com/chrome/dev/). There are no foreseeable plans to move this API from the dev channel into Chrome stable.

## Usage

The following code calls [`resolve()`](#method-resolve) to retrieve the IP address of `example.com`.

service-worker.js:

```
const resolveDNS = async () => {
    let record = await chrome.dns.resolve('example.com');
    console.log(record.address); // "192.0.2.172"
};

resolveDNS();
```

**Key point:** Do not include the scheme or trailing slash in the hostname. For example, `https://example.com/` is invalid.

## Types

### ResolveCallbackResolveInfo

#### Properties

- address
  
  string optional
  
  A string representing the IP address literal. Supplied only if resultCode indicates success.
- resultCode
  
  number
  
  The result code. Zero indicates success.

## Methods

### resolve()

```
chrome.dns.resolve(
  hostname: string,
): Promise<ResolveCallbackResolveInfo>
```

Resolves the given hostname or IP address literal.

#### Parameters

- hostname
  
  string
  
  The hostname to resolve.

#### Returns

- Promise&lt;[ResolveCallbackResolveInfo](#type-ResolveCallbackResolveInfo)&gt;