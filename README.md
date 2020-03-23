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
Provider can be created in ``.tf`` file 
```
provider "gcore" {
  username = "..."
  password = "..."
}
```
or be set an enviroment variable ``OS_PROJECT_USERNAME`` and ``OS_PROJECT_PASSWORD``. Also you can set ``HOST`` and ``OS_TIMEOUT`` to change default values.

Volume should have the fields:
```
size        int    required
source      string required, one of 'new-volume', 'image', 'snapshot'
name        string required
type_name   string optional, one of 'standard', 'ssd_hiiops', 'cold'
image_id    string optional
snapshot_id string optional
count       int    optional
```
also in volume body you should set 
1. ``project_id``(int) or ``project_name``(string)
2. ``region_id``(int) or ``region_name``(string)

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
OS_PROVIDER_USERNAME
OS_PROVIDER_PASSWORD
TEST_PROJECT_ID or TEST_PROJECT_NAME
TEST_REGION_ID or TEST_REGION_NAME
```
then run them:
```
go test -v
```