### Cloud Service For File Storage
Hey, welcome to my repository.
This is a basic cloud service for file storage written in golang (1.21.00) and gin. As you are observing, this is not the latest golang version, this is due to the deployment issues encountered during deploying using koyeb.

`Deployed At: https://cloudstorage-sambhavmahajan.koyeb.app`

#### Here's how it works!
- Create an account using register.
- Upload any file you like.
- You may logout or delete the file.
- You may visit again and login.

#### How to clone.
- Ensure golang version 1.21.00(atleast) is installed.
- Clone the repository: `git clone https://github.com/sambhavmahajan/Cloud-Service-For-File-Storage`
- `go mod tidy` to install any dependencies, if gin is not installed it will automatically install it.
- Set your own port and binding, modify router.Run(), default port is `8080` and default binding is `localhost`
  - Change Port: `router.Run(":9090")`
  - Change Bindng and Port: `router.Run("0.0.0.0:9090")`
- Visit `localhost:8080` or the appropriate binding and port.

If you find this project cool enough, please consider starring it on GitHub :)\
Licensed under the Apache 2.0 License. See [LICENSE](LICENSE) for details.
