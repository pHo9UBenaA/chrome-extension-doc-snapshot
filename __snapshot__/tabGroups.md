# chrome.tabGroups

## Description

Use the `chrome.tabGroups` API to interact with the browser's tab grouping system. You can use this API to modify and rearrange tab groups in the browser. To group and ungroup tabs, or to query what tabs are in groups, use the `chrome.tabs` API.

## Permissions

`tabGroups`

## Availability

Chrome 89+ MV3+

## Types

### Color

The group's color.

#### Enum

"grey"

"blue"

"red"

"yellow"

"green"

"pink"

"purple"

"cyan"

"orange"

### TabGroup

#### Properties

- collapsed
  
  boolean
  
  Whether the group is collapsed. A collapsed group is one whose tabs are hidden.
- color
  
  [Color](#type-Color)
  
  The group's color.
- id
  
  number
  
  The ID of the group. Group IDs are unique within a browser session.
- title
  
  string optional
  
  The title of the group.
- windowId
  
  number
  
  The ID of the window that contains the group.

## Properties

### TAB\_GROUP\_ID\_NONE

An ID that represents the absence of a group.

#### Value

-1

## Methods

### get()

Promise

```
chrome.tabGroups.get(
  groupId: number,
  callback?: function,
)
```

Retrieves details about the specified group.

#### Parameters

- groupId
  
  number
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (group: TabGroup) => void
  ```
  
  - group
    
    [TabGroup](#type-TabGroup)

#### Returns

- Promise&lt;[TabGroup](#type-TabGroup)&gt;
  
  Chrome 90+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### move()

Promise

```
chrome.tabGroups.move(
  groupId: number,
  moveProperties: object,
  callback?: function,
)
```

Moves the group and all its tabs within its window, or to a new window.

#### Parameters

- groupId
  
  number
  
  The ID of the group to move.
- moveProperties
  
  object
  
  - index
    
    number
    
    The position to move the group to. Use `-1` to place the group at the end of the window.
  - windowId
    
    number optional
    
    The window to move the group to. Defaults to the window the group is currently in. Note that groups can only be moved to and from windows with [`windows.WindowType`](https://developer.chrome.com/docs/extensions/reference/windows/#type-WindowType) type `"normal"`.
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (group?: TabGroup) => void
  ```
  
  - group
    
    [TabGroup](#type-TabGroup) optional
    
    Details about the moved group.

#### Returns

- Promise&lt;[TabGroup](#type-TabGroup) | undefined&gt;
  
  Chrome 90+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### query()

Promise

```
chrome.tabGroups.query(
  queryInfo: object,
  callback?: function,
)
```

Gets all groups that have the specified properties, or all groups if no properties are specified.

#### Parameters

- queryInfo
  
  object
  
  - collapsed
    
    boolean optional
    
    Whether the groups are collapsed.
  - color
    
    [Color](#type-Color) optional
    
    The color of the groups.
  - title
    
    string optional
    
    Match group titles against a pattern.
  - windowId
    
    number optional
    
    The ID of the parent window, or [`windows.WINDOW_ID_CURRENT`](https://developer.chrome.com/docs/extensions/reference/windows/#property-WINDOW_ID_CURRENT) for the [current window](https://developer.chrome.com/docs/extensions/reference/windows/#current-window).
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (result: TabGroup[]) => void
  ```
  
  - result
    
    [TabGroup](#type-TabGroup)\[]

#### Returns

- Promise&lt;[TabGroup](#type-TabGroup)\[]&gt;
  
  Chrome 90+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

### update()

Promise

```
chrome.tabGroups.update(
  groupId: number,
  updateProperties: object,
  callback?: function,
)
```

Modifies the properties of a group. Properties that are not specified in `updateProperties` are not modified.

#### Parameters

- groupId
  
  number
  
  The ID of the group to modify.
- updateProperties
  
  object
  
  - collapsed
    
    boolean optional
    
    Whether the group should be collapsed.
  - color
    
    [Color](#type-Color) optional
    
    The color of the group.
  - title
    
    string optional
    
    The title of the group.
- callback
  
  function optional
  
  The `callback` parameter looks like:
  
  ```
  (group?: TabGroup) => void
  ```
  
  - group
    
    [TabGroup](#type-TabGroup) optional
    
    Details about the updated group.

#### Returns

- Promise&lt;[TabGroup](#type-TabGroup) | undefined&gt;
  
  Chrome 90+
  
  Promises are supported in Manifest V3 and later, but callbacks are provided for backward compatibility. You cannot use both on the same function call. The promise resolves with the same type that is passed to the callback.

## Events

### onCreated

```
chrome.tabGroups.onCreated.addListener(
  callback: function,
)
```

Fired when a group is created.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (group: TabGroup) => void
  ```
  
  - group
    
    [TabGroup](#type-TabGroup)

### onMoved

```
chrome.tabGroups.onMoved.addListener(
  callback: function,
)
```

Fired when a group is moved within a window. Move events are still fired for the individual tabs within the group, as well as for the group itself. This event is not fired when a group is moved between windows; instead, it will be removed from one window and created in another.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (group: TabGroup) => void
  ```
  
  - group
    
    [TabGroup](#type-TabGroup)

### onRemoved

```
chrome.tabGroups.onRemoved.addListener(
  callback: function,
)
```

Fired when a group is closed, either directly by the user or automatically because it contained zero tabs.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (group: TabGroup) => void
  ```
  
  - group
    
    [TabGroup](#type-TabGroup)

### onUpdated

```
chrome.tabGroups.onUpdated.addListener(
  callback: function,
)
```

Fired when a group is updated.

#### Parameters

- callback
  
  function
  
  The `callback` parameter looks like:
  
  ```
  (group: TabGroup) => void
  ```
  
  - group
    
    [TabGroup](#type-TabGroup)