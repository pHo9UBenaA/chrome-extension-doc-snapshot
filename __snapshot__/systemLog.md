# chrome.systemLog

**Important:** This API works **only on ChromeOS**.

## Description

Use the `chrome.systemLog` API to record Chrome system logs from extensions.

## Permissions

`systemLog`

## Availability

Chrome 125+ ChromeOS only [Requires policy](https://support.google.com/chrome/a/answer/9296680)

## Types

### MessageOptions

#### Properties

- message
  
  string

## Methods

### add()

```
chrome.systemLog.add(
  options: MessageOptions,
): Promise<void>
```

Adds a new log record.

#### Parameters

- options
  
  [MessageOptions](#type-MessageOptions)
  
  The logging options.

#### Returns

- Promise&lt;void&gt;