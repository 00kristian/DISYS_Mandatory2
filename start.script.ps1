Start-Process cmd -ArgumentList "/C go run . Y :1000 :2000 :3000"
Start-Process cmd -ArgumentList "/C go run . n :2000 :3000 :1000"
Start-Process cmd -ArgumentList "/C go run . n :3000 :1000 :2000"