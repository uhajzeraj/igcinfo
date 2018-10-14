igctracer
=======

### IMT2681 Assignment 1 
### Student: Uran Hajzeraj, 501303 ######

[![Build Status](https://travis-ci.com/uhajzeraj/igcinfo.svg?branch=master)](https://travis-ci.com/uhajzeraj/igcinfo)

### GET /api

* Response type: application/json
* Response code: 200

```
{
  "uptime": <uptime>
  "info": "Service for IGC tracks."
  "version": "v1"
}
```

### POST /api/igc

* Response type: application/json
* Response code: 200 if everything is OK, appropriate error code otherwise

```
{
  "url": "<url>"
}
```

### GET /api/igc

* Response type: application/json
* Response code: 200 if everything is OK, appropriate error code otherwise

```
[<id1>, <id2>, ...]
```

### GET /api/igc/`<id>`

* Response type: application/json
* Response code: 200 if everything is OK, appropriate error code otherwise

```
"H_date": <date from File Header, H-record>,
"pilot": <pilot>,
"glider": <glider>,
"glider_id": <glider_id>,
"track_length": <calculated total track length>
}
```

### GET /api/igc/`<id>`/`<field>`

* Response type: text/plain
* Response code: 200 if everything is OK, appropriate error code otherwise
   * `<pilot>` for `pilot`
   * `<glider>` for `glider`
   * `<glider_id>` for `glider_id`
   * `<calculated total track length>` for `track_length`
   * `<H_date>` for `H_date`