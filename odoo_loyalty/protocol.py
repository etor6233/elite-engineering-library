# AUTHORED bounded JSON/IPC and representation glue. Source-derived rules live in engine.py.
import hashlib
import json
import re
import sys
from decimal import Decimal, localcontext
from pathlib import Path

from . import engine
from .orm_contract import Record, Records, build, coupon_for

REVISION = '99edb6dd82b7b560930c00b03b694ba700785370'
MAX_BYTES = 131072
DECIMAL = re.compile(r'-?(?:0|[1-9][0-9]{0,17})(?:\.[0-9]{1,6})?\Z')
IDENTITY = re.compile(r'[A-Za-z0-9][A-Za-z0-9_.:-]{0,127}\Z')


def exact(value, fields):
    if type(value) is not dict or set(value) != set(fields.split()):
        raise ValueError('unexpected or missing fields')


def identity(value):
    if type(value) is not str or not IDENTITY.fullmatch(value):
        raise ValueError('invalid identity')
    return value


def decimal(value, minimum=None, positive=False):
    if type(value) is not str or not DECIMAL.fullmatch(value):
        raise ValueError('bounded canonical decimal string required')
    number = Decimal(value)
    if number == 0 and value.startswith('-'):
        raise ValueError('negative zero is not canonical')
    if minimum is not None and number < minimum or positive and number <= 0:
        raise ValueError('decimal outside permitted range')
    return number


def flag(value):
    if type(value) is not bool:
        raise ValueError('boolean required')


def identities(values, maximum=100):
    if type(values) is not list or len(values) > maximum:
        raise ValueError('identity list exceeds limit')
    for value in values:
        identity(value)
    if len(set(values)) != len(values):
        raise ValueError('duplicate identity')


def canonical(value):
    return json.dumps(value, sort_keys=True, separators=(',', ':'), ensure_ascii=False).encode('utf-8')


def profile(program):
    exact(program, 'id kind currency currency_digits applies_on nominative trigger trigger_product_ids rules')
    identity(program['id'])
    if program['kind'] not in ('gift_card', 'loyalty') or program['applies_on'] not in ('future', 'current', 'both'):
        raise ValueError('unsupported program kind/application')
    if not re.fullmatch('[A-Z]{3}', program['currency']) or type(program['currency_digits']) is not int or program['currency_digits'] not in range(5):
        raise ValueError('invalid currency contract')
    if program['trigger'] not in ('auto', 'with_code'):
        raise ValueError('unsupported trigger')
    flag(program['nominative'])
    identities(program['trigger_product_ids'])
    rules = program['rules']
    if type(rules) is not list or not 1 <= len(rules) <= 20:
        raise ValueError('bounded nonempty rules required')
    ids = []
    for rule in rules:
        exact(rule, 'id mode points split minimum_qty minimum_amount tax_mode product_ids code_required')
        ids.append(identity(rule['id']))
        if rule['mode'] not in ('order', 'money', 'unit') or rule['tax_mode'] not in ('incl', 'excl'):
            raise ValueError('unsupported rule mode')
        decimal(rule['points'], positive=True)
        decimal(rule['minimum_qty'], minimum=0)
        decimal(rule['minimum_amount'], minimum=0)
        flag(rule['split'])
        flag(rule['code_required'])
        identities(rule['product_ids'])
        # Directly reflects LoyaltyRule._constraint_trigger_multi; not a new policy.
        if rule['split'] and program['applies_on'] == 'both':
            raise ValueError('split per unit is not allowed when applies_on is both')
    if len(ids) != len(set(ids)):
        raise ValueError('duplicate rule')


