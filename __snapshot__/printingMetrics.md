# chrome.printingMetrics

**Important:** This API works **only on ChromeOS**.

## Description

Use the `chrome.printingMetrics` API to fetch data about printing usage.

## Permissions

`printingMetrics`

## Availability

Chrome 79+ ChromeOS only [Requires policy](https://support.google.com/chrome/a/answer/9296680)

## Types

### ColorMode

#### Enum

"BLACK\_AND\_WHITE"  
Specifies that black and white mode was used.

"COLOR"  
Specifies that color mode was used.

### DuplexMode

#### Enum

"ONE\_SIDED"  
Specifies that one-sided printing was used.

"TWO\_SIDED\_LONG\_EDGE"  
Specifies that two-sided printing was used, flipping on long edge.

"TWO\_SIDED\_SHORT\_EDGE"  
Specifies that two-sided printing was used, flipping on short edge.

### MediaSize

#### Properties

- height
  
  number
  
  Height (in micrometers) of the media used for printing.
- vendorId
  
  string
  
  Vendor-provided ID, e.g. "iso\_a3\_297x420mm" or "na\_index-3x5\_3x5in". Possible values are values of "media" IPP attribute and can be found on [IANA page](https://www.iana.org/assignments/ipp-registrations/ipp-registrations.xhtml) .
- width
  
  number
  
  Width (in micrometers) of the media used for printing.

### Printer

#### Properties

- name
  
  string
  
  Displayed name of the printer.
- source
  
  [PrinterSource](#type-PrinterSource)
  
  The source of the printer.
- uri
  
  string
  
  The full path for the printer. Contains protocol, hostname, port, and queue.

### PrinterSource

The source of the printer.

#### Enum

"USER"  
Specifies that the printer was added by user.

"POLICY"  
Specifies that the printer was added via policy.

### PrintJobInfo

#### Properties

- completionTime
  
  number
  
  The job completion time (in milliseconds past the Unix epoch).
- creationTime
  
  number
  
  The job creation time (in milliseconds past the Unix epoch).
- id
  
  string
  
  The ID of the job.
- numberOfPages
  
  number
  
  The number of pages in the document.
- printer
  
  [Printer](#type-Printer)
  
  The info about the printer which printed the document.
- printer\_status
  
  [PrinterStatus](https://developer.chrome.com/docs/extensions/reference/printing/#type-PrinterStatus)
  
  Chrome 85+
  
  The status of the printer.
- settings
  
  [PrintSettings](#type-PrintSettings)
  
  The settings of the print job.
- source
  
  [PrintJobSource](#type-PrintJobSource)
  
  Source showing who initiated the print job.
- sourceId
  
  string optional
  
  ID of source. Null if source is PRINT\_PREVIEW or ANDROID\_APP.
- status
  
  [PrintJobStatus](#type-PrintJobStatus)
  
  The final status of the job.
- title
  
  string
  
  The title of the document which was printed.

### PrintJobSource

The source of the print job.

#### Enum

"PRINT\_PREVIEW"  
Specifies that the job was created from the Print Preview page initiated by the user.

"ANDROID\_APP"  
Specifies that the job was created from an Android App.

"EXTENSION"  
Specifies that the job was created by extension via Chrome API.

"ISOLATED\_WEB\_APP"  
Specifies that the job was created by an Isolated Web App via API.

### PrintJobStatus

Specifies the final status of the print job.

#### Enum

"FAILED"  
Specifies that the print job was interrupted due to some error.

"CANCELED"  
Specifies that the print job was canceled by the user or via API.

"PRINTED"  
Specifies that the print job was printed without any errors.

### PrintSettings

#### Properties

- color
  
  [ColorMode](#type-ColorMode)
  
  The requested color mode.
- copies
  
  number
  
  The requested number of copies.
- duplex
  
  [DuplexMode](#type-DuplexMode)
  
  The requested duplex mode.
- mediaSize
  
  [MediaSize](#type-MediaSize)
  
  The requested media size.

## Methods

### getPrintJobs()

Promise

```
chrome.printingMetrics.getPrintJobs(
  callback?: function,
)
```

Returns the list of the finished print jobs.

#### Parameters

- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (jobs: PrintJobInfo[]) => void
  ```
  
  - jobs
    
    [PrintJobInfo](#type-PrintJobInfo)\[]

#### Returns

- Promise&lt;[PrintJobInfo](#type-PrintJobInfo)\[]&gt;
  
  Chrome 96+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

## Events

### onPrintJobFinished

```
chrome.printingMetrics.onPrintJobFinished.addListener(
  callback: function,
)
```

Event fired when the print job is finished. This includes any of termination statuses: FAILED, CANCELED and PRINTED.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (jobInfo: PrintJobInfo) => void
  ```
  
  - jobInfo
    
    [PrintJobInfo](#type-PrintJobInfo)