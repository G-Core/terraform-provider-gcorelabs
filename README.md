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

### Tests
To run tests, set the environment variables:
```
TF_ACC (=1/true)
TEST_PROVIDER_JWT
TEST_PROJECT_ID or TEST_PROJECT_NAME
TEST_REGION_ID or TEST_REGION_NAME
```
then run them:
```
go test -v
```