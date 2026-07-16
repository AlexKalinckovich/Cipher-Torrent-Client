<#
.SYNOPSIS
    Запуск миграций базы данных с помощью Goose.

.DESCRIPTION
    Этот скрипт применяет (Up) или откатывает (Down) миграции SQL для базы данных MySQL.

    ПАРАМЕТРЫ:
    - Type "Up": Применяет новые миграции.
    - Type "Down": Откатывает примененные миграции.
    - Target: Опциональный номер версии миграции.
        * Для Up: применить миграции до этой версии (включительно). Если 0 или не указано — применить все.
        * Для Down: откатить миграции до этой версии. Если не указано — откатить одну последнюю. Если 0 — откатить всё.

.PARAMETER Path
    Путь к директории с файлами миграций.

.PARAMETER Type
    Тип действия: "Up" для применения или "Down" для отката.

.PARAMETER Target
    Целевая версия миграции (число). Опционально.

.PARAMETER EnvFilePath
    Путь к файлу переменных окружения (.env). По умолчанию: "./containerization/.env".

.EXAMPLE
    # Применить все новые миграции
    .\Run-GooseMigration.ps1 -Path "./migrations" -Type "Up"

.EXAMPLE
    # Применить миграции только до версии 2
    .\Run-GooseMigration.ps1 -Path "./migrations" -Type "Up" -Target 2

.EXAMPLE
    # Откатить только последнюю миграцию
    .\Run-GooseMigration.ps1 -Path "./migrations" -Type "Down"

.EXAMPLE
    # Откатить ВСЕ миграции до нуля
    .\Run-GooseMigration.ps1 -Path "./migrations" -Type "Down" -Target 0
#>

param(
    [Parameter(Mandatory = $true)]
    [string]$Path,

    [Parameter(Mandatory = $true)]
    [ValidateSet("Up", "Down")]
    [string]$Type,

    [int]$Target = -1,

    [string]$EnvFilePath = "./containerization/.env"
)

enum MigrationAction {
    Up
    Down
}

$migrationType = [MigrationAction]::$Type

class DatabaseConfig {
    [string]$RootPassword
    [string]$Database
    [string]$Port
    [string]$Host
}

class EnvironmentParser {
    [DatabaseConfig] Parse([string]$path) {
        if (-not [System.IO.File]::Exists($path)) {
            throw "Environment file not found at path: $path"
        }

        $lines = [System.IO.File]::ReadAllLines($path)
        $config = [DatabaseConfig]::new()

        foreach ($line in $lines) {
            $this.ProcessLine($line, $config)
        }

        if (-not $config.RootPassword -or -not $config.Database -or -not $config.Port) {
            throw "Missing required environment variables: DB_ROOT_PASSWORD, DB_DATABASE, DB_PORT"
        }

        if (-not $config.Host) {
            $config.Host = "127.0.0.1"
        }

        return $config
    }

    hidden [void] ProcessLine([string]$line, [DatabaseConfig]$config) {
        switch -Regex ($line) {
            '^\s*DB_ROOT_PASSWORD\s*=\s*(.*)$' { $config.RootPassword = $matches[1].Trim(); break }
            '^\s*DB_DATABASE\s*=\s*(.*)$'      { $config.Database = $matches[1].Trim(); break }
            '^\s*DB_PORT\s*=\s*(.*)$'          { $config.Port = $matches[1].Trim(); break }
            '^\s*DB_HOST\s*=\s*(.*)$'          { $config.Host = $matches[1].Trim(); break }
        }
    }
}

class ConnectionStringBuilder {
    static [string] Build([DatabaseConfig]$config) {
        return "root:$($config.RootPassword)@tcp($($config.Host):$($config.Port))/$($config.Database)?parseTime=true"
    }
}

class ArgumentBuilder {
    static [string[]] Build([string]$migrationDir, [string]$connectionString, [MigrationAction]$action, [int]$target) {
        $baseArgs = @('-dir', $migrationDir, 'mysql', $connectionString)

        if ($action -eq [MigrationAction]::Up) {
            if ($target -ge 0) {
                return $baseArgs + @('up-to', $target.ToString())
            } else {
                return $baseArgs + @('up')
            }
        } else {
            if ($target -ge 0) {
                return $baseArgs + @('down-to', $target.ToString())
            } else {
                return $baseArgs + @('down')
            }
        }
    }
}

class GooseExecutor {
    [int] TryExecute([string]$executable, [string[]]$arguments) {
        Write-Host "Starting Goose migration..." -ForegroundColor Cyan
        Write-Host "Arguments: $($arguments -join ' ')" -ForegroundColor Gray

        try {
            & $executable $arguments

            if ($LASTEXITCODE -ne 0) {
                Write-Host "ERROR: Goose migration failed with exit code $LASTEXITCODE" -ForegroundColor Red
                return 1
            }

            Write-Host "SUCCESS: Goose migration completed successfully." -ForegroundColor Green
            return 0
        } catch {
            Write-Host "ERROR: Execution failed unexpectedly: $($_.Exception.Message)" -ForegroundColor Red
            return 1
        }
    }
}

function Start-MigrationProcess {
    param(
        [string]$migrationPath,
        [MigrationAction]$migrationType,
        [int]$target,
        [string]$envPath
    )

    try {
        $parser = [EnvironmentParser]::new()
        $dbConfig = $parser.Parse($envPath)

        $connectionString = [ConnectionStringBuilder]::Build($dbConfig)
        $arguments = [ArgumentBuilder]::Build($migrationPath, $connectionString, $migrationType, $target)

        $executor = [GooseExecutor]::new()
        $exitCode = $executor.TryExecute('goose', $arguments)

        if ($exitCode -ne 0) {
            exit $exitCode
        }
    } catch {
        Write-Host "FATAL ERROR: $($_.Exception.Message)" -ForegroundColor Red
        exit 1
    }
}

Start-MigrationProcess -migrationPath $Path -migrationType $Type -target $Target -envPath $EnvFilePath