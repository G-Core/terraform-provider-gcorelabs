# Terraform volume create provider

### Quick start
Install go and terraform. Then
``
go build -o terraform-provider-gcore_v0.1
``
In Windows move the file to *%APPDATA%\terraform.d\plugins\windows_amd64*.
Then you can do:
```
terraform init
terraform plan
terraform apply
terraform destroy
```

### Description
This provider allows you to create a volume and to change its size or type.

### Import 
Existing volumes can be loaded from the cloud. Firstly, create a new volume record for a loading volume in a ``.tf`` file:
```
resource "gcore_volume" "<loading_volume_name>" {
}
```
then run in a teminal:
```
terraform import gcore_volume.<loading_volume_name> <project_id>:<region_id>:<loading_volume_uuid>
```

Example:
 in main.tf add:
    ```
    resource "gcore_volume" "foo" {
    }
    ```
 in a command line:
    ```
    terraform import gcore_volume.foo 2:1:7057f675-ed04-4001-9025-b58e34cd7327
    ```


### Tests
To run tests, set the environment variables:
```
TF_ACC (=1/true)
OS_PROVIDER_JWT
TEST_PROJECT_ID or TEST_PROJECT_NAME
TEST_REGION_ID or TEST_REGION_NAME
```
then run them:
```
go test -v
```