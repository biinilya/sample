# sample demo app (Run Keeper)



## Overview
This RunKeeper API allows you to manage users run data


### Version information
*Version* : 1.0.0


### Contact information
*Contact Email* : me@ilyabiin.com


### URI scheme
*BasePath* : /api/v1


### Tags

* user :  UserController operations for User

* user/:uid/record :  RecordController operations for Record






## Paths


### POST /user/

#### Description
create new User, returns key and secret used to authorize


#### Responses

|HTTP Code|Schema|
|---|---|
|**201**|[models.UserCredentialsView](#models-usercredentialsview)|


#### Tags

* user



### GET /user/

#### Description
Users Directory


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Query**|**filter**  *optional*|Filter users e.x. (key eq 'xxx')|string|
|**Query**|**limit**  *optional*|Limit number of records to (default 50)|integer (int64)|
|**Query**|**offset**  *optional*|Offset in records|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.UserInfoView](#models-userinfoview) > array|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user



### POST /user/sign_in

#### Description
Use this method to receive Access Token which is required for API access


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Body**|**body**  *required*|Credentials to Sign In|[models.UserCredentialsData](#models-usercredentialsdata)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.UserAccessTokenView](#models-useraccesstokenview)|
|**400**|bad request|No Content|
|**403**|forbidden|No Content|


#### Tags

* user



### GET /user/{uid}

#### Description
get User


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Path**|**uid**  *required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.UserInfoView](#models-userinfoview)|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user



### DELETE /user/{uid}

#### Description
delete the User


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Path**|**uid**  *required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.UserInfoView](#models-userinfoview)|
|**400**|bad request (uid is missing or is not a number)|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user



### PUT /user/{uid}/credentials

#### Description
generate new User credentials


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Path**|**uid**  *required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.UserCredentialsView](#models-usercredentialsview)|
|**400**|bad request (uid is missing or is not a number)|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user



### GET /user/{uid}/permission

#### Description
delete the User


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Path**|**uid**  *required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.PermissionView](#models-permissionview) > array|
|**400**|bad request (uid is missing or is not a number)|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user



### POST /user/{uid}/permission/{title}

#### Description
delete the User


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Path**|**title**  *required*|Title of Permission to Add to user|string|
|**Path**|**uid**  *required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.PermissionView](#models-permissionview) > array|
|**400**|bad request (uid is missing or is not a number)|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user



### DELETE /user/{uid}/permission/{title}

#### Description
delete the User


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Path**|**title**  *required*|Title of Permission to Add to user|string|
|**Path**|**uid**  *required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.PermissionView](#models-permissionview) > array|
|**400**|bad request (uid is missing or is not a number)|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user



### POST /user/{uid}/record/

#### Description
create new Record, returns new record data


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Path**|**uid**  *required*|User ID|integer (int64)|
|**Body**|**body**  *required*|Record|[models.RecordData](#models-recorddata)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**201**||[models.RecordView](#models-recordview)|
|**400**|Bad request|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|
|**503**|weather service is down|No Content|


#### Tags

* user/:uid/record



### GET /user/{uid}/record/

#### Description
Get all user-related records


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Path**|**uid**  *required*|User ID|integer (int64)|
|**Query**|**filter**  *optional*|Filter records e.x. (date eq '2017-01-01')|string|
|**Query**|**limit**  *optional*|Limit number of records to (default 50)|integer (int64)|
|**Query**|**offset**  *optional*|Offset in records|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.RecordView](#models-recordview) > array|
|**400**|Bad request|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user/:uid/record



### GET /user/{uid}/record/report/weekly

#### Description
Get average distance and duration per week


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Path**|**uid**  *required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.WeeklyReport](#models-weeklyreport) > array|
|**400**|Bad request|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user/:uid/record



### PUT /user/{uid}/record/{record_id}

#### Description
update existing record, returns updated record data


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Path**|**record_id**  *required*|Record ID|integer (int64)|
|**Path**|**uid**  *required*|User ID|integer (int64)|
|**Body**|**body**  *required*|Record|[models.RecordData](#models-recorddata)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.RecordView](#models-recordview)|
|**400**|Bad request|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user/:uid/record



### DELETE /user/{uid}/record/{record_id}

#### Description
update existing record, returns updated record data


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  *required*|Access Token|string|
|**Path**|**record_id**  *required*|Record ID|integer (int64)|
|**Path**|**uid**  *required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.RecordView](#models-recordview)|
|**400**|Bad request|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user/:uid/record





## Definitions


### models.PermissionView

|Name|Schema|
|---|---|
|**Description**  *optional*|string|
|**Title**  *optional*|string|



### models.RecordData

|Name|Schema|
|---|---|
|**date**  *optional*|string (string)|
|**distance**  *optional*|number (double)|
|**duration**  *optional*|number (double)|
|**latitude**  *optional*|number (double)|
|**longitude**  *optional*|number (double)|



### models.RecordView

|Name|Schema|
|---|---|
|**date**  *optional*|string (string)|
|**distance**  *optional*|number (double)|
|**duration**  *optional*|number (double)|
|**id**  *optional*|integer (int64)|
|**latitude**  *optional*|number (double)|
|**longitude**  *optional*|number (double)|
|**weather**  *optional*|[models.WeatherInfo](#models-weatherinfo)|



### models.UserAccessTokenView

|Name|Schema|
|---|---|
|**access-token**  *optional*|string|



### models.UserCredentialsData

|Name|Schema|
|---|---|
|**key**  *optional*|string|
|**secret**  *optional*|string|



### models.UserCredentialsView

|Name|Schema|
|---|---|
|**id**  *optional*|integer (int64)|
|**key**  *optional*|string|
|**secret**  *optional*|string|



### models.UserInfoView

|Name|Schema|
|---|---|
|**created**  *optional*|string (string)|
|**id**  *optional*|integer (int64)|
|**key**  *optional*|string|
|**updated**  *optional*|string (string)|



### models.WeatherInfo
*Type* : object



### models.WeeklyReport

|Name|Schema|
|---|---|
|**avg_distance**  *optional*|number (double)|
|**avg_duration**  *optional*|number (double)|
|**week**  *optional*|integer (int64)|
|**year**  *optional*|integer (int64)|