def order_contract(order, program):
    exact(order, 'order_id organization_id subject_id state public_subject currency total enabled_rule_ids lines')
    for key in ('order_id', 'organization_id', 'subject_id'):
        identity(order[key])
    if order['state'] not in ('draft', 'sent', 'sale', 'done') or order['currency'] != program['currency']:
        raise ValueError('state or currency binding differs')
    flag(order['public_subject'])
    identities(order['enabled_rule_ids'], 20)
    if not set(order['enabled_rule_ids']) <= {r['id'] for r in program['rules']}:
        raise ValueError('unknown enabled rule')
    decimal(order['total'], minimum=0)
    if type(order['lines']) is not list or not 1 <= len(order['lines']) <= 100:
        raise ValueError('bounded order lines required')
    ids = []
    for line in order['lines']:
        exact(line, 'id product_id quantity subtotal tax total reward_program_type reward_program_id reward_trigger threshold_excluded')
        ids.append(identity(line['id']))
        identity(line['product_id'])
        quantity = decimal(line['quantity'], minimum=0)
        if quantity != quantity.to_integral_value() or quantity > 10000:
            raise ValueError('only bounded normalized integer quantities admitted')
        subtotal, tax, total = (decimal(line[k]) for k in ('subtotal', 'tax', 'total'))
        if subtotal + tax != total:
            raise ValueError('line amounts do not bind')
        flag(line['threshold_excluded'])
        if line['reward_program_type'] not in ('', 'loyalty', 'gift_card', 'ewallet', 'promotion'):
            raise ValueError('unknown reward program')
        if line['reward_program_type']:
            identity(line['reward_program_id'])
            if line['reward_trigger'] not in ('auto', 'with_code'):
                raise ValueError('invalid reward trigger')
        elif line['reward_program_id'] != '' or line['reward_trigger'] != '' or total < 0:
            raise ValueError('normal line reward binding differs')
    if len(set(ids)) != len(ids) or sum(decimal(line['total']) for line in order['lines']) != decimal(order['total']):
        raise ValueError('duplicate line or order total mismatch')


def stringify(value):
    if isinstance(value, Decimal):
        if not value.is_finite():
            raise ValueError('non-finite result')
        return format(value, 'f')
    if type(value) is dict:
        return {str(k): (v if k == 'applied_minor_units' and type(v) is int else stringify(v)) for k,v in value.items()}
    if type(value) is list:
        return [stringify(v) for v in value]
    if type(value) is int and type(value) is not bool:
        return str(value)
    return value


