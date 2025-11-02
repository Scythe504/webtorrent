# FluxStream Windows Installer
# Usage:
#   irm https://raw.githubusercontent.com/scythe504/fluxstream/main/scripts/install.ps1 | iex

Write-Host "Detecting latest version..."
$repo = "Scythe504/fluxstream"
$apiUrl = "https://api.github.com/repos/$repo/releases/latest"

try {
    $release = Invoke-RestMethod -Uri $apiUrl -UseBasicParsing
    $latest = $release.tag_name
    Write-Host "Latest version: $latest"
} catch {
    Write-Host "❌ Failed to detect latest version."
    exit 1
}

# Remove 'v' if tag is like 'v0.1.2'
$version = $latest -replace '^v', ''

$tempDir = [System.IO.Path]::GetTempPath()
$file = "fluxstream_${version}_windows_amd64.zip"
$url = "https://github.com/$repo/releases/download/$latest/$file"

Write-Host "Downloading $file ..."
try {
    Invoke-WebRequest -Uri $url -OutFile "$tempDir\$file"
} catch {
    Write-Host "❌ Download failed: $url"
    exit 1
}

$installDir = "$env:LOCALAPPDATA\Programs\FluxStream"
Write-Host "Extracting to $installDir ..."
if (!(Test-Path $installDir)) {
    New-Item -ItemType Directory -Force -Path $installDir | Out-Null
}

try {
    Expand-Archive "$tempDir\$file" -DestinationPath $installDir -Force
} catch {
    Write-Host "❌ Extraction failed."
    exit 1
}

Write-Host "Adding FluxStream to PATH..."
$envPath = [Environment]::GetEnvironmentVariable("Path", "User")
if ($envPath -notlike "*$installDir*") {
    [Environment]::SetEnvironmentVariable("Path", "$envPath;$installDir", "User")
}

Write-Host "✅ FluxStream installed successfully!"
Write-Host ""
try {
    & "$installDir\fluxstream.exe" --version
} catch {
    Write-Host "⚠️ FluxStream was installed but not immediately found in PATH. Restart your terminal and run 'fluxstream --version'."
}
