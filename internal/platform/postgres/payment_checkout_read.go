package postgres

// AUTHORED scoped read composition over existing payment/order owners.
import (
	"context"
	"net/url"
	"strings"
	"time"

	"elite.local/enterprise/internal/paymentbridge"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerCheckoutReader struct {
	pool  *pgxpool.Pool
	scope paymentbridge.Scope
}

func NewCustomerCheckoutReader(pool *pgxpool.Pool, scope paymentbridge.Scope) (*CustomerCheckoutReader, error) {
	if pool == nil || scope.Validate() != nil {
		return nil, paymentbridge.ErrCheckoutConflict
	}
	return &CustomerCheckoutReader{pool, scope}, nil
}
func (s *CustomerCheckoutReader) CustomerCheckout(ctx context.Context, tenant, organization, subject, order string) (paymentbridge.CustomerCheckout, error) {
	var result paymentbridge.CustomerCheckout
	if s == nil || s.pool == nil || tenant != s.scope.TenantID || organization != s.scope.OrganizationID || strings.TrimSpace(subject) == "" || order == "" {
		return result, paymentbridge.ErrCheckoutNotAvailable
	}
	rows, err := s.pool.Query(ctx, `select o.order_id,b.provider_code,b.checkout_url,b.expires_at
 from sales.customer_order o
 join crm.customer_profile customer on customer.tenant_id=o.tenant_id and customer.customer_principal_id=o.customer_principal_id and customer.status='active'
 join payment.payment_attempt p on p.tenant_id=o.tenant_id and p.order_id=o.order_id and p.provider_code=$5 and p.state='pending' and p.amount_minor_units=(select provider_due_minor_units from payment.order_funding f where f.tenant_id=o.tenant_id and f.order_id=o.order_id and f.organization_id=o.organization_id) and p.currency=o.currency
 join payment.provider_checkout b on b.tenant_id=p.tenant_id and b.payment_attempt_id=p.payment_attempt_id and b.provider_code=p.provider_code and b.connection_id=$6 and b.account_ref=$7 and b.live_mode=$8 and b.session_id is not null and b.checkout_url is not null and b.checkout_url<>'' and b.expires_at>clock_timestamp() and b.checkout_state in ('open','preference-created')
 join integration.provider_connection c on c.tenant_id=b.tenant_id and c.connection_id=b.connection_id and c.provider_code=b.provider_code and c.state='active' and (c.organization_id is null or c.organization_id=o.organization_id)
 where o.tenant_id=$1 and o.organization_id=$2 and o.customer_principal_id=$3 and o.order_id=$4 and o.state in ('placed','confirmed','allocated') and o.currency=$9
 limit 2`, tenant, organization, subject, order, s.scope.ProviderCode, s.scope.ConnectionID, s.scope.AccountRef, s.scope.LiveMode, s.scope.Currency)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	count := 0
	for rows.Next() {
		count++
		if err = rows.Scan(&result.OrderID, &result.ProviderCode, &result.URL, &result.ExpiresAt); err != nil {
			return paymentbridge.CustomerCheckout{}, err
		}
	}
	if err = rows.Err(); err != nil {
		return paymentbridge.CustomerCheckout{}, err
	}
	if count != 1 || !result.Valid(time.Now()) {
		return paymentbridge.CustomerCheckout{}, paymentbridge.ErrCheckoutNotAvailable
	}
	if s.scope.ProviderCode == "mercadopago" {
		u, _ := url.Parse(result.URL)
		expected := "sandbox.mercadopago.com.ar"
		if s.scope.LiveMode {
			expected = "www.mercadopago.com.ar"
		}
		if !strings.EqualFold(u.Hostname(), expected) {
			return paymentbridge.CustomerCheckout{}, paymentbridge.ErrCheckoutNotAvailable
		}
	}
	return result, nil
}
