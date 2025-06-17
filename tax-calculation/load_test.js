import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 200,
  duration: '2m',
};

export default function () {
  const url = 'http://43.229.149.219:8080/tax-calculation';
  const payload = JSON.stringify({
    "income": {
      "monthly_income": 52000,
      "worked_month": 12,
      "bonus": 0,
      "freelance_income": 0
    },
    "deduction": {
      "spouse_deduction": true,
      "children": [
        {
          "birth_year": 2560,
          "studying": true
        }
      ],
      "pregnancy_expenses": 40000,
      "parents": {
        "own_parents": 1,
        "spouse_parents": 1
      },
      "disabled_care": 1,
      "secondary_cities": 15000,
      "shopdee_meekhun": 50000,
      "home_loan_interest": 100000,
      "purchase_otop_products": 5000,
      "purchase_from_community_enterprise": 5000,
      "purchase_from_social_enterprise": 10000,
      "purchase_with_vat_etax": 20000,
      "purchase_with_e_receipt": 10000
    }
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };
  const res = http.post(url, payload, params);
  check(res, {
    'is status 200': (r) => r.status === 200,
  });
    
  sleep(0.01);
}



// export default function () {
//   const res = http.get('http://test.k6.io/');
//   check(res, {
//     'is status 200': (r) => r.status === 200,
//   });
// }
