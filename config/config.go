package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	Cache     CacheConfig
	JWT       JWTConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
	Storage   StorageConfig
	Log       LogConfig
	Worker    WorkerConfig
	Queue     QueueConfig
	Migrate   MigrateConfig
	Payment   PaymentConfig
	SMS       SMSConfig
	Email     EmailConfig
}

type SMSConfig struct {
	Provider       string
	KavenegarAPIKey string
	KavenegarSender string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	FromAddress  string
}

type PaymentConfig struct {
	CallbackBaseURL  string
	ZarinpalMerchant string
	ZarinpalSandbox  bool
	StripeSecretKey  string
	StripeWebhookKey string
	StripeCurrency   string
}

type AppConfig struct {
	Name  string
	Env   string
	Debug bool
	URL   string
	Port  int
}

type DatabaseConfig struct {
	Driver     string
	Host       string
	Port       int
	Name       string
	User       string
	Password   string
	SSLMode    string
	Timezone   string
	SQLitePath string
}

type CacheConfig struct {
	Driver      string
	Redis       RedisConfig
	Memcached   MemcachedConfig
	Prefix      string
	DefaultTTL  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type MemcachedConfig struct {
	Host string
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string
}

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

type RateLimitConfig struct {
	Enabled bool
	Max     int
	Window  time.Duration
}

type StorageConfig struct {
	Driver        string
	LocalPath     string
	PublicBaseURL string
	S3            S3Config
}

type S3Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	Endpoint        string
	ForcePathStyle  bool
	PublicBaseURL   string
}

type LogConfig struct {
	Level    string
	Encoding string
}

type WorkerConfig struct {
	Concurrency int
}

type QueueConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

type MigrateConfig struct {
	Path string
}

