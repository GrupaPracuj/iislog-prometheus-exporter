$binPath = $PSScriptRoot + "/bin"
$exePath = $binPath + '/IISLogExporter.exe'
iex 'go get -v'
iex 'go get gopkg.in/natefinch/lumberjack.v2'
$buildCommand = 'go build  -o ' + $exePath
iex $buildCommand
$copyCommand = 'Copy-Item *.hcl ' + $binPath
iex $copyCommand