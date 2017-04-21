# README #

### What is this repository for? ###

* AccountCreation module of CMPUTE: Automates the process of AccountCreation .
* Version v1.0

###How to run?####
*The module accepts three parameters, AccessKey,SecretKey and a CSV file which has your credentials.
*The module loads the data from the CSV file  and generates the outputfile with name "results.csv"
*The input CSV file has to be maintained in a format that is given below.

		ACCESS KEY,             SECRECT KEY,                               REGION,        ACCOUNT NAME,         DOMAIN 
		xxxxxxxxxxxxxxxxxxxx,   ****************************************,  us-east-1,     Dev,                https://<tenantname>.cmpute.io
		xxxxxxxxxxxxxxxxxxxx,   ****************************************,  us-east-1,     Staging,            https://<tenantname>.cmpute.io



