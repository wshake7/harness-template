package orm

import (
	"admin/internal/config"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"database/sql"
	"go.uber.org/zap"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"
	gormCrud "orm-crud/gormc"
)

var Client *gormCrud.Client // 包级变量

func DB() *gorm.DB {
	return Client.DB
}

func New(config config.OrmConfig) *gormCrud.Client {
	var options []gormCrud.Option
	zapLogger := zap.L()
	l := zapgorm2.New(zapLogger)
	l.SetAsDefault()
	if config.IsLog {
		l.LogLevel = logger.Info
	}
	options = append(options,
		gormCrud.WithLogger(zapLogger.Sugar()),
		gormCrud.WithGormConfig(&gorm.Config{Logger: l, TranslateError: true}),
		gormCrud.WithDriverName(config.DriverName),
		gormCrud.WithDSN(config.DataSource),
		gormCrud.WithEnableTrace(true),
		gormCrud.WithEnableMetrics(true),
	)

	if config.IsAutoMigrate {
		options = append(options, gormCrud.WithAutoMigrate(models.Models...))
	}

	var err error
	Client, err = gormCrud.NewClient(options...)
	if err != nil {
		zapLogger.Error("init gorm client error", zap.Error(err))
		panic(err)
	}
	gromDb := Client.DB
	sqlDb, err := gromDb.DB()
	if err != nil {
		zapLogger.Error("init gorm client error", zap.Error(err))
		panic(err)
	}
	initSqlDB(sqlDb, config)
	if config.IsGenCode {
		genCode(gromDb, models.Models...)
	}
	query.SetDefault(gromDb)
	return Client
}

func initSqlDB(sqlDB *sql.DB, config config.OrmConfig) {
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	sqlDB.SetMaxIdleConns(config.MaxIdleConn)
	sqlDB.SetMaxOpenConns(config.MaxOpenConn)
}

func genCode(db *gorm.DB, modelSlice ...any) {
	cfg := gen.Config{
		OutPath:           "./fiberc/services/repo/query",
		OutFile:           "",
		ModelPkgPath:      "",
		WithUnitTest:      false,
		FieldNullable:     true,
		FieldCoverable:    true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
	}
	// 驼峰
	//cfg.WithJSONTagNameStrategy(func(columnName string) (tagContent string) {
	//	return strcase.LowerCamelCase(columnName)
	//})
	g := gen.NewGenerator(cfg)
	g.UseDB(db)
	g.ApplyBasic(modelSlice...)
	g.Execute()
}
