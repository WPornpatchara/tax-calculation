@startuml

title Tax Calculation - POST /core-calculate-tax/api/v1/calculate 

actor Requestor #deepskyblue

box "Core Layer" #LightBlue
  entity "core-tax-service" as core #DeepSkyBlue
end box

box "Database" #LightCyan
  database "tax_db" as db #salmon
end box

Requestor -> core : POST /core-calculate-tax/api/v1/calculate 
note right of Requestor
{
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
}
end note

activate core

group Calculate Total Income
  core -> core : total_income = (monthly_income * worked_month) + bonus + freelance_income
end group

group Calculate Total Deduction
  core -> core : total_deduction = sum all deductions
  note right of core
Includes:
Personal (60,000)
Spouse (60,000)
Child (30,000 + 30,000 for studying)
Pregnancy (40,000)
Parents (2 * 30,000)
Disabled Care (60,000)
Secondary Cities (15,000)
ShopDeeMeekhun (50,000)
Home Loan (100,000)
OTOP (5,000)
Community (5,000)
Social Enterprise (10,000)
VAT e-Tax (20,000)
e-Receipt (10,000)
end note
end group

group Tax Calculation
  core -> core : net_income = total_income - total_deduction
  alt net_income > 0
    core -> core : total_tax = apply tax brackets
  else
    core -> core : total_tax = 0
  end alt
end group

group Calculate Refund
  core -> core : refund = tax_paid - total_tax
  note right of core
If tax_paid > total_tax --> refund  
If tax_paid < total_tax --> need to pay more  
If equal --> no refund
end note
end group

group Save to DB
  core -> db : INSERT INTO tax_result_db(total_income, total_deduction, total_tax, refund)
  activate db
  db --> core : 
  deactivate db
end group

opt database returns error
  core --> Requestor : 
  note right of Requestor 
  {
      "code": 5000,
      "message": "Generic error"
  }
  end note
end opt

core --> Requestor : Response
note right of Requestor
{
  "code": 1000,
  "message": "success",
  "data": {
    "total_income": 624000,
    "total_deduction": 525000,
    "total_tax": 0,
    "refund": 1000
  }
}
end note
deactivate core
hide footbox
@enduml