def calculate(request):
    exact(request, 'schema operation program program_sha256 order data')
    if request['schema'] != 'elite.odoo-loyalty-calc.v1' or request['operation'] not in ('evaluate', 'reward', 'changes', 'reverse'):
        raise ValueError('unsupported contract')
    profile(request['program'])
    expected = hashlib.sha256(canonical(request['program'])).hexdigest()
    if request['program_sha256'] != expected:
        raise ValueError('program hash mismatch')
    with localcontext() as context:
        context.prec = 80
        order_contract(request['order'], request['program'])
        program, order = build(request['program'], request['order'])
        data = request['data']
        if request['operation'] == 'evaluate':
            exact(data, '')
            result = order._program_check_compute_points(Records([program]))[program]
            if 'error' not in result:
                result['points'] = engine.filtered_issuance(result, program)
        elif request['operation'] == 'reward':
            exact(data, 'coupon_id balance pending_earned pending_cost discountable reward')
            identity(data['coupon_id'])
            for key in ('balance', 'pending_earned', 'pending_cost'):
                decimal(data[key])
            discountable = decimal(data['discountable'], minimum=0)
            if discountable > order.amount_total:
                raise ValueError('discountable exceeds bound order total')
            reward_data = data['reward']
            exact(reward_data, 'mode discount required_points max_amount clear_wallet')
            if reward_data['mode'] not in ('per_point', 'per_order', 'percent'):
                raise ValueError('unsupported reward mode')
            decimal(reward_data['discount'], positive=True)
            decimal(reward_data['required_points'], positive=True)
            decimal(reward_data['max_amount'], minimum=0)
            flag(reward_data['clear_wallet'])
            coupon = coupon_for(program, data['coupon_id'], data['balance'])
            order.coupon_point_ids = Records([Record('coupon_point', 'pending', coupon_id=coupon, points=Decimal(data['pending_earned']))])
            cost = Record('cost', 'pending', coupon_id=coupon, points_cost=Decimal(data['pending_cost']), reward_id=True)
            order.order_line = order.order_line | Records([cost])
            reward = Record('reward', 'requested', program_id=program, currency_id=program.currency_id,
                            discount_mode=reward_data['mode'], discount=Decimal(reward_data['discount']),
                            required_points=Decimal(reward_data['required_points']), discount_max_amount=Decimal(reward_data['max_amount']),
                            clear_wallet=reward_data['clear_wallet'])
            # Same eligibility precondition as _get_claimable_rewards before invoking discount calculation.
            if order._get_real_points_for_coupon(coupon) < reward.required_points or discountable == 0:
                raise ValueError('reward is not claimable')
            result = engine.reward_from_bound_amounts(order, reward, coupon, discountable)
            result['value'] = program.currency_id.round(result['value'])
        elif request['operation'] == 'reverse':
            exact(data, 'entries')
            if type(data['entries']) is not list or not 1 <= len(data['entries']) <= 40:
                raise ValueError('bounded reversal entries required')
            result = []
            ids = set()
            for entry in data['entries']:
                exact(entry, 'account_id points_delta applied_minor_units')
                identity(entry['account_id'])
                if entry['account_id'] in ids:
                    raise ValueError('duplicate account reversal')
                ids.add(entry['account_id'])
                value = decimal(entry['points_delta'])
                minor = entry['applied_minor_units']
                if type(minor) is not int or not -(2**63)+1 <= minor <= (2**63)-1:
                    raise ValueError('bounded exact minor units required')
                result.append({'account_id':entry['account_id'],
                               'points_delta':engine.reverse_point_change(Decimal('0'), value),
                               'applied_minor_units':-minor})
        else:
            exact(data, 'issued used')
            coupons = {}
            for key in ('issued', 'used'):
                if type(data[key]) is not list or len(data[key]) > 100:
                    raise ValueError('changes exceed limit')
                for item in data[key]:
                    exact(item, 'coupon_id points')
                    identity(item['coupon_id'])
                    decimal(item['points'], minimum=0)
                    coupons.setdefault(item['coupon_id'], coupon_for(program, item['coupon_id'], '0'))
            order.coupon_point_ids = Records(Record('earned', str(i), coupon_id=coupons[x['coupon_id']], points=Decimal(x['points'])) for i,x in enumerate(data['issued']))
            order.order_line = Records(Record('spent', str(i), coupon_id=coupons[x['coupon_id']], points_cost=Decimal(x['points']), reward_id=True) for i,x in enumerate(data['used']))
            result = {coupon.id: amount for coupon,amount in order._get_point_changes().items()}
        return {'schema':'elite.odoo-loyalty-result.v1','source_revision':REVISION,
                'program_sha256':expected,'request_sha256':hashlib.sha256(canonical(request)).hexdigest(),
                'result':stringify(result)}


def pairs(entries):
    result = {}
    for key,value in entries:
        if key in result or key.casefold() in (x.casefold() for x in result):
            raise ValueError('duplicate key')
        result[key] = value
    return result


def main():
    try:
        raw = sys.stdin.buffer.read(MAX_BYTES+1)
        if len(raw) > MAX_BYTES:
            raise ValueError('request exceeds limit')
        request = json.loads(raw.decode('utf-8'), object_pairs_hook=pairs,
                             parse_constant=lambda _value: (_ for _ in ()).throw(ValueError('nonfinite JSON')))
        output = canonical(calculate(request))
        if len(output) > MAX_BYTES:
            raise ValueError('response exceeds limit')
        sys.stdout.buffer.write(output+b'\n')
        return 0
    except (ValueError, KeyError, TypeError, ArithmeticError):
        # No request data or upstream internals appear in errors.
        sys.stdout.buffer.write(b'{"error":"invalid_or_unsupported_calculation"}\n')
        return 2


if __name__ == '__main__':
    raise SystemExit(main())
