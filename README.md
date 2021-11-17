# DISYS_Mandatory2

## Guide to start up the application
### **Windows**

If you are on a windows computer:
- Open Windows powershell
- Navigate to directory of the assignment
- Run following command:
```
powershell -ExecutionPolicy bypass -file "start.script.ps1"
```
Note if any pop up windows about window firewall letting in a connection or something. Just press cancel, the token ring will run anyways

### **If you are running mac, linux or the above method does not work**

Unfortunately there's no experience in writing scripts for mac and no way for me to test it. So you would have to do it the old fashioned way through the command prompt / terminal

- Open three terminals
- Write the following (and do not execute yet):

Terminal 1:
```
go run . Y :1000 :2000 :3000
```
Terminal 2: 
```
go run . n :2000 :3000 :1000
```
Terminal 3: 
```
go run . Y :3000 :1000 :2000
```
- With rapid  and enter presses execute terminal 1 then 2 then 3 (Otherwise im pretty sure that you will get an error) 

Note: These steps might be easier if they were to be done in visual studio code. We've had the best succes there.
