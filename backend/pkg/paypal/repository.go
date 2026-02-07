package paypal

import (
	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/clarkzhu2020/aidecms/pkg/database"
	"gorm.io/gorm"
)

// Repository 支付数据仓库
type Repository struct {
	db *gorm.DB
}

// NewRepository 创建支付仓库
func NewRepository() *Repository {
	return &Repository{
		db: database.GetDB(),
	}
}

// CreatePayment 创建支付记录
func (r *Repository) CreatePayment(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

// GetPaymentByID 根据ID获取支付记录
func (r *Repository) GetPaymentByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetPaymentByOrderID 根据订单ID获取支付记录
func (r *Repository) GetPaymentByOrderID(orderID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetPaymentByPayPalOrderID 根据PayPal订单ID获取支付记录
func (r *Repository) GetPaymentByPayPalOrderID(paypalOrderID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("paypal_order_id = ?", paypalOrderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetPaymentByCaptureID 根据捕获ID获取支付记录
func (r *Repository) GetPaymentByCaptureID(captureID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("capture_id = ?", captureID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// ListPayments 获取支付列表
func (r *Repository) ListPayments(offset, limit int, status string) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	query := r.db.Model(&models.Payment{})

	// 状态过滤
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// UpdatePayment 更新支付记录
func (r *Repository) UpdatePayment(payment *models.Payment) error {
	return r.db.Save(payment).Error
}

// UpdatePaymentStatus 更新支付状态
func (r *Repository) UpdatePaymentStatus(id uint, status string) error {
	return r.db.Model(&models.Payment{}).Where("id = ?", id).Update("status", status).Error
}

// UpdatePaymentCaptureID 更新支付捕获ID
func (r *Repository) UpdatePaymentCaptureID(id uint, captureID string) error {
	return r.db.Model(&models.Payment{}).Where("id = ?", id).Update("capture_id", captureID).Error
}

// DeletePayment 删除支付记录（软删除）
func (r *Repository) DeletePayment(id uint) error {
	return r.db.Delete(&models.Payment{}, id).Error
}

// ListPaymentsByUserID 根据用户ID获取支付列表
func (r *Repository) ListPaymentsByUserID(userID uint, offset, limit int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	query := r.db.Model(&models.Payment{}).Where("user_id = ?", userID)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}

// CreateRefund 创建退款记录
func (r *Repository) CreateRefund(refund *models.PaymentRefund) error {
	return r.db.Create(refund).Error
}

// GetRefundByRefundID 根据退款ID获取退款记录
func (r *Repository) GetRefundByRefundID(refundID string) (*models.PaymentRefund, error) {
	var refund models.PaymentRefund
	err := r.db.Where("refund_id = ?", refundID).First(&refund).Error
	if err != nil {
		return nil, err
	}
	return &refund, nil
}

// ListRefundsByPaymentID 根据支付ID获取退款列表
func (r *Repository) ListRefundsByPaymentID(paymentID uint) ([]models.PaymentRefund, error) {
	var refunds []models.PaymentRefund
	err := r.db.Where("payment_id = ?", paymentID).Order("created_at DESC").Find(&refunds).Error
	if err != nil {
		return nil, err
	}
	return refunds, nil
}

// CreateWebhook 创建Webhook记录
func (r *Repository) CreateWebhook(webhook *models.PaymentWebhook) error {
	return r.db.Create(webhook).Error
}

// GetWebhookByEventID 根据事件ID获取Webhook记录
func (r *Repository) GetWebhookByEventID(eventID string) (*models.PaymentWebhook, error) {
	var webhook models.PaymentWebhook
	err := r.db.Where("event_id = ?", eventID).First(&webhook).Error
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// ListWebhooksByResourceID 根据资源ID获取Webhook列表
func (r *Repository) ListWebhooksByResourceID(resourceID string, limit int) ([]models.PaymentWebhook, error) {
	var webhooks []models.PaymentWebhook
	err := r.db.Where("resource_id = ?", resourceID).
		Order("created_at DESC").
		Limit(limit).
		Find(&webhooks).Error
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

// GetPaymentStats 获取支付统计信息
func (r *Repository) GetPaymentStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总支付金额
	var totalAmount float64
	r.db.Model(&models.Payment{}).
		Where("status = ?", "paid").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalAmount)
	stats["total_amount"] = totalAmount

	// 支付总数
	var totalCount int64
	r.db.Model(&models.Payment{}).
		Where("status = ?", "paid").
		Count(&totalCount)
	stats["total_count"] = totalCount

	// 按状态统计
	var statusStats []struct {
		Status string
		Count  int64
	}
	r.db.Model(&models.Payment{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&statusStats)

	statusMap := make(map[string]int64)
	for _, stat := range statusStats {
		statusMap[stat.Status] = stat.Count
	}
	stats["by_status"] = statusMap

	// 按货币统计
	var currencyStats []struct {
		Currency string
		Count    int64
		Amount   float64
	}
	r.db.Model(&models.Payment{}).
		Where("status = ?", "paid").
		Select("currency, COUNT(*) as count, SUM(amount) as amount").
		Group("currency").
		Scan(&currencyStats)

	currencyMap := make(map[string]interface{})
	for _, stat := range currencyStats {
		currencyMap[stat.Currency] = map[string]interface{}{
			"count":  stat.Count,
			"amount": stat.Amount,
		}
	}
	stats["by_currency"] = currencyMap

	return stats, nil
}

// SearchPayments 搜索支付记录
func (r *Repository) SearchPayments(keyword string, offset, limit int) ([]models.Payment, int64, error) {
	var payments []models.Payment
	var total int64

	query := r.db.Model(&models.Payment{}).Where(
		"order_id LIKE ? OR paypal_order_id LIKE ? OR description LIKE ?",
		"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%",
	)

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&payments).Error
	if err != nil {
		return nil, 0, err
	}

	return payments, total, nil
}
