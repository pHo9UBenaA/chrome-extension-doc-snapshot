# chrome.enterprise.login

This API is only for [extensions installed by a policy](https://cloud.google.com/blog/products/chrome-enterprise/how-to-manage-chrome-extensions-on-windows-and-mac).

**Important:** This API works **only on ChromeOS**.

## Description

Use the `chrome.enterprise.login` API to exit Managed Guest sessions. Note: This API is only available to extensions installed by enterprise policy in ChromeOS Managed Guest sessions.

## Permissions

`enterprise.login`

## Availability

Chrome 139+ ChromeOS only [Requires policy](https://support.google.com/chrome/a/answer/9296680)

## Methods

### exitCurrentManagedGuestSession()

```
chrome.enterprise.login.exitCurrentManagedGuestSession(): Promise<void>
```

Exits the current managed guest session.

#### Returns

- Promise&lt;void&gt;