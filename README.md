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