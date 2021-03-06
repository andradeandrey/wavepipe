API
===

wavepipe features a simple API which is used to retrieve metadata from media files, as well as endpoints
to retrieve a file stream from the server.

An information endpoint can be found at the root of the API, `/api`.  This endpoint contains API metadata
such as the current API version, supported API versions, and a link to this documentation.

At this time, the current API version is **v0**.  This API is **unstable**, and is subject to change.

API calls may respond to a variety of different HTTP methods.  This documentation will outline the function
of each method, but they are generally used as follows:
  - **GET**: retrieve one or more read-only resources from the API
  - **POST**: create a resource on the API
  - **PUT**: partially or fully update an existing resource on the API
  - **PATCH**: equivalent to **PUT**, partially or fully update an existing resource on the API
  - **DELETE**: delete a resource from the API

For additional security, wavepipe employs a very simple roles system.  In general, these roles are used as follows:
  - **Guest**: read-only access to the entire API, and no ability to update their own credentials
  - **User**: full API access, Last.fm scrobbling, the ability to update **only** their own credentials
  - **Administrator**: full API access, Last.fm scrobbling, full access to create/update/delete all users

If a user attempts to perform an action which is disallowed by their current role, they will receive a
`HTTP 403 Forbidden` error.

**Authentication:**

