# influxdb-test

An InfluxDB test server

## Run

Following command sequence is tested on an Digital Ocean Ubuntu 14.04 64bit droplet.

```
sudo apt-get update
sudo apt-get -y install git gcc npm nodejs
sudo ln -s /usr/bin/nodejs /usr/bin/node

wget http://get.influxdb.org/influxdb_0.9.0-rc30_amd64.deb
sudo dpkg -i influxdb_0.9.0-rc30_amd64.deb
sudo /etc/init.d/influxdb start

wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
tar zxvf go1.4.2.linux-amd64.tar.gz
mkdir ~/wd
cat - >> ~/.profile <<'EOF'
export GOROOT=~/go
export GOPATH=~/wd
export PATH=$GOROOT/bin:$GPPATH/bin:$PATH
EOF
source ~/.profile

cd $GOPATH/src/github.com/grafana/grafana
go run build.go setup
$GOPATH/bin/godep restore

npm install
npm install -g grunt-cli
grunt

go build .
```

`vi ~/wd/src/github.com/grafana/grafana/conf/defaults.ini`
Under `[auth.anonymous]`, modify `enabled = true`
Under `[security]`, Change `admin_user` and `admin_password`

```
curl -G http://localhost:8086/query --data-urlencode "q=CREATE DATABASE sonar"
curl -G http://localhost:8086/query --data-urlencode "q=CREATE USER sonar WITH PASSWORD 'sonar'"
grafana
```

```
cd $GOPATH/src/github.com/getlantern/influxdb-test/
go run main.go
```
