<#
Sync scripture JSON data from bcbooks/scriptures-json into internal/scripture/data for embedding.
Usage:
  pwsh ./sync-data.ps1
Requires: git in PATH.
#>
param()

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

Write-Host "Syncing scripture data from bcbooks/scriptures-json repository..."

$temp = New-Item -ItemType Directory -Path ([System.IO.Path]::GetTempPath()) -Name ("scriptures-json-" + [guid]::NewGuid())
try {
    Push-Location $temp.FullName
    Write-Host "Cloning scriptures-json repository..."
    git clone https://github.com/bcbooks/scriptures-json.git | Out-Null
    Set-Location scriptures-json

    $repoRoot = Split-Path -Parent (Split-Path -Parent $PSCommandPath) # go two up from script (repo root)
    $dataDir = Join-Path $repoRoot 'internal/scripture/data'
    if (-not (Test-Path $dataDir)) { New-Item -ItemType Directory -Path $dataDir | Out-Null }

    $files = @(
        'book-of-mormon.json',
        'doctrine-and-covenants.json',
        'pearl-of-great-price.json',
        'old-testament.json',
        'new-testament.json'
    )
    foreach ($f in $files) {
        Copy-Item -LiteralPath $f -Destination (Join-Path $dataDir $f) -Force
    }

    Write-Host "Scripture data synchronized successfully!" -ForegroundColor Green
    Write-Host "Data files updated in: $dataDir"
    Get-ChildItem -File $dataDir/*.json | Select-Object Name,Length | Format-Table -AutoSize
    Write-Host "`nTo see what changed, run: git diff --stat internal/scripture/data/" -ForegroundColor Yellow
}
finally {
    Pop-Location 2>$null
    if (Test-Path $temp) { Remove-Item -Recurse -Force $temp }
}
