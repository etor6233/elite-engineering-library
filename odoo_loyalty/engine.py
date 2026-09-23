# SPDX-License-Identifier: LGPL-3.0-only
# Derived from Odoo Community 19.0, commit 99edb6dd82b7b560930c00b03b694ba700785370.
# Original copyright and full LGPLv3/GPLv3 texts: COPYRIGHT and LICENSE.
# Modified for an isolated exact-decimal calculation module; see DERIVATION.json.
from collections import defaultdict
from decimal import Decimal, ROUND_DOWN, ROUND_HALF_UP, ROUND_HALF_EVEN

def _(message, **values):
    return message % values if values else message

class _Date:
    @staticmethod
    def today():
        # The contract rejects currency conversion; this argument is unused by same-currency adapter.
        return None

class fields:
    Date = _Date

def float_round(value, precision_digits=None, precision_rounding=None, rounding_method='HALF-UP'):
    # AUTHORED exact representation adapter; business rounding modes come from the mapped source.
    if (precision_digits is None) == (precision_rounding is None):
        raise ValueError('one rounding precision is required')
    step = Decimal(str(precision_rounding)) if precision_rounding is not None else Decimal(10) ** -precision_digits
    if step <= 0 or not step.is_finite():
        raise ValueError('invalid rounding precision')
    modes = {'DOWN': ROUND_DOWN, 'HALF-UP': ROUND_HALF_UP, 'HALF-EVEN': ROUND_HALF_EVEN}
    if rounding_method not in modes:
        raise ValueError('unsupported rounding mode')
    return (value / step).quantize(Decimal(1), rounding=modes[rounding_method]) * step

def _get_point_changes(self):
    """
    Returns the changes in points per coupon as a dict.

    Used when validating/cancelling an order
    """
    points_per_coupon = defaultdict(lambda: 0)
    for coupon_point in self.coupon_point_ids:
        points_per_coupon[coupon_point.coupon_id] += coupon_point.points
    for line in self.order_line:
        if not line.reward_id or not line.coupon_id:
            continue
        points_per_coupon[line.coupon_id] -= line.points_cost
    return points_per_coupon


def _get_real_points_for_coupon(self, coupon, post_confirm=False):
    """
    Returns the actual points usable for this coupon for this order. Set pos_confirm to True to include points for future orders.

    This is calculated by taking the points on the coupon, the points the order will give to the coupon (if applicable) and removing the points taken by already applied rewards.
    """
    self.ensure_one()
    points = coupon.points
    if self.state not in ('sale', 'done'):
        if coupon.program_id.applies_on != 'future':
            # Points that will be given by the order upon confirming the order
            points += self.coupon_point_ids.filtered(lambda p: p.coupon_id == coupon).points
        # Points already used by rewards
        points -= sum(self.order_line.filtered(lambda l: l.coupon_id == coupon).mapped('points_cost'))
    if any(rule.reward_point_mode == 'money' for rule in coupon.program_id.rule_ids):
        points = coupon.currency_id.round(points)
    return points


