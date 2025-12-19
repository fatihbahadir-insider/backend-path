package repository

import "backend-path/app/models"

type IAuditLogRepository interface {
	Create(auditLog *models.AuditLog) error
}


type AuditRepository struct {}

func NewAuditRepository() *AuditRepository {
	return &AuditRepository{}
}

func (r *AuditRepository) Create(auditLog *models.AuditLog) error {
	return DB.Create(auditLog).Error
}