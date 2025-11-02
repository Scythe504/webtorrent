$repo = "scythe504/fluxstream"
$installDir = "$env:LOCALAPPDATA\Programs\FluxStream"
$tempDir = "$env:TEMP"
New-Item -ItemType Directory -Force -Path $installDir | Out-Null

Write-Host "Detecting latest version..."
$latest = (Invoke-RestMethod "https://api.github.com/repos/$repo/releases/latest").tag_name
Write-Host "Latest version: $latest"

# Detect architecture
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "arm64" }
$file = "fluxstream_${latest}_windows_${arch}.zip"
$url = "https://github.com/$repo/releases/download/$latest/$file"

Write-Host "Downloading $file ..."
Invoke-WebRequest -Uri $url -OutFile "$tempDir\$file"

Write-Host "Extracting to $installDir ..."
Expand-Archive "$tempDir\$file" -DestinationPath $installDir -Force

# Optionally add to PATH
if ($env:PATH -notlike "*$installDir*") {
    Write-Host "Adding FluxStream to PATH"
    [Environment]::SetEnvironmentVariable("PATH", "$env:PATH;$installDir", "User")
}

Write-Host "FluxStream installed successfully!"
& "$installDir\fluxstream.exe" --version
