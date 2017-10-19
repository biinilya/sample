# Toptal sample project


<a name="overview"></a>
## Overview
beego has a very cool tools to autogenerate documents for your API


### Version information
*Version* : 1.0.0


### Contact information
*Contact Email* : astaxie@gmail.com


### License information
*License* : Apache 2.0  
*License URL* : http://www.apache.org/licenses/LICENSE-2.0.html  
*Terms of service* : http://beego.me/


### URI scheme
*BasePath* : /api/v1


### Tags

* permissions :  PermissionsController operations for Permissions





<a name="paths"></a>
## Paths

<a name="permissionscontroller-post"></a>
### POST /permissions/

#### Description
create Permissions


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Body**|**body**  <br>*required*|body for Permissions content|[models.Permissions](#models-permissions)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**201**|{int} models.Permissions|No Content|
|**403**|body is empty|No Content|


#### Tags

* permissions


<a name="permissionscontroller-get-all"></a>
### GET /permissions/

#### Description
get Permissions


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Query**|**fields**  <br>*optional*|Fields returned. e.g. col1,col2 ...|string|
|**Query**|**limit**  <br>*optional*|Limit the size of result set. Must be an integer|string|
|**Query**|**offset**  <br>*optional*|Start position of result set. Must be an integer|string|
|**Query**|**order**  <br>*optional*|Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ...|string|
|**Query**|**query**  <br>*optional*|Filter. e.g. col1:v1,col2:v2 ...|string|
|**Query**|**sortby**  <br>*optional*|Sorted-by fields. e.g. col1,col2 ...|string|


#### Responses

|HTTP Code|Schema|
|---|---|
|**200**|[models.Permissions](#models-permissions)|
|**403**|No Content|


#### Tags

* permissions


<a name="permissionscontroller-get-one"></a>
### GET /permissions/{id}

#### Description
get Permissions by id


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Path**|**id**  <br>*required*|The key for staticblock|string|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.Permissions](#models-permissions)|
|**403**|:id is empty|No Content|


#### Tags

* permissions


<a name="permissionscontroller-put"></a>
### PUT /permissions/{id}

#### Description
update the Permissions


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Path**|**id**  <br>*required*|The id you want to update|string|
|**Body**|**body**  <br>*required*|body for Permissions content|[models.Permissions](#models-permissions)|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**||[models.Permissions](#models-permissions)|
|**403**|:id is not int|No Content|


#### Tags

* permissions


<a name="permissionscontroller-delete"></a>
### DELETE /permissions/{id}

#### Description
delete the Permissions


#### Parameters

|Type|Name|Description|Schema|
|---|---|---|---|
|**Path**|**id**  <br>*required*|The id you want to delete|string|


#### Responses

|HTTP Code|Description|Schema|
|---|---|---|
|**200**|{string} delete success!|No Content|
|**403**|id is empty|No Content|


#### Tags

* permissions




<a name="definitions"></a>
## Definitions

<a name="models-permissions"></a>
### models.Permissions

|Name|Schema|
|---|---|
|**Description**  <br>*optional*|string|
|**Id**  <br>*optional*|integer (int64)|
|**Title**  <br>*optional*|string|