func Load(path string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType("env")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			if !fileMissing(err) {
				return nil, fmt.Errorf("read config: %w", err)
			}
		}
	}

	cfg := &Config{
		App: AppConfig{
			Name:  v.GetString("APP_NAME"),
			Env:   v.GetString("APP_ENV"),
			Debug: v.GetBool("APP_DEBUG"),
			URL:   v.GetString("APP_URL"),
			Port:  v.GetInt("APP_PORT"),
		},
		Database: DatabaseConfig{
			Driver:     strings.ToLower(v.GetString("DB_DRIVER")),
			Host:       v.GetString("DB_HOST"),
			Port:       v.GetInt("DB_PORT"),
			Name:       v.GetString("DB_NAME"),
			User:       v.GetString("DB_USER"),
			Password:   v.GetString("DB_PASSWORD"),
			SSLMode:    v.GetString("DB_SSLMODE"),
			Timezone:   v.GetString("DB_TIMEZONE"),
			SQLitePath: v.GetString("DB_SQLITE_PATH"),
		},
		Cache: CacheConfig{
			Driver: strings.ToLower(v.GetString("CACHE_DRIVER")),
			Redis: RedisConfig{
				Host:     v.GetString("REDIS_HOST"),
				Port:     v.GetInt("REDIS_PORT"),
				Password: v.GetString("REDIS_PASSWORD"),
				DB:       v.GetInt("REDIS_DB"),
			},
			Memcached: MemcachedConfig{
				Host: v.GetString("MEMCACHED_HOST"),
			},
			Prefix:     v.GetString("CACHE_PREFIX"),
			DefaultTTL: v.GetDuration("CACHE_DEFAULT_TTL"),
		},
		JWT: JWTConfig{
			Secret:     v.GetString("JWT_SECRET"),
			AccessTTL:  v.GetDuration("JWT_ACCESS_TTL"),
			RefreshTTL: v.GetDuration("JWT_REFRESH_TTL"),
			Issuer:     v.GetString("JWT_ISSUER"),
		},
		CORS: CORSConfig{
			AllowOrigins:     splitCSV(v.GetString("CORS_ALLOW_ORIGINS")),
			AllowMethods:     splitCSV(v.GetString("CORS_ALLOW_METHODS")),
			AllowHeaders:     splitCSV(v.GetString("CORS_ALLOW_HEADERS")),
			AllowCredentials: v.GetBool("CORS_ALLOW_CREDENTIALS"),
			MaxAge:           v.GetInt("CORS_MAX_AGE"),
		},
		RateLimit: RateLimitConfig{
			Enabled: v.GetBool("RATE_LIMIT_ENABLED"),
			Max:     v.GetInt("RATE_LIMIT_MAX"),
			Window:  v.GetDuration("RATE_LIMIT_WINDOW"),
		},
		Storage: StorageConfig{
			Driver:        strings.ToLower(v.GetString("STORAGE_DRIVER")),
			LocalPath:     v.GetString("STORAGE_LOCAL_PATH"),
			PublicBaseURL: v.GetString("STORAGE_PUBLIC_BASE_URL"),
			S3: S3Config{
				Region:          v.GetString("AWS_REGION"),
				AccessKeyID:     v.GetString("AWS_ACCESS_KEY_ID"),
				SecretAccessKey: v.GetString("AWS_SECRET_ACCESS_KEY"),
				Bucket:          v.GetString("AWS_S3_BUCKET"),
				Endpoint:        v.GetString("AWS_S3_ENDPOINT"),
				ForcePathStyle:  v.GetBool("AWS_S3_FORCE_PATH_STYLE"),
				PublicBaseURL:   v.GetString("AWS_S3_PUBLIC_BASE_URL"),
			},
		},
		Log: LogConfig{
			Level:    v.GetString("LOG_LEVEL"),
			Encoding: v.GetString("LOG_ENCODING"),
		},
		Worker: WorkerConfig{
			Concurrency: v.GetInt("WORKER_CONCURRENCY"),
		},
		Queue: QueueConfig{
			RedisAddr:     v.GetString("QUEUE_REDIS_ADDR"),
			RedisPassword: v.GetString("QUEUE_REDIS_PASSWORD"),
			RedisDB:       v.GetInt("QUEUE_REDIS_DB"),
		},
		Migrate: MigrateConfig{
			Path: v.GetString("MIGRATIONS_PATH"),
		},
		Payment: PaymentConfig{
			CallbackBaseURL:  v.GetString("PAYMENT_CALLBACK_BASE_URL"),
			ZarinpalMerchant: v.GetString("ZARINPAL_MERCHANT_ID"),
			ZarinpalSandbox:  v.GetBool("ZARINPAL_SANDBOX"),
			StripeSecretKey:  v.GetString("STRIPE_SECRET_KEY"),
			StripeWebhookKey: v.GetString("STRIPE_WEBHOOK_SECRET"),
			StripeCurrency:   strings.ToLower(v.GetString("STRIPE_CURRENCY")),
		},
		SMS: SMSConfig{
			Provider:        strings.ToLower(v.GetString("SMS_PROVIDER")),
			KavenegarAPIKey: v.GetString("KAVENEGAR_API_KEY"),
			KavenegarSender: v.GetString("KAVENEGAR_SENDER"),
		},
		Email: EmailConfig{
			SMTPHost:     v.GetString("SMTP_HOST"),
			SMTPPort:     v.GetInt("SMTP_PORT"),
			SMTPUser:     v.GetString("SMTP_USER"),
			SMTPPassword: v.GetString("SMTP_PASSWORD"),
			FromAddress:  v.GetString("SMTP_FROM"),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.App.Port <= 0 {
		return fmt.Errorf("APP_PORT must be > 0")
	}
	if c.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	switch c.Database.Driver {
	case "postgres", "mysql", "sqlserver", "sqlite":
	default:
		return fmt.Errorf("unsupported DB_DRIVER: %s", c.Database.Driver)
	}
	switch c.Cache.Driver {
	case "redis", "memcached":
	default:
		return fmt.Errorf("unsupported CACHE_DRIVER: %s", c.Cache.Driver)
	}
	switch c.Storage.Driver {
	case "local", "s3":
	default:
		return fmt.Errorf("unsupported STORAGE_DRIVER: %s", c.Storage.Driver)
	}
	return nil
}

func (c *Config) Address() string {
	return fmt.Sprintf(":%d", c.App.Port)
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("APP_NAME", "CRM-StarterKit")
	v.SetDefault("APP_ENV", "development")
	v.SetDefault("APP_DEBUG", true)
	v.SetDefault("APP_URL", "http://localhost:8080")
	v.SetDefault("APP_PORT", 8080)

	v.SetDefault("DB_DRIVER", "postgres")
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_NAME", "starterkit")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASSWORD", "secret")
	v.SetDefault("DB_SSLMODE", "disable")
	v.SetDefault("DB_TIMEZONE", "UTC")
	v.SetDefault("DB_SQLITE_PATH", "./storage/database.sqlite")

	v.SetDefault("CACHE_DRIVER", "redis")
	v.SetDefault("REDIS_HOST", "localhost")
	v.SetDefault("REDIS_PORT", 6379)
	v.SetDefault("REDIS_PASSWORD", "")
	v.SetDefault("REDIS_DB", 0)
	v.SetDefault("MEMCACHED_HOST", "localhost:11211")
	v.SetDefault("CACHE_PREFIX", "crm:")
	v.SetDefault("CACHE_DEFAULT_TTL", "1h")

	v.SetDefault("JWT_SECRET", "change-me-to-a-long-random-secret")
	v.SetDefault("JWT_ACCESS_TTL", "15m")
	v.SetDefault("JWT_REFRESH_TTL", "168h")
	v.SetDefault("JWT_ISSUER", "crm-starterkit")

	v.SetDefault("CORS_ALLOW_ORIGINS", "*")
	v.SetDefault("CORS_ALLOW_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	v.SetDefault("CORS_ALLOW_HEADERS", "Origin,Content-Type,Accept,Authorization")
	v.SetDefault("CORS_ALLOW_CREDENTIALS", false)
	v.SetDefault("CORS_MAX_AGE", 86400)

	v.SetDefault("RATE_LIMIT_ENABLED", true)
	v.SetDefault("RATE_LIMIT_MAX", 100)
	v.SetDefault("RATE_LIMIT_WINDOW", "1m")

	v.SetDefault("STORAGE_DRIVER", "local")
	v.SetDefault("STORAGE_LOCAL_PATH", "./storage/uploads")
	v.SetDefault("STORAGE_PUBLIC_BASE_URL", "http://localhost:8080/uploads")

	v.SetDefault("AWS_REGION", "us-east-1")
	v.SetDefault("AWS_S3_FORCE_PATH_STYLE", false)

	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("LOG_ENCODING", "json")

	v.SetDefault("WORKER_CONCURRENCY", 10)
	v.SetDefault("QUEUE_REDIS_ADDR", "localhost:6379")
	v.SetDefault("QUEUE_REDIS_PASSWORD", "")
	v.SetDefault("QUEUE_REDIS_DB", 1)

	v.SetDefault("MIGRATIONS_PATH", "./migrations")

	v.SetDefault("PAYMENT_CALLBACK_BASE_URL", "http://localhost:8080/api/v1/payments/callback")
	v.SetDefault("ZARINPAL_MERCHANT_ID", "")
	v.SetDefault("ZARINPAL_SANDBOX", true)
	v.SetDefault("STRIPE_SECRET_KEY", "")
	v.SetDefault("STRIPE_WEBHOOK_SECRET", "")
	v.SetDefault("STRIPE_CURRENCY", "usd")

	v.SetDefault("SMS_PROVIDER", "kavenegar")
	v.SetDefault("KAVENEGAR_API_KEY", "")
	v.SetDefault("KAVENEGAR_SENDER", "")
	v.SetDefault("SMTP_HOST", "")
	v.SetDefault("SMTP_PORT", 587)
	v.SetDefault("SMTP_USER", "")
	v.SetDefault("SMTP_PASSWORD", "")
	v.SetDefault("SMTP_FROM", "")
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func fileMissing(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "no such file") || strings.Contains(msg, "cannot find the file") || strings.Contains(msg, "The system cannot find")
}
