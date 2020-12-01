Terraform G-Core Labs Provider
------------------------------
- Website: https://gcorelabs.com 
- Technical partner: https://itkey.com

<img src="https://gcorelabs.com/img/logo.svg" data-src="https://gcorelabs.com/img/logo.svg" alt="G-Core Labs" width="500px" width="500px"> 
=========================================================================================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.13.x
-	[Go](https://golang.org/doc/install) 1.14 (to build the provider plugin)

Building the provider
---------------------
```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers
$ cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/G-Core/terraform-provider-gcorelabs.git
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-gcorelabs
$ make build
```

Using the provider
------------------
To use the provider, prepare configuration files based on examples

```sh
$ cp ./examples/... .
$ terraform init
```

Thank You
