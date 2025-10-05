# chrome.wallpaper

**Important:** This API works **only on ChromeOS**.

## Description

Use the `chrome.wallpaper` API to change the ChromeOS wallpaper.

## Permissions

`wallpaper`

You must declare the "wallpaper" permission in the app's [manifest](/docs/extensions/reference/manifest) to use the wallpaper API. For example:

```
{
  "name": "My extension",
  ...
  "permissions": [
    "wallpaper"
  ],
  ...
}
```

## Availability

Chrome 43+ ChromeOS only

## Examples

For example, to set the wallpaper as the image at `https://example.com/a_file.png`, you can call `chrome.wallpaper.setWallpaper` this way:

```
chrome.wallpaper.setWallpaper(
  {
    'url': 'https://example.com/a_file.jpg',
    'layout': 'CENTER_CROPPED',
    'filename': 'test_wallpaper'
  },
  function() {}
);
```

## Types

### WallpaperLayout

Chrome 44+

The supported wallpaper layouts.

#### Enum

"STRETCH"

"CENTER"

"CENTER\_CROPPED"

## Methods

### setWallpaper()

```
chrome.wallpaper.setWallpaper(
  details: object,
): Promise<ArrayBuffer | undefined>
```

Sets wallpaper to the image at *url* or *wallpaperData* with the specified *layout*

#### Parameters

- details
  
  object
  
  - data
    
    ArrayBuffer optional
    
    The jpeg or png encoded wallpaper image as an ArrayBuffer.
  - filename
    
    string
    
    The file name of the saved wallpaper.
  - layout
    
    [WallpaperLayout](#type-WallpaperLayout)
    
    The supported wallpaper layouts.
  - thumbnail
    
    boolean optional
    
    True if a 128x60 thumbnail should be generated. Layout and ratio are not supported yet.
  - url
    
    string optional
    
    The URL of the wallpaper to be set (can be relative).

#### Returns

- Promise&lt;ArrayBuffer | undefined&gt;
  
  Chrome 96+