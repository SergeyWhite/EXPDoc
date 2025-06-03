package psql

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Database struct {
	DB *sql.DB
}

// DatabaseConfig содержит конфигурацию для подключения к базе данных
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewDatabase создает новое подключение к PostgreSQL
func NewDatabase(config DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройки пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	return &Database{DB: db}, nil
}

// Close закрывает подключение к базе данных
func (d *Database) Close() error {
	if d.DB != nil {
		return d.DB.Close()
	}
	return nil
}

// Migration представляет одну миграцию
type Migration struct {
	Version int
	Name    string
	Content string
}

// createMigrationsTable создает таблицу для отслеживания миграций
func (d *Database) createMigrationsTable() error {
	query := `
  CREATE TABLE IF NOT EXISTS schema_migrations (
   version INTEGER PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  );
 `
	_, err := d.DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	return nil
}

// getAppliedMigrations возвращает список примененных миграций
func (d *Database) getAppliedMigrations() (map[int]bool, error) {
	applied := make(map[int]bool)

	rows, err := d.DB.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}

	return applied, nil
}

// loadMigrations загружает миграции из папки migrations
func (d *Database) loadMigrations(migrationsDir string) ([]Migration, error) {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []Migration

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}

		// Парсим имя файла для получения версии
		// Ожидаемый формат: 001_create_users_table.sql
		parts := strings.SplitN(file.Name(), "_", 2)
		if len(parts) != 2 {
			log.Printf("Skipping file with invalid name format: %s", file.Name())
			continue
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			log.Printf("Skipping file with invalid version number: %s", file.Name())
			continue
		}

		// Читаем содержимое файла
		filePath := filepath.Join(migrationsDir, file.Name())
		content, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		// Получаем имя миграции (убираем расширение .sql)
		name := strings.TrimSuffix(file.Name(), ".sql")

		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			Content: string(content),
		})
	}

	// Сортируем миграции по версии
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// RunMigrations выполняет все неприменённые миграции из папки migrations
func (d *Database) RunMigrations(migrationsDir string) error {
	// Создаем таблицу миграций если её нет
	if err := d.createMigrationsTable(); err != nil {
		return err
	}
	// Загружаем миграции из файлов
	migrations, err := d.loadMigrations(migrationsDir)
	if err != nil {
		return err
	}

	if len(migrations) == 0 {
		log.Println("No migration files found")
		return nil
	}

	// Получаем список примененных миграций
	applied, err := d.getAppliedMigrations()
	if err != nil {
		return err
	}

	// Выполняем неприменённые миграции
	for _, migration := range migrations {
		if applied[migration.Version] {
			log.Printf("Migration %d (%s) already applied, skipping", migration.Version, migration.Name)
			continue
		}

		log.Printf("Applying migration %d: %s", migration.Version, migration.Name)

		// Начинаем транзакцию
		tx, err := d.DB.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", migration.Version, err)
		}

		// Выполняем миграцию
		if _, err := tx.Exec(migration.Content); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %d (%s): %w", migration.Version, migration.Name, err)
		}

		// Записываем информацию о применённой миграции
		if _, err := tx.Exec("INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
			migration.Version, migration.Name); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		// Подтверждаем транзакцию
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
		}

		log.Printf("Successfully applied migration %d: %s", migration.Version, migration.Name)
	}

	log.Println("All migrations completed successfully")
	return nil
}

// GetMigrationStatus возвращает статус миграций
func (d *Database) GetMigrationStatus(migrationsDir string) error {
	migrations, err := d.loadMigrations(migrationsDir)
	if err != nil {
		return err
	}

	applied, err := d.getAppliedMigrations()
	if err != nil {
		return err
	}

	log.Println("Migration Status:")
	log.Println("================")

	for _, migration := range migrations {
		status := "PENDING"
		if applied[migration.Version] {
			status = "APPLIED"
		}
		log.Printf("Version %d: %s [%s]", migration.Version, migration.Name, status)
	}

	return nil
}

// InitializeDatabase инициализирует подключение к базе данных с настройками из контейнера
func InitializeDatabase() (*Database, error) {
	config := DatabaseConfig{
		Host:     "localhost", // или "postgres" если запускается в Docker Compose
		Port:     "5432",
		User:     "admin",
		Password: "admin",
		DBName:   "appdb",
		SSLMode:  "disable",
	}

	return NewDatabase(config)
}

// Пример использования
func Init_db() {
	// Инициализируем подключение к базе данных
	db, err := InitializeDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Запускаем миграции
	migrationsDir := "./migrations"
	if err := db.RunMigrations(migrationsDir); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Показываем статус миграций
	if err := db.GetMigrationStatus(migrationsDir); err != nil {
		log.Fatal("Failed to get migration status:", err)
	}

	log.Println("Database initialized successfully!")
}