In order to use the wavepipe API, all requests must be authenticated.  The first step is to generate a new
session via the [Login](#login) API.  Login username and password can be passed using either HTTP Basic or
via POST body.  In addition, an optional client parameter may be passed, which will identify this session
with the given name. Both methods can be demonstrated with `curl` as follows:

```
$ curl -X POST -u test:test http://localhost:8080/api/v0/login
$ curl -X POST -d "username=test&password=test&client=testclient" http://localhost:8080/api/v0/login
```

Example [Session](http://godoc.org/github.com/mdlayher/wavepipe/data#Session) output is as follows:

```json
{
	"error": null,
	"session": {
		"id": 1,
		"userId": 1,
		"client": "testclient",
		"expire": 1397713157,
		"key": "abcdef0123456789abcdef0123456789"
	}
}
```

Upon successful login, a session key is generated, which is used to authenticate subsequent requests.
It should be noted that unless the service is secured via HTTPS, this token can be compromised by other
users on the same network.  For this reason, it is recommended to place wavepipe behind SSL.

This method can be demonstrated with `curl` as follows.

```
$ curl http://localhost:8080/api/v0/albums?s=abcdef0123456789abcdef0123456789
```

Sessions which are not used for one week will expire.  Each subsequent API request with a specified session will
update the expiration time to one week in the future.

**Table of Contents:**

| Name | Versions | Description |
| :--: | :------: | :---------: |
| [Albums](#albums) | v0 | Used to retrieve information about albums from wavepipe. |
| [Art](#art) | v0 | Used to retrieve a binary data stream of an art file from wavepipe. |
| [Artists](#artists) | v0 | Used to retrieve information about artists from wavepipe. |
| [Folders](#folders) | v0 | Used to retrieve information about folders from wavepipe. |
| [LastFM](#lastfm) | v0 | Used to scrobble songs from wavepipe to Last.fm. |
| [Login](#login) | v0 | Used to generate a new API session on wavepipe. |
| [Logout](#logout) | v0 | Used to destroy the current API session from wavepipe. |
| [Search](#search) | v0 | Used to retrieve artists, albums, songs, and folders which match a specified search query. |
| [Songs](#songs) | v0 | Used to retrieve information about songs from wavepipe. |
| [Status](#status) | v0 | Used to retrieve current server status from wavepipe, as well as server metrics, if specified. |
| [Stream](#stream) | v0 | Used to retrieve a raw, non-transcoded, binary data stream of a media file from wavepipe. |
| [Transcode](#transcode) | v0 | Used to retrieve transcoded binary data stream of a media file from wavepipe. |
| [Users](#users) | v0 | Used to retrieve information about users from wavepipe. |
| [Waveform](#waveform) | v0 | Used to generate and return a waveform image of a media file from wavepipe. |

## Albums
Used to retrieve information about albums from wavepipe.  If an ID is specified, information will be
retrieved about a single album.

**Versions:** `v0`

**URL:** `GET /api/v0/albums/:id`

**Examples:**
  - `GET http://localhost:8080/api/v0/albums/`
  - `GET http://localhost:8080/api/v0/albums/1`
  - `GET http://localhost:8080/api/v0/albums?limit=0,100`

**Query Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| limit | v0 | integer,integer | | Comma-separated integer pair which limits the number of returned results.  First integer is the offset, second integer is the item count. |

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error)/null | Information about any errors that occurred.  Value is null if no error occurred. |
| albums | \[\][Album](http://godoc.org/github.com/mdlayher/wavepipe/data#Album) | Array of Album objects returned by the API. |
| songs | \[\][Song](http://godoc.org/github.com/mdlayher/wavepipe/data#Song)/null | If ID is specified, array of Song objects attached to this album.  Value is null if no ID specified. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 400 | invalid comma-separated integer pair for limit | A valid integer pair could not be parsed from the limit parameter. Input must be in the form "x,y". |
| 400 | invalid integer album ID | A valid integer could not be parsed from the ID. |
| 404 | album ID not found | An album with the specified ID does not exist. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## Art
Used to retrieve a binary data stream of an art file from wavepipe.  An ID **must** be specified to access an art stream.
Successful calls with return a binary stream, and unsuccessful ones will return a JSON error.

**Versions:** `v0`

**URL:** `GET /api/v0/art/:id`

**Examples:**
  - `GET http://localhost:8080/api/v0/art/1`
  - `GET http://localhost:8080/api/v0/art/1?size=500`

**Query Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| size | v0 | integer | | Scale the art to the specified width in pixels. The art's original aspect ratio will be preserved. |

**Return Binary:** Binary data stream containing the art file stream.

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error) | Information about any errors that occurred. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 400 | no integer art ID provided | No integer ID was sent in request. |
| 400 | invalid art stream ID | A valid integer could not be parsed from the ID. |
| 400 | invalid integer size | A valid integer could not be parsed from the size parameter. |
| 400 | negative integer size | A negative integer was passed to the size parameter. Size **must** be a positive integer. |
| 404 | art ID not found | An art file with the specified ID does not exist. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## Artists
Used to retrieve information about artists from wavepipe.  If an ID is specified, information will be
retrieved about a single artist.

**Versions:** `v0`

**URL:** `GET /api/v0/artists/:id`

**Examples:**
  - `GET http://localhost:8080/api/v0/artists/`
  - `GET http://localhost:8080/api/v0/artists/1`
  - `GET http://localhost:8080/api/v0/artists?limit=0,100`
  - `GET http://localhost:8080/api/v0/artists/1?songs=true`

**Query Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| limit | v0 | string "integer,integer" | | Comma-separated integer pair which limits the number of returned results.  First integer is the offset, second integer is the item count. |
| songs | v0 | boolean | | If true, returns all songs attached to this artist. |

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error)/null | Information about any errors that occurred.  Value is null if no error occurred. |
| artists | \[\][Artist](http://godoc.org/github.com/mdlayher/wavepipe/data#Artist) | Array of Artist objects returned by the API. |
| albums | \[\][Album](http://godoc.org/github.com/mdlayher/wavepipe/data#Album)/null | If ID is specified, array of Album objects attached to this artist. |
| songs | \[\][Song](http://godoc.org/github.com/mdlayher/wavepipe/data#Song)/null | If parameter `songs` is true, array of Song objects attached to this artist.  Value is null if parameter `songs` is false or not specified. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 400 | invalid comma-separated integer pair for limit | A valid integer pair could not be parsed from the limit parameter. Input must be in the form "x,y". |
| 400 | invalid integer artist ID | A valid integer could not be parsed from the ID. |
| 404 | artist ID not found | An artist with the specified ID does not exist. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## Folders
Used to retrieve information about folders from wavepipe.  If an ID is specified, information will be
retrieved about a single folder.

**Versions:** `v0`

**URL:** `GET /api/v0/folders/:id`

**Examples:**
  - `GET http://localhost:8080/api/v0/folders/`
  - `GET http://localhost:8080/api/v0/folders/1`
  - `GET http://localhost:8080/api/v0/folders?limit=0,100`

**Query Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| limit | v0 | integer,integer | | Comma-separated integer pair which limits the number of returned results.  First integer is the offset, second integer is the item count. |

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error)/null | Information about any errors that occurred.  Value is null if no error occurred. |
| folders | \[\][Folder](http://godoc.org/github.com/mdlayher/wavepipe/data#Folder) | Array of Folder objects returned by the API. |
| subfolders | \[\][Folder](http://godoc.org/github.com/mdlayher/wavepipe/data#Folder) | If ID is specified, array of Folder objects which are children to the current folder, returned by the API. Value is null if no ID is specified. |
| songs | \[\][Song](http://godoc.org/github.com/mdlayher/wavepipe/data#Song)/null | If ID is specified, array of Song objects attached to this folder.  Value is null if no ID specified. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 400 | invalid comma-separated integer pair for limit | A valid integer pair could not be parsed from the limit parameter. Input must be in the form "x,y". |
| 400 | invalid integer folder ID | A valid integer could not be parsed from the ID. |
| 404 | folder ID not found | An folder with the specified ID does not exist. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## LastFM
Used to scrobble songs from wavepipe to Last.fm.  The user must first complete a `login` action with their Last.fm
credentials, and then the `nowplaying` and `scrobble` actions may be used.  After the initial `login`, wavepipe
will store an API key for the user, and use this key for future requests.

Ideally, a `nowplaying` action will be triggered by clients as soon as the track begins playing on that client.
After a fair amount of time has passed (for example, 50-75% of the song), a `scrobble` request should be triggered
to commit the play to Last.fm.

Last.fm actions are only allowed for users with the role `User` or `Administrator`.  `Guest` users are not permitted
to use Last.fm functionality, and will receive a `HTTP 403 Forbidden` error when accessing this API call.

**Versions:** `v0`

**URL:** `POST /api/v0/lastfm/:action/:id`

**Examples:**
  - `POST http://localhost:8080/api/v0/lastfm/login "username=test&password=test"`
  - `POST http://localhost:8080/api/v0/lastfm/nowplaying/1`
  - `POST http://localhost:8080/api/v0/lastfm/scrobble/1`
  - `POST http://localhost:8080/api/v0/lastfm/scrobble/1?timestamp=1403384162`

**POST Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| username | v0 | string | | Username used to authenticate to Last.fm via wavepipe. Only used for the `login` action. |
| password | v0 | string | | Password used to authenticate to Last.fm via wavepipe. Only used for the `login` action. |

**Query Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| timestamp | v0 | integer | | Optional integer UNIX timestamp, which can be used to specify a past timestamp. The current timestamp is used if not specified. |

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error)/null | Information about any errors that occurred.  Value is null if no error occurred. |
| url | string | String containing the URL required to authorize wavepipe's Last.fm token for this user. Only returned on the `login` action, or if other actions are accessed while the token is not authorized. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 400 | no string action provided | No action was specified in the URL.  An action **must** be specified to use Last.fm functionality. |
| 400 | invalid string action provided | An unknown action was specified in the URL.  Valid actions are `login`, `nowplaying`, and `scrobble`. |
| 400 | login: no username provided | No Last.fm username was passed via POST body. Only returned on `login` action. |
| 400 | login: no password provided | No Last.fm password was passed via POST body. Only returned on `login` action. |
| 400 | no integer song ID provided | No integer ID was sent in request. Only returned on `nowplaying` and `scrobble` actions. |
| 400 | invalid integer song ID | A valid integer could not be parsed from the ID. Only returned on `nowplaying` and `scrobble` actions. |
| 401 | action: last.fm authentication failed | Could not authenticate to Last.fm. Could be due to invalid username/password, or an invalid API token. |
| 401 | action: user must authenticate to last.fm | User attempted to perform `nowplaying` or `scrobble` action, without first completing `login` action. |
| 401 | action: last.fm token not yet authorized | User must authorize wavepipe to access their Last.fm account, via the provided URL. |
| 403 | permission denied | The current user is forbidden from performing this action. |
| 404 | song ID not found | A song with the specified ID does not exist. Only returned on `nowplaying` and `scrobble` actions. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## Login
Used to generate a new API session on wavepipe.  Credentials may be provided either via query string,
or using a HTTP Basic username and password combination.

**Versions:** `v0`

**URL:** `POST /api/v0/login`

**Examples:**
  - `POST http://localhost:8080/api/v0/login "username=test&password=test&client=testclient"`

**POST Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| username | v0 | string | X | Username used to authenticate to wavepipe. Can also be passed via HTTP Basic. |
| password | v0 | string | X | Associated password used to authenticate to wavepipe. Can also be passed via HTTP Basic. |
| client | v0 | string | | Optional client name used to identify this session. |

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error)/null | Information about any errors that occurred.  Value is null if no error occurred. |
| session | [Session](http://godoc.org/github.com/mdlayher/wavepipe/data#Session) | Session object which contains the public and secret keys used to authenticate further API calls. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 401 | authentication failed: X | API authentication failed. Could be due to malformed, missing, or bad credentials. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## Logout
Used to destroy the current API session from wavepipe.

**Versions:** `v0`

**URL:** `POST /api/v0/logout`

**Examples:**
  - `POST http://localhost:8080/api/v0/logout`

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error)/null | Information about any errors that occurred.  Value is null if no error occurred. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## Search
Used to retrieve artists, albums, songs, and folders which match a specified search query.  A search query **must** be
specified to retrieve results.

**Versions:** `v0`

**URL:** `GET /api/v0/search/:query`

**Examples:**
  - `GET http://localhost:8080/api/v0/search/boston`
  - `GET http://localhost:8080/api/v0/search/boston?type=artists,songs`

**Query Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| type | v0 | string | | Comma-separated string containing object types (`artists`, `albums`, `songs`, `folders`) to return search results. If not specified, equivalent to `artists,albums,songs,folders`. |

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error) | Information about any errors that occurred. |
| artists | \[\][Artist](http://godoc.org/github.com/mdlayher/wavepipe/data#Artist) | Array of Artist objects with titles matching the search query. |
| albums | \[\][Album](http://godoc.org/github.com/mdlayher/wavepipe/data#Album) | Array of Album objects with titles matching the search query. |
| songs | \[\][Song](http://godoc.org/github.com/mdlayher/wavepipe/data#Song) | Array of Song objects with titles matching the search query. |
| folders | \[\][Folder](http://godoc.org/github.com/mdlayher/wavepipe/data#Folder) | Array of Folder objects with titles matching the search query. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 400 | no search query specified | No search query was specified in the URL. A search query **must** be specified to retrieve results. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## Songs
Used to retrieve information about songs from wavepipe.  If an ID is specified, information will be
retrieved about a single song.

**Versions:** `v0`

**URL:** `GET /api/v0/songs/:id`

**Examples:**
  - `GET http://localhost:8080/api/v0/songs/`
  - `GET http://localhost:8080/api/v0/songs/1`
  - `GET http://localhost:8080/api/v0/songs?limit=0,100`
  - `GET http://localhost:8080/api/v0/songs?random=10`

**Query Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| limit | v0 | integer,integer | | Comma-separated integer pair which limits the number of returned results.  First integer is the offset, second integer is the item count. |
| random | v0 | integer | | If specified, wavepipe will return N random songs instead of the entire list, where N is the integer specified in this parameter. |

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error)/null | Information about any errors that occurred.  Value is null if no error occurred. |
| songs | \[\][Song](http://godoc.org/github.com/mdlayher/wavepipe/data#Song) | Array of Song objects returned by the API. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 400 | invalid comma-separated integer pair for limit | A valid integer pair could not be parsed from the limit parameter. Input must be in the form "x,y". |
| 400 | invalid integer for random | A valid integer could not be parsed from the random parameter. |
| 400 | invalid integer song ID | A valid integer could not be parsed from the ID. |
| 404 | song ID not found | A song with the specified ID does not exist. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## Status
Used to retrieve current server status from wavepipe, as well as server metrics, if specified.

**Versions:** `v0`

**URL:** `GET /api/v0/status`

**Examples:**
  - `GET http://localhost:8080/api/v0/status`
  - `GET http://localhost:8080/api/v0/status?metrics=all`
  - `GET http://localhost:8080/api/v0/status?metrics=database`

**Query Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| metrics | v0 | string | | Comma-separated string containing metric types (`all`, `database`, `network`) to return. If not specified, no metrics will be returned. |

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error)/null | Information about any errors that occurred.  Value is null if no error occurred. |
| status | [Status](http://godoc.org/github.com/mdlayher/wavepipe/common#Status) | Status object containing current server information, returned by the API. |
| metrics | [Metrics](http://godoc.org/github.com/mdlayher/wavepipe/common#Metrics)/null | Metrics object containing current server metrics, returned by the API. Value is null unless parameter `metrics` contains a comma-separated list of metric types. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## Stream
Used to retrieve a raw, non-transcoded, binary data stream of a media file from wavepipe.  An ID **must** be specified to access a file stream.  Successful calls with return a binary stream, and unsuccessful ones will return a JSON error.

**Versions:** `v0`

**URL:** `GET /api/v0/stream/:id`

**Examples:**
  - `GET http://localhost:8080/api/v0/stream/1`

**Return Binary:** Binary data stream containing the media file stream.

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error) | Information about any errors that occurred. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 400 | no integer stream ID provided | No integer ID was sent in request. |
| 400 | invalid integer stream ID | A valid integer could not be parsed from the ID. |
| 404 | song ID not found | A song with the specified ID does not exist. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## Transcode
Used to retrieve a transcoded binary data stream of a media file from wavepipe.  An ID **must** be specified to access a file stream.  Successful calls with return a binary stream, and unsuccessful ones will return a JSON error.

**Versions:** `v0`

**URL:** `GET /api/v0/transcode/:id`

**Examples:**
  - `GET http://localhost:8080/api/v0/transcode/1`
  - `GET http://localhost:8080/api/v0/transcode/1?codec=MP3&quality=320`

**Query Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| codec | v0 | string | | The codec selected for use by the transcoder.  If not specified, defaults to **MP3**.  Options are: **MP3**, OGG, OPUS (lowercase variants will be automatically capitalized). |
| quality | v0 | string/integer | | The quality selected for use by the transcoder.  String options specify VBR encodings, while integer options specify CBR encodings.  If not specified, defaults to **192**. |

**Available Codecs:**

| Codec | Versions | Type | Options | Description |
| :---: | :------: | :--: | :-----: | :---------: |
| MP3 | v0 | CBR | 128, **192** (default), 256, 320 | Generates a constant bitrate encode using LAME. |
| MP3 | v0 | VBR | V0 (~245kbps), V2 (~190kbps), V4 (~165kbps) | Generates a variable bitrate encode using a specific LAME quality level. |
| OGG | v0 | CBR | 128, **192** (default), 256, 320, 500 | Generates a constant bitrate encode using Ogg Vorbis. |
| OGG | v0 | VBR | Q10 (~500kbps), Q8 (~256kbps), Q6 (~192kbps) | Generates a variable bitrate encode using a specific Ogg Vorbis quality level. |
| OPUS | v0 | CBR | 128, **192** (default), 256, 320, 500 | Generates a constant bitrate encode using Ogg Opus. |
| OPUS | v0 | VBR | Q10 (~500kbps), Q8 (~256kbps), Q6 (~192kbps) | Generates a variable bitrate encode using a specific Ogg Opus quality level. |

**Return Binary:** Binary data stream containing the transcoded media file stream.

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error) | Information about any errors that occurred. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 400 | no integer transcode ID provided | No integer ID was sent in request. |
| 400 | invalid integer transcode ID | A valid integer could not be parsed from the ID. |
| 400 | invalid transcoder codec: X | A non-existant transcoder codec was passed via the codec parameter. |
| 400 | invalid quality for codec X: X | A non-existant quality setting for the specified codec was passed via the quality parameter. |
| 404 | song ID not found | A song with the specified ID does not exist. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |
| 503 | ffmpeg not found, transcoding disabled | ffmpeg binary could not be detected in system PATH, so the transcoding subsystem is disabled. |
| 503 | ffmpeg codec libmp3lame not found, MP3 transcoding disabled | ffmpeg was not compiled with libmp3lame codec, so MP3 transcoding is disabled. |
| 503 | ffmpeg codec libvorbis not found, OGG transcoding disabled | ffmpeg was not compiled with libvorbis codec, so Ogg Vorbis transcoding is disabled. |
| 503 | ffmpeg codec libopus not found, OPUS transcoding disabled | ffmpeg was not compiled with libopus codec, so Ogg Opus transcoding is disabled. |

## Users
Used to retrieve information about users from wavepipe.  If an ID is specified, information will be
retrieved about a single user.

In addition, this API call may be used to create, modify, or delete existing users.  Different functionality is
available for each user role:
  - Users with the role `Administrator` have full control over all users.
  - Users with the role `User` may update their own account, but may not create, delete, or update other users.
  - Users with the role `Guest` have no access to create, delete, or update any user.

If a user is disallowed from performing an action, they will receive a `HTTP 403 Forbidden` error.

**Versions:** `v0`

**URL:** `GET/POST/PUT/PATCH/DELETE /api/v0/users/:id`

**Examples:**
  - `GET http://localhost:8080/api/v0/users/`
  - `GET http://localhost:8080/api/v0/users/1`
  - `POST http://localhost:8080/api/v0/users "username=test&password=test&role=2"`
  - `PUT http://localhost:8080/api/v0/users/1 "username=test2&password=test2"`
  - `PATCH http://localhost:8080/api/v0/users/1 "username=test3"`
  - `DELETE http://localhost:8080/api/v0/users/1`

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error)/null | Information about any errors that occurred.  Value is null if no error occurred. |
| users | \[\][User](http://godoc.org/github.com/mdlayher/wavepipe/data#User) | Array of User objects returned by the API. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 400 | invalid integer user ID | A valid integer could not be parsed from the ID. |
| 400 | invalid integer role ID | A valid integer could not be parsed from the role ID, or an invalid role was specified. |
| 400 | missing required parameter: username | No username specified in POST body during user creation. |
| 400 | missing required parameter: password | No password specified in POST body during user creation. |
| 400 | missing required parameter: role | No role specified in POST body during user creation. |
| 403 | permission denied | The current user is forbidden from performing this action. |
| 403 | cannot delete current user | User attempted to delete itself, which is forbidden. |
| 404 | user ID not found | A user with the specified ID does not exist. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |

## Waveform
Used to generate and return a waveform image of a media file from wavepipe.  An ID **must** be specified to access a
file stream.  Successful calls with return a binary stream, and unsuccessful ones will return a JSON error.

**Versions:** `v0`

**URL:** `GET /api/v0/waveform/:id`

**Examples:**
  - `GET http://localhost:8080/api/v0/waveform/1`
  - `GET http://localhost:8080/api/v0/waveform/1?fg=%23FF0000&bg=%230000FF&alt=%2300FF00`
  - `GET http://localhost:8080/api/v0/waveform/1?size=1024x256`
  - `GET http://localhost:8080/api/v0/waveform/1?size=1024x0`

**Query Parameters:**

| Name | Versions | Type | Required | Description |
| :--: | :------: | :--: | :------: | :---------: |
| bg | v0 | string | | The hex background color for the waveform image. If not specified, defaults to **#FFFFFF** (white). Invalid hex strings will be ignored, and the default will be used. |
| fg | v0 | string | | The hex foreground color for the waveform image. If not specified, defaults to **#000000** (black). Invalid hex strings will be ignored, and the default will be used. |
| alt | v0 | string | | The hex alternate color for the waveform image. Creates a striping effect with the foreground color. If not specified, defaults to the foreground color. Invalid hex strings will be ignored, and the default will be used. |
| size | v0 | integerxinteger | | Scale the waveform to the specified width and height in pixels. If height is 0, the waveform's original aspect ratio will be preserved. |

**Return Binary:** Binary data stream containing a waveform image generated from a media file stream.

**Return JSON:**

| Name | Type | Description |
| :--: | :--: | :---------: |
| error | [Error](http://godoc.org/github.com/mdlayher/wavepipe/api#Error) | Information about any errors that occurred. |

**Possible errors:**

| Code | Message | Description |
| :--: | :-----: | :---------: |
| 400 | unsupported API version: vX | Attempted access to an invalid version of this API, or to a version before this API existed. |
| 400 | no integer song ID provided | No integer ID was sent in request. |
| 400 | invalid integer song ID | A valid integer could not be parsed from the ID. |
| 400 | invalid x-separated integer pair for size | A valid integer pair could not be parsed from the size parameter. Input must be in the form "XxY". |
| 404 | song ID not found | A song with the specified ID does not exist. |
| 500 | server error | An internal error occurred. wavepipe will log these errors to its console log. |
| 501 | unsupported audio format | The song is in an unsupported format, which cannot be decoded to a waveform. |
