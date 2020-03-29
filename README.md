# Terraform volume create provider

Description
-----------
This provider allows you to create a volume and to change its size or type.

Building the provider
---------------------
Install terraform and Go. Clone the provider, then build it
``
go build -o terraform-provider-gcore_v0.1
``

Using the provider
------------------
###### Provider
Provider can be created in ``.tf`` file or set the variables in an environment 

| Name in tf file    | Alternative name in environment       | type    | Required |
| :----------------: |:-------------------------:| :------:| :-------:|
| username           | GCORE_PROVIDER_USERNAME   | string  | true     | 
| password           | GCORE_PROVIDER_PASSWORD   | string  | true     | 
| platform_url       | GCORE_PLATFORM_URL        | string  | false    | 
| host               | GCORE_HOST                | string  | false    | 
| timeout            | GCORE_TIMEOUT             | integer | false    | 

example:
```
provider "gcore" {
  username = "..."
  password = "..."
}
```

##### Volume 

Volume should have the fields:

| Name | Type  |Required  | Alternative name      | Alternative type    | Note |
| :----------------: |:-------------------------:| :------:| :-------:| :------:| :-------:| 
|size       | int |   true| | | |
|source     | string |true | | | one of 'new-volume', 'image', 'snapshot' |
|name       | string | true| | | |
|project_id| int | true| prject_name| string | |
|region_id| int| true| region_name| string | |
|type_name |  string | false | | | one of 'standard', 'ssd_hiiops', 'cold' |
|image_id  |  string | false | | | it must be used when source is 'image' |
|snapshot_id | string | false | | | is must be used when source is 'snapshot' |
|count   | int | false | | | set it if you want create more objects than one |


Terraform commands
------------------
### Import 
Existing volumes can be loaded from the cloud. Firstly, create a new volume record for a loading volume in a ``.tf`` file:
```
resource "gcore_volumeV1" "<loading_volume_name>" {
}
```

then run in a teminal:
```
terraform import gcore_volumeV1.<loading_volume_name> <project_id>:<region_id>:<loading_volume_uuid>
```

   ###### Example:
   in main.tf add:
      ```
      resource "gcore_volumeV1" "foo" {
      }
      ```
   
   then in a command line:
      ```
      terraform import gcore_volumeV1.foo 2:1:7057f675-ed04-4001-9025-b58e34cd7327
      ```
   where ``project id = 2``, ``regiont id = 1``. Project id and region id will be saved in state of this volume.

### Update 
You can update volume size or volume type.

Test the provider
-----------------
To run tests, set the environment variables:
```
TF_ACC
GCORE_PROVIDER_USERNAME
GCORE_PROVIDER_PASSWORD
TEST_PROJECT_ID or TEST_PROJECT_NAME
TEST_REGION_ID or TEST_REGION_NAME
```
then run them:
```
go test -v
```