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
|**Header**|**X-Request-Id**  <br>*required*|Access Token|string|


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
|**Body**|**body**  <br>*required*|Credentials to Sign In|[models.UserCredentialsView](#models-usercredentialsview)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.UserAccessTokenView](#models-useraccesstokenview)|
|**400**|bad request|No Content|
|**401**|wrong credentials|No Content|


#### Tags

* user


<a name="usercontroller-get-one"></a>
### GET /user/{uid}

#### Description
get User


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Header**|**X-Request-Id**  <br>*required*|Access Token|string|
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
|**Header**|**X-Request-Id**  <br>*required*|Access Token|string|
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
|**Header**|**X-Request-Id**  <br>*required*|Access Token|string|
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




<a name="definitions"></a>
## Definitions

<a name="models-useraccesstokenview"></a>
### models.UserAccessTokenView

|Name|Schema|
|---|---|
|**access-token**  <br>*optional*|string|


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





