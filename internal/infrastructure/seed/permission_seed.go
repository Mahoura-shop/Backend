package seed

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type PermissionSeeder struct {
	db database.Database
}

func NewPermissionSeeder(db database.Database) *PermissionSeeder {
	return &PermissionSeeder{db: db}
}

var predefinedPermissions = []entity.Permission{
	{Name: "product:create", Description: "ایجاد محصول", Category: "product"},
	{Name: "product:edit", Description: "ویرایش محصول", Category: "product"},
	{Name: "product:delete", Description: "حذف محصول", Category: "product"},
	{Name: "product:see", Description: "مشاهده محصولات", Category: "product"},
	{Name: "product:batch_price", Description: "به‌روزرسانی گروهی قیمت", Category: "product"},
	{Name: "product:batch_inventory", Description: "به‌روزرسانی گروهی موجودی", Category: "product"},
	{Name: "category:see", Description: "مشاهده دسته‌بندی‌ها", Category: "category"},
	{Name: "category:create", Description: "ایجاد دسته‌بندی", Category: "category"},
	{Name: "category:edit", Description: "ویرایش دسته‌بندی", Category: "category"},
	{Name: "category:delete", Description: "حذف دسته‌بندی", Category: "category"},
	{Name: "brand:see", Description: "مشاهده برندها", Category: "brand"},
	{Name: "brand:create", Description: "ایجاد برند", Category: "brand"},
	{Name: "brand:edit", Description: "ویرایش برند", Category: "brand"},
	{Name: "brand:delete", Description: "حذف برند", Category: "brand"},
	{Name: "order:see", Description: "مشاهده سفارشات", Category: "order"},
	{Name: "order:update_status", Description: "به‌روزرسانی وضعیت سفارش", Category: "order"},
	{Name: "order:cancel", Description: "لغو سفارش", Category: "order"},
	{Name: "users:see", Description: "مشاهده کاربران", Category: "users"},
	{Name: "users:edit_role", Description: "تغییر نقش کاربر", Category: "users"},
	{Name: "users:wallet", Description: "مدیریت کیف پول کاربر", Category: "users"},
	{Name: "users:ban", Description: "مسدود کردن کاربر", Category: "users"},
	{Name: "users:unban", Description: "رفع مسدودیت کاربر", Category: "users"},
	{Name: "rbac:create", Description: "ایجاد نقش", Category: "rbac"},
	{Name: "rbac:edit", Description: "ویرایش نقش", Category: "rbac"},
	{Name: "rbac:see", Description: "مشاهده نقش‌ها", Category: "rbac"},
	{Name: "rbac:delete", Description: "حذف نقش", Category: "rbac"},
	{Name: "contact:see", Description: "مشاهده پیام‌های تماس", Category: "contact"},
	{Name: "update:currencies", Description: "به‌روزرسانی ارزها", Category: "update"},
	{Name: "adminlogs:see", Description: "مشاهده لاگ‌های ادمین", Category: "adminlogs"},
}

func (s *PermissionSeeder) SeedPermissions() {
	for _, perm := range predefinedPermissions {
		var existing entity.Permission
		result := s.db.GetDB().Where("name = ?", perm.Name).First(&existing)
		if result.Error != nil {
			s.db.GetDB().Create(&perm)
		} else if existing.Category != perm.Category || existing.Description != perm.Description {
			s.db.GetDB().Model(&existing).Updates(entity.Permission{
				Description: perm.Description,
				Category:    perm.Category,
			})
		}
	}
}
