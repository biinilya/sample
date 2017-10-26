# Toptal demo app (Run Keeper)


<a name="overview"></a>
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





<a name="paths"></a>
## Paths

<a name="usercontroller-post"></a>
### POST /user/

#### Description
create new User, returns key and secret used to authorize


#### Responses

|HTTP Code|Schema|
|---|---|
|**201**|[models.UserCredentialsView](#models-usercredentialsview)|


#### Tags

* user


<a name="usercontroller-get-all"></a>
### GET /user/

#### Description
Users Directory


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Query**|**filter**  <br>*optional*|Filter users e.x. (key eq 'xxx')|string|
|**Query**|**limit**  <br>*optional*|Limit number of records to (default 50)|integer (int64)|
|**Query**|**offset**  <br>*optional*|Offset in records|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.UserInfoView](#models-userinfoview) > array|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user


<a name="usercontroller-signin"></a>
### POST /user/sign_in

#### Description
Use this method to receive Access Token which is required for API access


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Body**|**body**  <br>*required*|Credentials to Sign In|[models.UserCredentialsData](#models-usercredentialsdata)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.UserAccessTokenView](#models-useraccesstokenview)|
|**400**|bad request|No Content|
|**403**|forbidden|No Content|


#### Tags

* user


<a name="usercontroller-get-one"></a>
### GET /user/{uid}

#### Description
get User


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Path**|**uid**  <br>*required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.UserInfoView](#models-userinfoview)|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user


<a name="usercontroller-delete"></a>
### DELETE /user/{uid}

#### Description
delete the User


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Path**|**uid**  <br>*required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.UserInfoView](#models-userinfoview)|
|**400**|bad request (uid is missing or is not a number)|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user


<a name="usercontroller-put"></a>
### PUT /user/{uid}/credentials

#### Description
generate new User credentials


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Path**|**uid**  <br>*required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.UserCredentialsView](#models-usercredentialsview)|
|**400**|bad request (uid is missing or is not a number)|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user


<a name="usercontroller-add-permission"></a>
### GET /user/{uid}/permission

#### Description
delete the User


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Path**|**uid**  <br>*required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.PermissionView](#models-permissionview) > array|
|**400**|bad request (uid is missing or is not a number)|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user


<a name="usercontroller-add-permission"></a>
### POST /user/{uid}/permission/{title}

#### Description
delete the User


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Path**|**title**  <br>*required*|Title of Permission to Add to user|string|
|**Path**|**uid**  <br>*required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.PermissionView](#models-permissionview) > array|
|**400**|bad request (uid is missing or is not a number)|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user


<a name="usercontroller-del-permission"></a>
### DELETE /user/{uid}/permission/{title}

#### Description
delete the User


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Path**|**title**  <br>*required*|Title of Permission to Add to user|string|
|**Path**|**uid**  <br>*required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.PermissionView](#models-permissionview) > array|
|**400**|bad request (uid is missing or is not a number)|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user


<a name="recordcontroller-post"></a>
### POST /user/{uid}/record/

#### Description
create new Record, returns new record data


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Path**|**uid**  <br>*required*|User ID|integer (int64)|
|**Body**|**body**  <br>*required*|Record|[models.RecordData](#models-recorddata)|


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


<a name="recordcontroller-getall"></a>
### GET /user/{uid}/record/

#### Description
Get all user-related records


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Path**|**uid**  <br>*required*|User ID|integer (int64)|
|**Query**|**filter**  <br>*optional*|Filter records e.x. (date eq '2017-01-01')|string|
|**Query**|**limit**  <br>*optional*|Limit number of records to (default 50)|integer (int64)|
|**Query**|**offset**  <br>*optional*|Offset in records|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.RecordView](#models-recordview) > array|
|**400**|Bad request|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user/:uid/record


<a name="recordcontroller-weeklyreport"></a>
### GET /user/{uid}/record/report/weekly

#### Description
Get average distance and duration per week


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Path**|**uid**  <br>*required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||< [models.WeeklyReport](#models-weeklyreport) > array|
|**400**|Bad request|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user/:uid/record


<a name="recordcontroller-put"></a>
### PUT /user/{uid}/record/{record_id}

#### Description
update existing record, returns updated record data


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Path**|**record_id**  <br>*required*|Record ID|integer (int64)|
|**Path**|**uid**  <br>*required*|User ID|integer (int64)|
|**Body**|**body**  <br>*required*|Record|[models.RecordData](#models-recorddata)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.RecordView](#models-recordview)|
|**400**|Bad request|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user/:uid/record


<a name="recordcontroller-delete"></a>
### DELETE /user/{uid}/record/{record_id}

#### Description
update existing record, returns updated record data


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Access-Token**  <br>*required*|Access Token|string|
|**Path**|**record_id**  <br>*required*|Record ID|integer (int64)|
|**Path**|**uid**  <br>*required*|User ID|integer (int64)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.RecordView](#models-recordview)|
|**400**|Bad request|No Content|
|**401**|unauthorized|No Content|
|**403**|forbidden|No Content|


#### Tags

* user/:uid/record




<a name="definitions"></a>
## Definitions

<a name="models-permissionview"></a>
### models.PermissionView

|Name|Schema|
|---|---|
|**Description**  <br>*optional*|string|
|**Title**  <br>*optional*|string|


<a name="models-recorddata"></a>
### models.RecordData

|Name|Schema|
|---|---|
|**date**  <br>*optional*|string (string)|
|**distance**  <br>*optional*|number (double)|
|**duration**  <br>*optional*|number (double)|
|**latitude**  <br>*optional*|number (double)|
|**longitude**  <br>*optional*|number (double)|


<a name="models-recordview"></a>
### models.RecordView

|Name|Schema|
|---|---|
|**date**  <br>*optional*|string (string)|
|**distance**  <br>*optional*|number (double)|
|**duration**  <br>*optional*|number (double)|
|**id**  <br>*optional*|integer (int64)|
|**latitude**  <br>*optional*|number (double)|
|**longitude**  <br>*optional*|number (double)|
|**weather**  <br>*optional*|[models.WeatherInfo](#models-weatherinfo)|


<a name="models-useraccesstokenview"></a>
### models.UserAccessTokenView

|Name|Schema|
|---|---|
|**access-token**  <br>*optional*|string|


<a name="models-usercredentialsdata"></a>
### models.UserCredentialsData

|Name|Schema|
|---|---|
|**key**  <br>*optional*|string|
|**secret**  <br>*optional*|string|


<a name="models-usercredentialsview"></a>
### models.UserCredentialsView

|Name|Schema|
|---|---|
|**id**  <br>*optional*|integer (int64)|
|**key**  <br>*optional*|string|
|**secret**  <br>*optional*|string|


<a name="models-userinfoview"></a>
### models.UserInfoView

|Name|Schema|
|---|---|
|**created**  <br>*optional*|string (string)|
|**id**  <br>*optional*|integer (int64)|
|**key**  <br>*optional*|string|
|**updated**  <br>*optional*|string (string)|


<a name="models-weatherinfo"></a>
### models.WeatherInfo
*Type* : object


<a name="models-weeklyreport"></a>
### models.WeeklyReport

|Name|Schema|
|---|---|
|**avg_distance**  <br>*optional*|number (double)|
|**avg_duration**  <br>*optional*|number (double)|
|**week**  <br>*optional*|integer (int64)|
|**year**  <br>*optional*|integer (int64)|





