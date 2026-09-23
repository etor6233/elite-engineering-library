import copy
import hashlib
import json
import random
import subprocess
import sys
import unittest
from decimal import Decimal
from pathlib import Path

from odoo_loyalty import engine
from odoo_loyalty.orm_contract import Records, Record, build, coupon_for
from odoo_loyalty.oracle_loader import original_functions
from odoo_loyalty.protocol import calculate, canonical, pairs


def program(kind='gift_card'):
    # ADAPTED fixture configuration: Odoo tests/common.py:104-130, source-only arithmetic.
    return {'id':'gift-1' if kind == 'gift_card' else 'loyalty-1','kind':kind,'currency':'ARS',
            'currency_digits':2,'applies_on':'future' if kind == 'gift_card' else 'both',
            'nominative':kind == 'loyalty','trigger':'auto','trigger_product_ids':['gift-50'],
            'rules':[{'id':'rule-1','mode':'money','points':'1','split':kind == 'gift_card',
                      'minimum_qty':'1','minimum_amount':'0','tax_mode':'incl',
                      'product_ids':['gift-50'] if kind == 'gift_card' else [],'code_required':False}]}


def order(quantity=2, total='100', product='gift-50'):
    return {'order_id':'order-1','organization_id':'org-1','subject_id':'subject-1','state':'draft',
            'public_subject':False,'currency':'ARS','total':total,'enabled_rule_ids':[],
            'lines':[{'id':'line-1','product_id':product,'quantity':str(quantity),'subtotal':total,'tax':'0','total':total,
                      'reward_program_type':'','reward_program_id':'','reward_trigger':'','threshold_excluded':False}]}


def request(p=None, o=None, operation='evaluate', data=None):
    p = p or program()
    return {'schema':'elite.odoo-loyalty-calc.v1','operation':operation,'program':p,
            'program_sha256':hashlib.sha256(canonical(p)).hexdigest(),'order':o or order(),'data':data or {}}


