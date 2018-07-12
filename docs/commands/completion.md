# Gotree: toolkit and api for phylogenetic tree manipulation

## Commands

### completion
This command generates auto-completion commands for bash or zsh. 

#### Usage

```
Usage:
  gotree completion SHELL
```

#### Bash
* Install bash-completion:
```
# MacOS
brew install bash-completion
# Linux
yum install bash-completion -y
apt-get install bash-completion
```

* Activate gotree bash completion
```
# Once
source <(gotree completion bash)
# Permanently
mkdir ~/.gotree
gotree completion bash > ~/.gotree/completion.bash.inc
printf "
# gotree shell completion
source '$HOME/.gotree/completion.bash.inc'
" >> $HOME/.bashrc
```

#### Zsh (not tested)

```
# Once
source <(kubectl completion zsh)
# Permanently
gotree completion zsh > "${fpath[1]}/_gotree"
```
