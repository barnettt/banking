# Udemy Golang Training Banking App

Following the Udemy course to understand how to use golang.

A simple banking app demonstrating the use of golang constructs.

This was learning exercise, I have added a bit of my own thought n certain elements, such as the transcation
manager. I felt it better to separate this concern.
The application shows how to create an application following the hexagonal development pattern, where domain
objects and components expose an interface contract enabling interaction.
Banking consists of 2 apps: Banking and Banking-auth.
#Banking
Banking covers the simple use case of crud operations for a customer i.e. deposit and withdraw from an account.
An admin can create a customer and an account for a customer.
#Baning-auth
Is the security component which uses the bearer token in the request to validate a customer.
The customer at login time, is returned a token that can be used in subsequent requests.
A claim for the customer is created ensuring that the only access their own accounts.

The Banking app, on receipt of a request will call out to the Banking-auth app to verify the cutomers 
rights and permissions before allowing any operations for the customer.
#Configuration
Requires a mysql database see the db/banking.sql.
###### Parameters
The following configuration parameters are set on the command line:
```
SERVER_HOST=localhost;
SERVER_PORT=8000;
DB_HOST=localhost;
DB_PORT=3306;
DB_USER=banking;
DB_PASSWD=banking;
DB_PROTOCOL=tcp;
DB_NAME=banking;
DB_DRIVER_NAME=mysql;
AUTH_SERVER_PORT=8001
```
