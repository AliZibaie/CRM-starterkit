package modulesconfig

import (
	"crm-starterkit/internal/core/module"
	"crm-starterkit/internal/modules/announcement"
	"crm-starterkit/internal/modules/auditlog"
	"crm-starterkit/internal/modules/auth"
	"crm-starterkit/internal/modules/banner"
	"crm-starterkit/internal/modules/cart"
	"crm-starterkit/internal/modules/catalog"
	"crm-starterkit/internal/modules/cms"
	"crm-starterkit/internal/modules/crm"
	"crm-starterkit/internal/modules/cronbuilder"
	"crm-starterkit/internal/modules/discount"
	"crm-starterkit/internal/modules/media"
	"crm-starterkit/internal/modules/notification"
	"crm-starterkit/internal/modules/order"
	"crm-starterkit/internal/modules/payment"
	"crm-starterkit/internal/modules/rbac"
	"crm-starterkit/internal/modules/requestlog"
	"crm-starterkit/internal/modules/review"
	"crm-starterkit/internal/modules/setting"
	"crm-starterkit/internal/modules/sms"
	"crm-starterkit/internal/modules/support"
	"crm-starterkit/internal/modules/user"
	"crm-starterkit/internal/modules/wishlist"
)

var EnabledModules = []module.Module{
	rbac.New(),
	user.New(),
	auth.New(),
	cronbuilder.New(),
	sms.New(),
	notification.New(),
	announcement.New(),
	support.New(),
	banner.New(),
	catalog.New(),
	cart.New(),
	discount.New(),
	order.New(),
	payment.New(),
	media.New(),
	cms.New(),
	review.New(),
	wishlist.New(),
	crm.New(),
	setting.New(),
	auditlog.New(),
	requestlog.New(),
}