class DerivedTests(unittest.TestCase):
    def test_original_functions_are_source_ast(self):
        manifest=json.loads((Path(__file__).parent/'odoo_loyalty/DERIVATION.json').read_text())
        for mapping in manifest['functions']:
            path=Path(__file__).parent/'odoo_loyalty/upstream'/mapping['source_path']
            raw=path.read_bytes()
            self.assertEqual(hashlib.sha256(raw).hexdigest(),mapping['source_file_sha256'])
            lo,hi=mapping['source_lines']
            self.assertEqual(hashlib.sha256(b''.join(raw.splitlines(keepends=True)[lo-1:hi])).hexdigest(),mapping['source_segment_sha256'])

    def test_upstream_gift_card_quantity_expectations(self):
        # Actual expected card counts from test_buy_gift_card.py:12-36; full Odoo ORM suite is not claimed.
        for quantity in [1,2,1]:
            result=calculate(request(o=order(quantity,str(quantity*50))))['result']
            cards=[Decimal(p) for p in result['points'] if Decimal(p)]
            self.assertEqual(cards,[Decimal(50)]*quantity)

    def test_rule_modes_and_source_oracle(self):
        original=original_functions()
        rng=random.Random(40219)
        for mode in ['order','money','unit']:
            for _ in range(80):
                quantity=rng.randint(1,20)
                total=Decimal(rng.randint(1,100000))/100
                p=program('loyalty')
                p['rules'][0].update(mode=mode,points=str(Decimal(rng.randint(1,100))/10),minimum_amount=str(rng.randint(0,10)))
                o=order(quantity,format(total,'f'),'product-A')
                exact_p,exact_o=build(p,o)
                float_p,float_o=build(p,o,float,original)
                actual=exact_o._program_check_compute_points(Records([exact_p]))[exact_p]
                oracle=float_o._program_check_compute_points(Records([float_p]))[float_p]
                self.assertEqual(actual.keys(),oracle.keys())
                if 'points' in actual:
                    self.assertEqual([x.quantize(Decimal('.01')) if isinstance(x,Decimal) else Decimal(x) for x in actual['points']],
                                     [Decimal(str(x)).quantize(Decimal('.01')) for x in oracle['points']])
                else:
                    self.assertEqual(actual,oracle)

    def test_source_filters_discount_and_payment_lines(self):
        p=program('loyalty')
        o=order(1,'100','product-A')
        for kind,amount in [('promotion','-10'),('gift_card','-20')]:
            o['lines'].append({'id':kind,'product_id':'product-A','quantity':'1','subtotal':amount,'tax':'0','total':amount,
                               'reward_program_type':kind,'reward_program_id':kind,'reward_trigger':'with_code','threshold_excluded':False})
        o['total']='70'
        # Source excludes gift-card tender from money-spend rewards, but includes the discount line.
        self.assertEqual(calculate(request(p,o))['result']['points'],['90.00'])
        o['lines'][1].update(subtotal='-0.25',total='-0.25')
        o['total']='79.75'
        self.assertEqual(calculate(request(p,o))['result']['points'],['99.75'])

    def test_threshold_code_products_and_public_subject(self):
        p=program()
        p['rules'][0].update(code_required=True,minimum_amount='100',minimum_qty='2')
        r=request(p,order())
        self.assertIn('error',calculate(r)['result'])
        r['order']['enabled_rule_ids']=['rule-1']
        self.assertNotIn('error',calculate(r)['result'])
        for amount,quantity in [('99.99',2),('100',1)]:
            r['order']=order(quantity,amount)
            r['order']['enabled_rule_ids']=['rule-1']
            self.assertIn('error',calculate(r)['result'])
        r=request(program(),order(product='other-product'))
        self.assertIn('error',calculate(r)['result'])
        r=request(program('loyalty'))
        r['order']['public_subject']=True
        self.assertIn('error',calculate(r)['result'])

    def test_available_points_future_and_finalized(self):
        original=original_functions()
        for kind in ['gift_card','loyalty']:
            for state,expected in [('draft',Decimal(80) if kind=='gift_card' else Decimal(90)),('sale',Decimal(100))]:
                for number,functions in [(Decimal,None),(float,original)]:
                    p,o=build(program(kind),order(),number,functions)
                    o.state=state
                    card=coupon_for(p,'card-1','100',number)
                    o.coupon_point_ids=Records([Record('earned','e',coupon_id=card,points=number('10'))])
                    o.order_line=Records([Record('used','u',coupon_id=card,points_cost=number('20'))])
                    self.assertEqual(Decimal(str(o._get_real_points_for_coupon(card))),expected)

    def test_gift_card_over_under_and_multiple_bound_amounts(self):
        # Source test_pay_with_gift_card.py:12-69 expects total reduction == before/after balance delta.
        for order_amount,balance,expected in [('115','100','100'),('5.75','100','5.75'),('2300','200','200')]:
            data={'coupon_id':'card-1','balance':balance,'pending_earned':'0','pending_cost':'0','discountable':order_amount,
                  'reward':{'mode':'per_point','discount':'1','required_points':'1','max_amount':'0','clear_wallet':False}}
            result=calculate(request(o=order(1,order_amount),operation='reward',data=data))['result']
            self.assertEqual(Decimal(result['value']),Decimal(expected))
            self.assertEqual(Decimal(result['points_cost']),Decimal(expected))
            self.assertEqual(Decimal(balance)-Decimal(result['points_cost']),Decimal(balance)-Decimal(expected))

    def test_changes_preserves_card_identity_and_reversal_amount(self):
        data={'issued':[{'coupon_id':'card-A','points':'100'},{'coupon_id':'card-A','points':'5'},{'coupon_id':'card-B','points':'50'}],
              'used':[{'coupon_id':'card-A','points':'30'},{'coupon_id':'card-B','points':'10'}]}
        result=calculate(request(operation='changes',data=data))['result']
        self.assertEqual(result,{'card-A':'75','card-B':'40'})
        # The durable layer must record the negation of this same committed change for source cancel intent.
        self.assertEqual({k:-Decimal(v) for k,v in result.items()},{'card-A':Decimal(-75),'card-B':Decimal(-40)})

    def test_exact_decimal_large_value_does_not_inherit_float_loss(self):
        p=program('loyalty')
        o=order(1,'99999999999999999.99','product-A')
        self.assertEqual(calculate(request(p,o))['result']['points'],['99999999999999999.99'])
        self.assertNotEqual(Decimal(str(float('99999999999999999.99'))),Decimal('99999999999999999.99'))

    def test_invalid_hash_amount_binding_and_no_partial_contract(self):
        mutations=[lambda r:r.update(program_sha256='0'*64),
                   lambda r:r['program']['rules'][0].update(points='0'),
                   lambda r:r['program'].update(kind='ewallet'),
                   lambda r:r['order'].update(currency='USD'),
                   lambda r:r['order'].update(total='101'),
                   lambda r:r['order']['lines'][0].update(quantity='1.5'),
                   lambda r:r['order']['lines'][0].update(tax='1'),
                   lambda r:r['order'].update(public_subject=None),
                   lambda r:r['program'].update(currency_digits=True),
                   lambda r:r['order']['lines'][0].update(subtotal=100),
                   lambda r:r['order'].update(enabled_rule_ids=['unknown']),
                   lambda r:r['program'].update(tax_rate='21'),
                   lambda r:r['order'].update(lines=r['order']['lines']*2)]
        for mutate in mutations:
            r=request();mutate(r)
            with self.subTest(r=r),self.assertRaises((ValueError,TypeError)):
                calculate(r)

    def test_duplicate_json_and_exact_ipc_output(self):
        with self.assertRaises(ValueError):
            json.loads('{"schema":1,"Schema":2}',object_pairs_hook=pairs)
        r=request()
        process=subprocess.run([sys.executable,'-B','-m','odoo_loyalty.protocol'],input=canonical(r),capture_output=True,cwd=Path(__file__).parent,timeout=5)
        self.assertEqual(process.returncode,0,process.stderr)
        self.assertEqual(json.loads(process.stdout),calculate(r))
        bad=subprocess.run([sys.executable,'-B','-m','odoo_loyalty.protocol'],input=b'{"sensitive":"never-echo-this"}',capture_output=True,cwd=Path(__file__).parent,timeout=5)
        self.assertEqual(bad.returncode,2)
        self.assertNotIn(b'never-echo-this',bad.stdout+bad.stderr)


if __name__ == '__main__':
    unittest.main()
