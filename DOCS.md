# Documentation
## Installation:
### To install the latest version via homebrew (MacOS only):
##### 1. Install [Homebrew](https://brew.sh/)
##### 2. Run the following commands:
```shell
brew tap hazelops/ize
```

```shell
brew install ize
```

Now you can run `ize` from command shell by typing `ize` in console.

#### 4. Update:
4.1 Uninstall previous version:

```shell
brew uninstall ize
```

4.2 Update version in brew repo: `

```shell
brew tap hazelops/ize
```

4.3 Install:

```shell
brew install ize
```

### Installation via public apt repository URL (Ubuntu):
##### 1. To add public apt repository run:
 ```shell
echo "deb [trusted=yes] https://apt.fury.io/hazelops/ /" | sudo tee /etc/apt/sources.list.d/fury.list
```

##### 2. After this, you should update information. Run:
```shell
sudo apt-get update
```

##### 3. To install the latest version of `ize` app, you should run:
```shell
sudo apt-get install ize 
```

##### 4. If you wish to install certain version of the `ize` you should add version like this:
 ```shell
sudo apt-get install ize=<version>
 ```

##### 6. To remove `ize` app - run this command:
```shell
sudo apt-get purge ize
```

### Installation via Scoop (Windows only)
#### Prerequisites:
You must install Scoop (command-line installer for Windows). To install Scoop see the following [manual](https://github.com/ScoopInstaller/Scoop#installation)

To install `ize` with Scoop:

1. Add bucket with the ize source:

```shell
scoop bucket add hazelops/scoop-ize
```
2. Install ize:

```shell
scoop install ize
```




### Installation from source:
#### Prerequisites:
- GO version should be 1.16+
- `GOPATH` environment variable is set to `~/go`

To install Ize from source download code or clone it from this repo. After this you should run:

```shell
go mod download
make install
```


### Autocomplete:
You could use integrated option to add autocompletion to Ize commands (bash, fish, zsh, powershell). In this manual we will describe it only for zsh and bash.

To add autocompletion script, use the following manual:

##### 1. ZSH:
If shell completion is not already enabled in your environment you will need to enable it. You should execute the following once:

```shell
echo "autoload -U compinit; compinit" >>  ~/.zshrc
```

To load completions for every new session, execute once:

###### 1.1 macOS:
```shell
ize gen completion zsh > /usr/local/share/zsh/site-functions/_ize
```

###### 1.2 Linux:
You will need root privileges.

```shell
sudo zsh
```
Input your root password and run:

```shell
ize gen completion zsh > "${fpath[1]}/_ize"
```

To take effect for this setup you should run `source ~/.zshrc` or simply restart shell.


##### 2. Bash:
Autocompletion script depends on the `bash-completion` package.

If it is not installed already, you can install it via your OS’s package manager.

To load completions for every new session, you should execute once:

###### 2.1 MacOS:
```shell
ize gen completion bash > /usr/local/etc/bash_completion.d/ize
```

###### 2.2 Linux:
You will need root privileges.

```shell
sudo bash
```
Input your root password and run:

```shell
ize gen completion bash > /etc/bash_completion.d/ize
```

To take effect for this setup you should run `source ~/.bashrc` or simply restart shell.
