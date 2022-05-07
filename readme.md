# This is an example of Golang x AWS SDK
This Golang program first login an AWS service account. Then, the program try to read aws's ec2 information.

---
## Install golang 1.17 version
```bash
# Install Golang binary
wget https://go.dev/dl/go1.17.5.linux-amd64.tar.gz
tar -zxvf go1.17.5.linux-amd64.tar.gz -C /usr/local/
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bash_profile
mkdir /go
echo "export GOPATH=/go" >> ~/.bash_profile
source ~/.bash_profile
rm -f go1.17.5.linux-amd64.tar.gz
```

---
## Quick Start
```bash 
git clone https://github.com/alexshinningsun/read-aws-ec2.git
cd read-aws-ec2
go mod download
go build -o app .
./app
```