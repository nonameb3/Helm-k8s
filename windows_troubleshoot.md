# Windows Troubleshooting Guide

This guide provides Windows-specific instructions for the health-service Kubernetes deployment.

## Prerequisites for Windows

### Required Software
```powershell
# Install using winget (Windows Package Manager)
winget install Docker.DockerDesktop
winget install Kubernetes.kubectl
winget install Helm.Helm
winget install Git.Git

# Or install manually:
# - Docker Desktop: https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe
# - kubectl: https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/
# - Helm: https://github.com/helm/helm/releases
```

### Recommended Terminal
Use **Git Bash** or **PowerShell** for best compatibility. Git Bash provides Unix-like commands.

## Windows-Specific Command Replacements

### DNS Configuration (Local Development)

#### PowerShell (Run as Administrator)
```powershell
# Add DNS entries
Add-Content -Path "C:\Windows\System32\drivers\etc\hosts" -Value "127.0.0.1 health-service-dev.local"
Add-Content -Path "C:\Windows\System32\drivers\etc\hosts" -Value "127.0.0.1 health-service-staging.local"
Add-Content -Path "C:\Windows\System32\drivers\etc\hosts" -Value "127.0.0.1 health-service-prod.local"

# Verify entries were added
Get-Content "C:\Windows\System32\drivers\etc\hosts" | Select-String "health-service"
```

#### Command Prompt (Run as Administrator)
```cmd
echo 127.0.0.1 health-service-dev.local >> C:\Windows\System32\drivers\etc\hosts
echo 127.0.0.1 health-service-staging.local >> C:\Windows\System32\drivers\etc\hosts
echo 127.0.0.1 health-service-prod.local >> C:\Windows\System32\drivers\etc\hosts
```

### DNS Cleanup

#### PowerShell (Run as Administrator)
```powershell
# Remove health-service DNS entries
$hostsPath = "C:\Windows\System32\drivers\etc\hosts"
(Get-Content $hostsPath) | Where-Object {$_ -notmatch 'health-service.*\.local'} | Set-Content $hostsPath

# Verify cleanup
Get-Content $hostsPath | Select-String "health-service"
```

### Load Testing Commands

#### PowerShell
```powershell
# Multiple concurrent requests
1..6 | ForEach-Object {
    Start-Job -ScriptBlock { curl http://health-service-dev.local/load/20000 }
}

# Wait for jobs and get results
Get-Job | Wait-Job | Receive-Job
Get-Job | Remove-Job
```

#### Git Bash (Recommended)
```bash
# Works exactly like Linux/Mac
for i in {1..6}; do
  curl http://health-service-dev.local/load/20000 &
done
```

## Windows-Specific Tools

### Apache Bench Alternative
```powershell
# Install using Chocolatey
choco install apache-httpd

# Or download from Apache Lounge
# https://www.apachelounge.com/download/

# Usage (same as Linux/Mac)
ab -n 10000 -c 50 http://health-service-dev.local/health
```

### Alternative Load Testing with PowerShell
```powershell
# Simple load test script
$url = "http://health-service-dev.local/health"
$requests = 100
$concurrent = 10

1..$concurrent | ForEach-Object {
    Start-Job -ScriptBlock {
        param($url, $count)
        for ($i = 1; $i -le $count; $i++) {
            try {
                Invoke-WebRequest -Uri $url -Method GET -TimeoutSec 30
            } catch {
                Write-Host "Request failed: $_"
            }
        }
    } -ArgumentList $url, ($requests / $concurrent)
}

Get-Job | Wait-Job | Receive-Job
Get-Job | Remove-Job
```

## Common Windows Issues

### Issue 1: kubectl not found
```powershell
# Check if kubectl is in PATH
kubectl version --client

# If not found, add to PATH or use full path
$env:PATH += ";C:\Program Files\kubectl"
```

### Issue 2: Docker not running
```powershell
# Check Docker status
docker version

# Start Docker Desktop if not running
Start-Process "C:\Program Files\Docker\Docker\Docker Desktop.exe"
```

### Issue 3: Permission denied on hosts file
```powershell
# Always run PowerShell/CMD as Administrator for hosts file changes
# Right-click -> "Run as administrator"
```

### Issue 4: Port forwarding not working
```powershell
# Check if ports are available
netstat -an | findstr :8080
netstat -an | findstr :80

# Kill process using port if needed
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

### Issue 5: Git Bash path issues
```bash
# In Git Bash, use Unix-style paths
cd /c/Users/username/project

# Convert Windows paths
winpty docker run -it ubuntu  # Add winpty for interactive containers
```

## File Path Conversions

| Linux/Mac | Windows PowerShell | Git Bash |
|-----------|-------------------|----------|
| `/etc/hosts` | `C:\Windows\System32\drivers\etc\hosts` | `/c/Windows/System32/drivers/etc/hosts` |
| `/home/user` | `C:\Users\username` | `/c/Users/username` |
| `./script.sh` | `.\script.ps1` | `./script.sh` |

## Recommended Windows Workflow

1. **Use Git Bash** for most commands (closest to Linux/Mac experience)
2. **Use PowerShell as Administrator** for system-level changes (hosts file, services)
3. **Use Docker Desktop** GUI for container management
4. **Use Windows Terminal** for better terminal experience

## Testing the Setup

```powershell
# Test all components
kubectl cluster-info
helm version
docker version

# Test DNS resolution
nslookup health-service-dev.local
ping health-service-dev.local

# Test endpoints
curl http://health-service-dev.local/health
curl http://health-service-staging.local/health
curl http://health-service-prod.local/health
```

## Emergency Fixes

### Reset hosts file
```powershell
# Backup current hosts file
Copy-Item "C:\Windows\System32\drivers\etc\hosts" "C:\Windows\System32\drivers\etc\hosts.backup"

# Reset to default (PowerShell as Administrator)
@"
# Copyright (c) 1993-2009 Microsoft Corp.
#
# This is a sample HOSTS file used by Microsoft TCP/IP for Windows.
#
127.0.0.1       localhost
::1             localhost
"@ | Set-Content "C:\Windows\System32\drivers\etc\hosts"
```

### Reset Docker Desktop
```powershell
# Reset Docker to factory settings
# Docker Desktop -> Settings -> Reset -> Reset to factory defaults
```