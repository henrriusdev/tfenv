# Detect architecture
$ARCH = $env:PROCESSOR_ARCHITECTURE
if ($ARCH -eq "AMD64") {
    $ARCH = "amd64"
} elseif ($ARCH -eq "ARM64") {
    $ARCH = "arm64"
} else {
    Write-Host "❌ Unsupported architecture: $ARCH"
    exit 1
}

# Set download URL based on architecture
$BINARY_NAME = "tfenv-windows-$ARCH.exe"
$DOWNLOAD_URL = "https://github.com/henrriusdev/tfenv/releases/latest/download/$BINARY_NAME"
$DEST_PATH = "C:\Program Files\tfenv\tfenv.exe"

# Create directory if it doesn't exist
if (!(Test-Path "C:\Program Files\tfenv")) {
    New-Item -ItemType Directory -Path "C:\Program Files\tfenv" | Out-Null
}

# Download the correct binary
Write-Host "🔽 Downloading $BINARY_NAME from $DOWNLOAD_URL..."
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $DEST_PATH

# Ensure it's executable
Set-ExecutionPolicy Unrestricted -Scope Process -Force
Write-Host "✅ Installation complete! Run 'tfenv' to use the CLI."
