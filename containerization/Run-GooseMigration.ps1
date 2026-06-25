param(
    [Parameter(Mandatory = $true)]
    [string]$Path,

    [Parameter(Mandatory = $true)]
    [ValidateSet("Up", "Down")]
    [string]$Type,

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
}

class EnvironmentParser {
    [DatabaseConfig] Parse([string]$path) {
        $lines = [System.IO.File]::ReadAllLines($path)
        $config = [DatabaseConfig]::new()
        foreach ($line in $lines) {
            $this.ProcessLine($line, $config)
        }
        if (-not $config.RootPassword -or -not $config.Database -or -not $config.Port) {
            throw "Missing required environment variables: DB_ROOT_PASSWORD, DB_DATABASE, DB_PORT"
        }
        return $config
    }

    hidden [void] ProcessLine([string]$line, [DatabaseConfig]$config) {
        switch -Regex ($line) {
            '^\s*DB_ROOT_PASSWORD\s*=\s*(.*)$' { $config.RootPassword = $matches[1].Trim(); break }
            '^\s*DB_DATABASE\s*=\s*(.*)$'      { $config.Database = $matches[1].Trim(); break }
            '^\s*DB_PORT\s*=\s*(.*)$'          { $config.Port = $matches[1].Trim(); break }
        }
    }
}

class ConnectionStringBuilder {
    static [string] Build([DatabaseConfig]$config) {
        return "root:$($config.RootPassword)@tcp(127.0.0.1:$($config.Port))/$($config.Database)?parseTime=true"
    }
}

class ArgumentBuilder {
    static [string[]] Build([string]$migrationDir, [string]$connectionString, [MigrationAction]$action) {
        $actionString = if ($action -eq [MigrationAction]::Up) { 'up' } else { 'down' }
        return @('-dir', $migrationDir, 'mysql', $connectionString, $actionString)
    }
}

class GooseExecutor {
    [int] TryExecute([string]$executable, [string[]]$arguments) {
        try {
            & $executable $arguments
            if ($LASTEXITCODE -ne 0) {
                throw "goose exited with code $LASTEXITCODE"
            }
            return 0
        } catch {
            Write-Error "Execution failed: $($_.Exception.Message)"
            return 1
        }
    }
}

function Start-MigrationProcess {
    param(
        [string]$migrationPath,
        [MigrationAction]$migrationType,
        [string]$envPath
    )
    $parser = [EnvironmentParser]::new()
    $dbConfig = $parser.Parse($envPath)

    $connectionString = [ConnectionStringBuilder]::Build($dbConfig)
    $arguments = [ArgumentBuilder]::Build($migrationPath, $connectionString, $migrationType)

    $executor = [GooseExecutor]::new()
    $exitCode = $executor.TryExecute('goose', $arguments)
    if ($exitCode -ne 0) {
        exit $exitCode
    }
}

Start-MigrationProcess -migrationPath $Path -migrationType $Type -envPath $EnvFilePath