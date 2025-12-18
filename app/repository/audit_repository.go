package repository

import "backend-path/app/models"

type AuditRepository struct {}

func (r *AuditRepository) Create(auditLog *models.AuditLog) error {
	return DB.Create(auditLog).Error
}