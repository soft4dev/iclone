# Install script for your Go CLI tool (Windows)
# Usage: irm https://raw.githubusercontent.com/your-username/your-repo/main/install.ps1 | iex

$ErrorActionPreference = 'Stop'

# Configuration - UPDATE THESE VALUES
$Repo = "soft4dev/clonei"  # GitHub repository (e.g., "owner/repo")
$BinName = "clonei"              # Binary name (e.g., "mycli")

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "x86_64" }
        "ARM64" { return "arm64" }
        "x86" { return "i386" }
        default {
            Write-Host "Unsupported architecture: $arch" -ForegroundColor Red
            exit 1
        }
    }
}

# Get the latest release version
function Get-LatestVersion {
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
        return $response.tag_name
    }
    catch {
        Write-Host "Error: Could not fetch latest version" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
        exit 1
    }
}

# Main installation function
function Install-Binary {
    $arch = Get-Architecture
    $version = Get-LatestVersion
    
    Write-Host "Installing $BinName $version..." -ForegroundColor Green

    # Construct download URL
    $archiveName = "${BinName}_Windows_${arch}.zip"
    $downloadUrl = "https://github.com/$Repo/releases/download/$version/$archiveName"
    
    # Create temporary directory
    $tmpDir = Join-Path $env:TEMP ([System.IO.Path]::GetRandomFileName())
    New-Item -ItemType Directory -Path $tmpDir | Out-Null
    
    $archivePath = Join-Path $tmpDir $archiveName
    
    try {
        Write-Host "Downloading from $downloadUrl..." -ForegroundColor Yellow
        
        # Download the archive
        Invoke-WebRequest -Uri $downloadUrl -OutFile $archivePath -UseBasicParsing
        
        Write-Host "Extracting..." -ForegroundColor Yellow
        
        # Extract the archive
        Expand-Archive -Path $archivePath -DestinationPath $tmpDir -Force
        
        # Determine install directory
        $installDir = if ($env:BIN_DIR) { $env:BIN_DIR } else { Join-Path $env:USERPROFILE ".local\bin" }
        
        # Create install directory if it doesn't exist
        if (-not (Test-Path $installDir)) {
            New-Item -ItemType Directory -Path $installDir | Out-Null
        }
        
        # Move binary to install directory
        $binaryName = "$BinName.exe"
        $sourcePath = Join-Path $tmpDir $binaryName
        $destPath = Join-Path $installDir $binaryName
        
        Write-Host "Installing to $installDir..." -ForegroundColor Yellow
        
        # Remove existing binary if present
        if (Test-Path $destPath) {
            Remove-Item $destPath -Force
        }
        
        Move-Item -Path $sourcePath -Destination $destPath -Force
        
        Write-Host "✓ $BinName $version installed successfully!" -ForegroundColor Green
        Write-Host ""
        
        # Check if install directory is in PATH
        $userPath = [Environment]::GetEnvironmentVariable("Path", "User")
        if ($userPath -notlike "*$installDir*") {
            Write-Host "Note: Adding $installDir to your PATH..." -ForegroundColor Yellow
            
            # Add to PATH
            $newPath = if ($userPath) { "$userPath;$installDir" } else { $installDir }
            [Environment]::SetEnvironmentVariable("Path", $newPath, "User")
            
            Write-Host "✓ Added to PATH. Restart your terminal for changes to take effect." -ForegroundColor Green
            Write-Host ""
        }
        
        Write-Host "Run '$BinName --help' to get started" -ForegroundColor Green
        Write-Host "(You may need to restart your terminal first)" -ForegroundColor Yellow
    }
    catch {
        Write-Host "Installation failed: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
    finally {
        # Cleanup
        if (Test-Path $tmpDir) {
            Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
        }
    }
}

# Run installation
Install-Binary