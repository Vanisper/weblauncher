if (!(Test-Path ".env")) {
    Copy-Item ".env.example" ".env"
    Write-Host "Created .env from .env.example" -ForegroundColor Green
} else {
    Write-Host ".env already exists." -ForegroundColor Cyan
}

Write-Host ""
Write-Host "Tip: Create .env.local for local overrides (not tracked by git)" -ForegroundColor DarkGray
Write-Host "     Edit .env or .env.local with your configuration." -ForegroundColor Yellow
