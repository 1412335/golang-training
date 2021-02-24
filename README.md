# golang-trainning

# Version Management
```
wget https://dl.google.com/go/go1.15.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.15.linux-amd64.tar.gz
mv /usr/local/go /usr/local/go1.15
mkdir -p /home/namnnp/go1.15

# vi ~/.profile 
vi ~/.zshrc 
###
# Golang v.1.15
export PATH=/usr/local/go1.15/bin:/home/namnnp/go1.15/bin:$PATH
export GOROOT=/usr/local/go1.15
export GOPATH=/home/namnnp/go1.15
###

source ~/.zshrc 
```

# Docs
- https://github.com/softrams/k6-docker