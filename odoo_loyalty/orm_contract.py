# AUTHORED interface glue for the isolated source-derived calculations.
# This is a bounded recordset protocol, not Odoo ORM or an implementation of its database.
from decimal import Decimal, ROUND_HALF_UP
from types import MethodType


class NullRecord:
    def __bool__(self):
        return False
    def __getattr__(self, _name):
        return self
    def __eq__(self, other):
        return isinstance(other, NullRecord)
    def __hash__(self):
        return 0


NULL = NullRecord()


class Record:
    def __init__(self, kind, identity, **values):
        self._kind = kind
        self.id = identity
        self.__dict__.update(values)
    def __hash__(self):
        return hash((self._kind, self.id))
    def __eq__(self, other):
        return isinstance(other, Record) and (self._kind, self.id) == (other._kind, other.id)
    def ensure_one(self):
        return self


class Records:
    def __init__(self, values=()):
        self.values = tuple(dict.fromkeys(values))
    def __iter__(self):
        return iter(self.values)
    def __bool__(self):
        return bool(self.values)
    def __len__(self):
        return len(self.values)
    def __contains__(self, value):
        return value in self.values
    def __or__(self, other):
        return Records(self.values + other.values)
    def __sub__(self, other):
        return Records(x for x in self.values if x not in other)
    def filtered(self, condition):
        if isinstance(condition, str):
            return Records(x for x in self.values if getattr(x, condition))
        return Records(x for x in self.values if condition(x))
    def mapped(self, field):
        values = []
        for record in self.values:
            value = record
            for name in field.split('.'):
                value = getattr(value, name)
            if isinstance(value, Records):
                values.extend(value.values)
            else:
                values.append(value)
        if all(isinstance(x, (Record, NullRecord)) for x in values):
            return Records(x for x in values if x)
        return values
    def __getattr__(self, field):
        if field == 'ids':
            return [x.id for x in self.values]
        result = self.mapped(field)
        if isinstance(result, Records):
            return result
        if not result:
            # Odoo scalar access on an empty recordset is False; numeric empty fields add as zero.
            return False
        if len(result) != 1:
            raise ValueError('recordset scalar requires one record')
        return result[0]
    def _get_valid_products(self, products):
        # Only explicit product-id filters are supported by the profile. Categories/tags/domains are rejected at IPC.
        return {rule: Records(p for p in products if not rule.product_ids or p.id in rule.product_ids)
                for program in self.values for rule in program.rule_ids}


class Currency(Record):
    def __init__(self, code, digits, number=Decimal):
        super().__init__('currency', code, name=code)
        self.digits = digits
        self.number = number
    def _convert(self, amount, other, _company, _date):
        if self.id != other.id:
            raise ValueError('currency conversion is outside this module')
        return amount
    def round(self, amount):
        if self.number is Decimal:
            return amount.quantize(Decimal(10) ** -self.digits, rounding=ROUND_HALF_UP)
        from .oracle_loader import original_float_round
        return original_float_round(amount, precision_digits=self.digits)


class Unit(Record):
    def _compute_quantity(self, quantity, other):
        if self.id != other.id:
            raise ValueError('UoM conversion is outside this module')
        return quantity


class Environment:
    def __getitem__(self, key):
        if key != 'sale.order.line':
            raise ValueError('unsupported ORM surface')
        return Records()


class Order(Record):
    def _get_not_rewarded_order_lines(self):
        return self.order_line.filtered(lambda line: line.product_id and not line.reward_id)
    def _get_no_effect_on_threshold_lines(self):
        return self.order_line.filtered(lambda line: line.threshold_excluded)
    def _get_order_line_price(self, line, price_type):
        return sum(line._get_lines_with_price().mapped(price_type))
    def _allow_nominative_programs(self):
        return False


def build(program_data, order_data, number=Decimal, functions=None):
    from . import engine
    functions = functions or engine
    currency = Currency(program_data['currency'], program_data['currency_digits'], number)
    company = Record('company', order_data['organization_id'])
    program = Record('program', program_data['id'], program_type=program_data['kind'],
                     applies_on=program_data['applies_on'], is_nominative=program_data['nominative'],
                     trigger=program_data['trigger'], currency_id=currency,
                     is_payment_program=program_data['kind'] in ('gift_card', 'ewallet'),
                     trigger_product_ids=program_data['trigger_product_ids'])
    rules = []
    for data in program_data['rules']:
        rule = Record('rule', data['id'], program_id=program,
                      reward_point_mode=data['mode'], reward_point_amount=number(data['points']),
                      reward_point_split=data['split'], minimum_qty=number(data['minimum_qty']),
                      minimum_amount=number(data['minimum_amount']),
                      minimum_amount_tax_mode=data['tax_mode'], product_ids=data['product_ids'],
                      mode='with_code' if data['code_required'] else 'auto')
        rule._compute_amount = lambda target, rule=rule: currency._convert(rule.minimum_amount, target, company, None)
        rules.append(rule)
    program.rule_ids = Records(rules)
    unit = Unit('unit', 'normalized-unit')
    products = {data['product_id']: Record('product', data['product_id'], uom_id=unit) for data in order_data['lines']}
    lines = []
    for data in order_data['lines']:
        reward = NULL
        if data['reward_program_type']:
            reward_program = program if data['reward_program_id'] == program.id else Record(
                'program', data['reward_program_id'], program_type=data['reward_program_type'], trigger=data['reward_trigger'])
            reward = Record('reward', data['id'], reward_type='discount', program_id=reward_program)
        line = Record('line', data['id'], product_id=products[data['product_id']],
                      product_uom_id=unit, product_uom_qty=number(data['quantity']),
                      price_subtotal=number(data['subtotal']), price_tax=number(data['tax']),
                      price_total=number(data['total']), combo_item_id=False,
                      is_reward_line=bool(reward), reward_id=reward,
                      threshold_excluded=data['threshold_excluded'],
                      coupon_id=NULL, points_cost=number('0'))
        line._get_lines_with_price = lambda line=line: Records([line])
        lines.append(line)
    order = Order('order', order_data['order_id'], state=order_data['state'],
                  partner_id=Record('partner', order_data['subject_id'], is_public=order_data['public_subject']),
                  company_id=company, currency_id=currency, amount_total=number(order_data['total']),
                  order_line=Records(lines), coupon_point_ids=Records(),
                  code_enabled_rule_ids=Records(r for r in rules if r.id in order_data['enabled_rule_ids']),
                  env=Environment())
    for name in ('_get_point_changes', '_get_real_points_for_coupon', '_program_check_compute_points'):
        setattr(order, name, MethodType(getattr(functions, name), order))
    return program, order


def coupon_for(program, identity, points, number=Decimal):
    return Record('coupon', identity, program_id=program, currency_id=program.currency_id, points=number(points))