def _program_check_compute_points(self, programs):
    """
    Checks the program validity from the order lines aswell as computing the number of points to add.

    Returns a dict containing the error message or the points that will be given with the keys 'points'.
    """
    self.ensure_one()

    # Prepare quantities
    order_lines = self._get_not_rewarded_order_lines().filtered(
        lambda line: not line.combo_item_id
    )
    products = order_lines.product_id
    products_qties = dict.fromkeys(products, 0)
    for line in order_lines:
        product_qty = line.product_uom_id._compute_quantity(
            line.product_uom_qty, line.product_id.uom_id
        )
        products_qties[line.product_id] += product_qty
    # Contains the products that can be applied per rule
    products_per_rule = programs._get_valid_products(products)

    # Prepare amounts
    so_products_per_rule = programs._get_valid_products(self.order_line.product_id)
    lines_per_rule = defaultdict(lambda: self.env['sale.order.line'])
    # Skip lines that have no effect on the minimum amount to reach.
    for line in self.order_line - self._get_no_effect_on_threshold_lines():
        is_discount = line.reward_id.reward_type == 'discount'
        reward_program = line.reward_id.program_id
        # Skip lines for automatic discounts, as well as combo item lines.
        if (is_discount and reward_program.trigger == 'auto') or line.combo_item_id:
            continue
        for program in programs:
            # Skip lines for the current program's discounts.
            if is_discount and reward_program == program:
                continue
            for rule in program.rule_ids:
                # Skip lines to which the rule doesn't apply.
                if line.product_id in so_products_per_rule.get(rule, []):
                    lines_per_rule[rule] |= line._get_lines_with_price()

    result = {}
    for program in programs:
        # Used for error messages
        # By default False, but True if no rules and applies_on current -> misconfigured coupons program
        code_matched = not bool(program.rule_ids) and program.applies_on == 'current' # Stays false if all triggers have code and none have been activated
        minimum_amount_matched = code_matched
        product_qty_matched = code_matched
        points = 0
        # Some rules may split their points per unit / money spent
        #  (i.e. gift cards 2x50$ must result in two 50$ codes)
        rule_points = []
        program_result = result.setdefault(program, dict())
        for rule in program.rule_ids:
            # prevent bottomless ewallet spending
            if program.program_type == 'ewallet' and not program.trigger_product_ids:
                break
            if rule.mode == 'with_code' and rule not in self.code_enabled_rule_ids:
                continue
            code_matched = True
            rule_amount = rule._compute_amount(self.currency_id)
            untaxed_amount = sum(lines_per_rule[rule].mapped('price_subtotal'))
            tax_amount = sum(lines_per_rule[rule].mapped('price_tax'))
            if rule_amount > (rule.minimum_amount_tax_mode == 'incl' and (untaxed_amount + tax_amount) or untaxed_amount):
                continue
            minimum_amount_matched = True
            if not products_per_rule.get(rule):
                continue
            rule_products = products_per_rule[rule]
            ordered_rule_products_qty = sum(products_qties[product] for product in rule_products)
            if ordered_rule_products_qty < rule.minimum_qty or not rule_products:
                continue
            product_qty_matched = True
            if not rule.reward_point_amount:
                continue
            # Count all points separately if the order is for the future and the split option is enabled
            if program.applies_on == 'future' and rule.reward_point_split and rule.reward_point_mode != 'order':
                if rule.reward_point_mode == 'unit':
                    rule_points.extend(rule.reward_point_amount for _ in range(int(ordered_rule_products_qty)))
                elif rule.reward_point_mode == 'money':
                    for line in self.order_line:
                        if (
                            line.is_reward_line
                            or line.combo_item_id
                            or line.product_id not in rule_products
                            or line.product_uom_qty <= 0
                        ):
                            continue
                        line_price_total = self._get_order_line_price(line, 'price_total')
                        points_per_unit = float_round(
                            (rule.reward_point_amount * line_price_total / line.product_uom_qty),
                            precision_digits=2, rounding_method='DOWN')
                        if not points_per_unit:
                            continue
                        rule_points.extend([points_per_unit] * int(line.product_uom_qty))
            else:
                # All checks have been passed we can now compute the points to give
                if rule.reward_point_mode == 'order':
                    points += rule.reward_point_amount
                elif rule.reward_point_mode == 'money':
                    # Compute amount paid for rule
                    # NOTE: this accounts for discounts -> 1 point per $ * (100$ - 30%) will
                    # result in 70 points
                    amount_paid = Decimal('0.0')
                    rule_products = so_products_per_rule.get(rule, [])
                    for line in self.order_line - self._get_no_effect_on_threshold_lines():
                        if line.combo_item_id or line.reward_id.program_id.program_type in [
                            'ewallet', 'gift_card', program.program_type
                        ]:
                            continue
                        line_price_total = self._get_order_line_price(line, 'price_total')
                        amount_paid += (
                            line_price_total if line.product_id in rule_products
                            else Decimal('0.0')
                        )

                    points += float_round(rule.reward_point_amount * amount_paid, precision_digits=2, rounding_method='DOWN')
                elif rule.reward_point_mode == 'unit':
                    points += rule.reward_point_amount * ordered_rule_products_qty
        # NOTE: for programs that are nominative we always allow the program to be 'applied' on the order
        #  with 0 points so that `_get_claimable_rewards` returns the rewards associated with those programs
        if not program.is_nominative:
            if not code_matched:
                program_result['error'] = _("This program requires a code to be applied.")
            elif not minimum_amount_matched:
                program_result['error'] = _(
                    "To take advantage of this offer, your order must include at least %(amount)s %(currency)s of the eligible products.",
                    amount=min(program.rule_ids.mapped('minimum_amount')),
                    currency=program.currency_id.name,
                )
            elif not product_qty_matched:
                program_result['error'] = _("You don't have the required product quantities on your sales order.")
        elif self.partner_id.is_public and not self._allow_nominative_programs():
            program_result['error'] = _("This program is not available for public users.")
        if 'error' not in program_result:
            points_result = [points] + rule_points
            program_result['points'] = points_result
    return result


def reward_from_bound_amounts(self, reward, coupon, discountable):
    reward_currency = reward.currency_id
    reward_program = reward.program_id
    max_discount = reward_currency._convert(reward.discount_max_amount, self.currency_id, self.company_id, fields.Date.today()) or Decimal('Infinity')
    # discount should never surpass the order's current total amount
    max_discount = min(self.amount_total, max_discount)
    if reward.discount_mode == 'per_point':
        points = self._get_real_points_for_coupon(coupon)
        if not reward_program.is_payment_program:
            # Rewards cannot be partially offered to customers
            points = points // reward.required_points * reward.required_points
        max_discount = min(max_discount,
            reward_currency._convert(reward.discount * points,
                self.currency_id, self.company_id, fields.Date.today()))
    elif reward.discount_mode == 'per_order':
        max_discount = min(max_discount,
            reward_currency._convert(reward.discount, self.currency_id, self.company_id, fields.Date.today()))
    elif reward.discount_mode == 'percent':
        max_discount = min(max_discount, discountable * (reward.discount / 100))

    # Discount per taxes
    point_cost = reward.required_points if not reward.clear_wallet else self._get_real_points_for_coupon(coupon)
    if reward.discount_mode == 'per_point' and not reward.clear_wallet:
        # Calculate the actual point cost if the cost is per point
        converted_discount = self.currency_id._convert(min(max_discount, discountable), reward_currency, self.company_id, fields.Date.today())
        point_cost = coupon.currency_id.round(converted_discount / reward.discount)

    return {"value": min(max_discount, discountable), "points_cost": point_cost}



# ADAPTED sale_order.py:1138-1140. Container arguments replace ORM attributes.
def filtered_issuance(status, program):
    all_point_changes = [p for p in status['points'] if p]
    if not all_point_changes and program.is_nominative:
        all_point_changes = [0]
    return all_point_changes

# ADAPTED exact cancellation arithmetic from _action_cancel: coupon.points -= changes.
# Argument/return plumbing replaces an ORM property update; immutable journal
# delta is computed against zero rather than mutating or deleting historical rows.
def reverse_point_change(points, changes):
    points -= changes
    return points
