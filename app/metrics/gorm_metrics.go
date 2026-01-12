package metrics

import (
    "time"

    "gorm.io/gorm"
)

func RegisterGormCallbacks(db *gorm.DB) {
    db.Callback().Create().Before("gorm:create").Register("metrics:before_create", beforeCallback())
    db.Callback().Create().After("gorm:create").Register("metrics:after_create", afterCallback("create"))

    db.Callback().Query().Before("gorm:query").Register("metrics:before_query", beforeCallback())
    db.Callback().Query().After("gorm:query").Register("metrics:after_query", afterCallback("query"))

    db.Callback().Update().Before("gorm:update").Register("metrics:before_update", beforeCallback())
    db.Callback().Update().After("gorm:update").Register("metrics:after_update", afterCallback("update"))

    db.Callback().Delete().Before("gorm:delete").Register("metrics:before_delete", beforeCallback())
    db.Callback().Delete().After("gorm:delete").Register("metrics:after_delete", afterCallback("delete"))
}

func beforeCallback() func(*gorm.DB) {
    return func(db *gorm.DB) {
        db.InstanceSet("metrics:start_time", time.Now())
    }
}

func afterCallback(operation string) func(*gorm.DB) {
    return func(db *gorm.DB) {
        startTime, ok := db.InstanceGet("metrics:start_time")
        if !ok {
            return
        }

        duration := time.Since(startTime.(time.Time)).Seconds()
        table := db.Statement.Table
        if table == "" {
            table = "unknown"
        }

        DatabaseQueriesTotal.WithLabelValues(operation, table).Inc()
        DatabaseQueryDuration.WithLabelValues(operation, table).Observe(duration)
    }
}