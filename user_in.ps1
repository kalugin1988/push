
$login = $env:username
$device_name = $env:computername
$DOMAIN = $env:USERDOMAIN
$time = [datetime]::Now;
$mac = (Get-NetAdapter -Physical | Where-Object {$_.Status -eq "Up"} | Select-Object MacAddress, IPv4Address) | ForEach-Object {"$($_.MacAddress)"}
$string = "Logged on $computer at $time";

if (-not ($login -match '[a-zA-Z]')) {
    $login = 'student'
  }
#$login = $login -replace 'Ученик','Extybr' 
& "C:\push.exe" --job Login --instance $device_name --install "$login=1" --url https://push.it25.su

