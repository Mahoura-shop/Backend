package enum

type PermissionType uint
type PermissionCategory uint

const (
	PermissionAll PermissionType = iota + 1
	PermissionGeneral

	// Admin Role Permissions
	// User Management
	UserViewAll
	UserBanUnban
	UserChangeRole
	UserViewRoles
	UserManageRolePermissions
	UserRemoveRole
	UserCreateRole

	// Ticket Management
	TicketViewAll
	TicketRespond
	TicketClose
	TicketComment

	// Profile Management
	ProfileViewPrivate
	ProfileUpdate
)

const (
	CategoryGeneral PermissionCategory = iota + 1
	CategoryUser
	CategoryTicket
	CategoryReport
	CategoryProfile
)

var permissionNames = map[PermissionType]string{
	PermissionAll:     "general.all",
	PermissionGeneral: "general.general",

	// User Management
	UserViewAll:               "user.viewAll",
	UserBanUnban:              "user.banUnban",
	UserChangeRole:            "user.changeRole",
	UserViewRoles:             "user.viewRoles",
	UserManageRolePermissions: "user.manageRolePermissions",
	UserRemoveRole:            "user.removeRole",
	UserCreateRole:            "user.createRole",

	// Ticket Management
	TicketViewAll: "ticket.viewAll",
	TicketRespond: "ticket.respond",
	TicketClose:   "ticket.close",
	TicketComment: "ticket.comment",

	// Profile Management
	ProfileViewPrivate: "profile.viewPrivate",
	ProfileUpdate:      "profile.update",
}

var permissionDescriptions = map[PermissionType]string{
	PermissionAll:     "دسترسی کامل به سیستم",
	PermissionGeneral: "دسترسی عمومی",

	// User Management
	UserViewAll:               "مشاهده لیست کاربران وب‌سایت",
	UserBanUnban:              "مسدود/رفع مسدودیت کاربران",
	UserChangeRole:            "تغییر نقش کاربر",
	UserViewRoles:             "مشاهده لیست نقش‌ها و مجوزها",
	UserManageRolePermissions: "تغییر مجوزهای نقش",
	UserRemoveRole:            "حذف نقش",
	UserCreateRole:            "ایجاد نقش جدید",

	// Ticket Management
	TicketViewAll: "مشاهده لیست تیکت‌ها",
	TicketRespond: "پاسخ به تیکت",
	TicketClose:   "بستن تیکت",
	TicketComment: "افزودن نظر به تیکت",

	// Profile Management
	ProfileViewPrivate: "مشاهده اطلاعات خصوصی پروفایل",
	ProfileUpdate:      "به‌روزرسانی پروفایل",
}

var permissionCategories = map[PermissionType]PermissionCategory{
	PermissionAll:     CategoryGeneral,
	PermissionGeneral: CategoryGeneral,

	// User Management
	UserViewAll:               CategoryUser,
	UserBanUnban:              CategoryUser,
	UserChangeRole:            CategoryUser,
	UserViewRoles:             CategoryUser,
	UserManageRolePermissions: CategoryUser,
	UserRemoveRole:            CategoryUser,
	UserCreateRole:            CategoryUser,

	// Ticket Management
	TicketViewAll: CategoryTicket,
	TicketRespond: CategoryTicket,
	TicketClose:   CategoryTicket,
	TicketComment: CategoryTicket,

	// Profile Management
	ProfileViewPrivate: CategoryProfile,
	ProfileUpdate:      CategoryProfile,
}

func (perm PermissionType) String() string {
	if description, ok := permissionNames[perm]; ok {
		return description
	}
	return ""
}

func (perm PermissionType) Description() string {
	if description, ok := permissionDescriptions[perm]; ok {
		return description
	}
	return ""
}

func (perm PermissionType) Category() PermissionCategory {
	if category, ok := permissionCategories[perm]; ok {
		return category
	}
	return CategoryGeneral
}

func (category PermissionCategory) String() string {
	switch category {
	case CategoryGeneral:
		return "عمومی"
	case CategoryUser:
		return "مدیریت کاربران"
	case CategoryTicket:
		return "مدیریت تیکت"
	case CategoryReport:
		return "مدیریت گزارشات"
	case CategoryProfile:
		return "مدیریت پروفایل"
	}
	return "unknown"
}

func GetAllPermissionTypes() []PermissionType {
	return []PermissionType{
		PermissionAll, PermissionGeneral,

		// User Management
		UserViewAll, UserBanUnban, UserChangeRole, UserViewRoles,
		UserManageRolePermissions, UserRemoveRole, UserCreateRole,

		// Ticket Management
		TicketViewAll, TicketRespond, TicketClose, TicketComment,

		// Profile Management
		ProfileViewPrivate, ProfileUpdate,
	}
}

func GetAdminPermissionTypes() []PermissionType {
	return []PermissionType{
		PermissionAll, PermissionGeneral,

		// User Management
		UserViewAll, UserBanUnban, UserChangeRole, UserViewRoles,
		UserManageRolePermissions, UserRemoveRole, UserCreateRole,

		// Ticket Management
		TicketViewAll, TicketRespond, TicketClose, TicketComment,
	}
}
