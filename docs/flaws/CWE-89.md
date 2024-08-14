# CWE-89: Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')

VerademoGO uses a call DB::QueryRow(), which constructs a dynamic SQL query using a variable that is derived from unverified input. Attackers can exploit this by executing arbitrary SQL statements.

## Mitigate
* Validate all input and ensure that it matches the expected format. 

## Remediate 
* Utilize parameterized statements rather than SQL queries.

# Resources 
* [CWE-89](https://cwe.mitre.org/data/definitions/89.html)
* [OWASP] (https://owasp.org/www-community/attacks/SQL_Injection)