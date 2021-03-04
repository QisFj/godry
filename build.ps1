Set-Location $PSScriptRoot

# clean outputdirectory
if (Test-Path ./output)
{
    Remove-Item ./output -Recurse -Force
}

# go build
Get-ChildItem ./cmd -directory | ForEach-Object {
    Write-Output "<-- start building $_"
    go build -v -o ./output/$_.exe ./cmd/$_
    if ($LASTEXITCODE -ne 0)
    {
        Write-Output " --> build $_ failed"
        exit 1
    }
    Write-Output " --> build $_ success"
}