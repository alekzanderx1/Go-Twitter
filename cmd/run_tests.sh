curl -L http://127.0.0.1:12380/ping -XPUT -d 'ping'
curl -L http://127.0.0.1:12380/tweets -XPUT -d ' {"chi.bharathsai@gmail.com":[{"Text":"Hi, this is my first tweet!","CreatedBy":"chi.bharathsai@gmail.com","CreatedTimestamp":"2022-12-13 22:27:18"},{"Text":"Argentina won, yay!","CreatedBy":"chi.bharathsai@gmail.com","CreatedTimestamp":"2022-12-13 22:27:30"}],"idk":[{"Text":"Who am I?","CreatedBy":"idk","CreatedTimestamp":"2022-12-14 12:58:53"},{"Text":"I am spamming..","CreatedBy":"idk","CreatedTimestamp":"2022-12-15 09:24:01"}],"test":[{"Text":"This app is awesome!","CreatedBy":"test","CreatedTimestamp":"2022-12-18 11:19:39"}],"tjmax":[{"Text":"Hello World!","CreatedBy":"tjmax","CreatedTimestamp":"2022-12-18 11:20:11"}]}'
curl -L http://127.0.0.1:12380/users -XPUT -d '{"chi.bharathsai@gmail.com":{"Username":"chi.bharathsai@gmail.com","Name":"Bharath","Password":"!Y4qJiEfvs6uP26","Following":{}},"sid":{"Username":"sid","Name":"sid","Password":"123","Following":{}},"syedahmad":{"Username":"syedahmad","Name":"Syed","Password":"123","Following":{}},"test":{"Username":"test","Name":"test","Password":"123","Following":{}},"tjmax":{"Username":"tjmax","Name":"Tejas","Password":"123","Following":{}}}'

cd ../web/users/
go test
cd ../tweets
go test
cd ../authentication
go test
echo "All tests passed! OK"