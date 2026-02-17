param(
  [string]$BaseUrl = "http://localhost:8080/api/v1",
  [string]$Token = "",
  [string]$VoyageId = "",
  [string]$CruiseId = "",
  [string]$CabinId = "",
  [string]$CabinTypeId = "",
  [string]$WechatpaySignature = "",
  [string]$WechatpayTimestamp = "",
  [string]$WechatpayNonce = "",
  [string]$WechatpaySerial = "",
  [string]$WechatCallbackBody = "{}",
  [switch]$SkipCallback
)

$ErrorActionPreference = "Stop"

function Ensure-Command {
  param([string]$Name)
  if (-not (Get-Command $Name -ErrorAction SilentlyContinue)) {
    throw "未找到命令 '$Name'，请先安装。"
  }
}

Ensure-Command "newman"

$root = Split-Path -Parent $PSScriptRoot
$collection = Join-Path $PSScriptRoot "CruiseBooking-Prelaunch.postman_collection.json"
$envFile = Join-Path $PSScriptRoot "CruiseBooking-Local.postman_environment.json"

if (-not (Test-Path $collection)) {
  throw "Collection 不存在: $collection"
}
if (-not (Test-Path $envFile)) {
  throw "Environment 不存在: $envFile"
}

$envJson = Get-Content $envFile -Raw | ConvertFrom-Json

function Set-EnvValue {
  param(
    [object]$Env,
    [string]$Key,
    [string]$Value
  )
  foreach ($item in $Env.values) {
    if ($item.key -eq $Key) {
      $item.value = $Value
      return
    }
  }
}

Set-EnvValue -Env $envJson -Key "baseUrl" -Value $BaseUrl
Set-EnvValue -Env $envJson -Key "token" -Value $Token
Set-EnvValue -Env $envJson -Key "voyageId" -Value $VoyageId
Set-EnvValue -Env $envJson -Key "cruiseId" -Value $CruiseId
Set-EnvValue -Env $envJson -Key "cabinId" -Value $CabinId
Set-EnvValue -Env $envJson -Key "cabinTypeId" -Value $CabinTypeId
Set-EnvValue -Env $envJson -Key "wechatpaySignature" -Value $WechatpaySignature
Set-EnvValue -Env $envJson -Key "wechatpayTimestamp" -Value $WechatpayTimestamp
Set-EnvValue -Env $envJson -Key "wechatpayNonce" -Value $WechatpayNonce
Set-EnvValue -Env $envJson -Key "wechatpaySerial" -Value $WechatpaySerial
Set-EnvValue -Env $envJson -Key "wechatCallbackBody" -Value $WechatCallbackBody

$tmpEnv = Join-Path $env:TEMP "CruiseBooking-Local.runtime.postman_environment.json"
$envJson | ConvertTo-Json -Depth 20 | Set-Content -Path $tmpEnv -Encoding UTF8

$folderArgs = @()
if ($SkipCallback) {
  $folderArgs += @("--folder", "0. Health")
  $folderArgs += @("--folder", "1. Order")
  $folderArgs += @("--folder", "2. Payment")
  $folderArgs += @("--folder", "4. Order Cancel")
}

Write-Host "[INFO] 开始执行 Newman..."
$newmanArgs = @(
  "run", $collection,
  "-e", $tmpEnv,
  "--reporters", "cli,junit",
  "--reporter-junit-export", (Join-Path $root "postman\newman-report.xml")
) + $folderArgs

& newman @newmanArgs
if ($LASTEXITCODE -ne 0) {
  throw "Newman 执行失败，退出码: $LASTEXITCODE"
}

Write-Host "[INFO] 执行完成，报告: postman/newman-report.xml"
