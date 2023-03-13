### [Tools collections](https://kubernetes.io/docs/tasks/tools/)
https://kubernetes.io/docs/tasks/tools/

### bash auto-completion on Linux

##### Install bash-completion
```shell
apt-get install bash-completion
yum install bash-completion
```

##### add the following to your ~/.bashrc file
```shell
source /usr/share/bash-completion/bash_completion
```

##### Enable kubectl autocompletion
```shell
echo 'source <(kubectl completion bash)' >>~/.bashrc
```

##### Alias for kubectl
```shell
echo 'alias k=kubectl' >>~/.bashrc
echo 'complete -F __start_kubectl k' >>~/.bashrc
```

### bash auto-completion on macOS
```shell
bash_complation@1 --> bash@3.2
bash_completion@2 --> bash@4.1+
```

##### Upgrade Bash
```shell
echo $BASH_VERSION
brew install bash
```

##### Install bash-completion
```shell
brew install bash-completion@2
```

##### Store the configuration for the automcompletion 
```shell
kubectl completion bash > ~/.kube/kubectl_autocompletion
```

##### add the following to your ~/.bash_p file
```shell
if [ -f /usr/local/share/bash-completion/bash_completion ]; then
. /usr/local/share/bash-completion/bash_completion
fi
source ~/.kube/kubectl_autocompletion
```

After reloading your shell, kubectl autocompletion should be working.