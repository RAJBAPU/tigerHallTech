# For SignIn And SignUp

1. /signup: The user provides the required credentials and makes a POST request to the endpoint to create a new account. 
The server receives the request, validates the credentials, adds the user to the database, and sends a verification code to the user’s email address.

Request: 
{
    "name":"Ayush Raj",
    "email": "ayursaj@gmail.com",
    "password": "miqcnqke",
    "passwordConfirm": "miqcnqke"

}

Success Response:
{
    "code": 200,
    "msg": "success",
    "model": null
}

2. /verifyemail/:verificationCode :
In the email, user recieves a verification code.
However, since there is no frontend application, you need to copy the verification code from the redirect URL and manually make the request to the server.
Paste the verification code in the URL and make a GET request to the /verifyemail/:verificationCode endpoint to verify the email address.
Success Response:
{
    "code": 200,
    "msg": "success",
    "model": null
}

3./login: Once the account has been verified, the user can provide the email and password used in registering for the account and make a POST request to the /login endpoint to sign into the account.
Request:
{ 
    "email": "ayursaj@gmail.com",
    "password": "miqcnqke"
}

Success Response:
{
    "code": 200,
    "msg": "success",
    "model": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTQ0MjE4MzUsImlhdCI6MTY5MDgyMTgzNSwibmJmIjoxNjkwODIxODM1LCJzdWIiOjIwfQ.vmwj1ttEsyD2Nl3xuA7hk-uGReZmY4RO7jMgv2eHWis"
}

4./me : Here, the user can make a GET request to the /me endpoint with the token received from the server to retrieve his credentials.

{
    "code": 200,
    "msg": "success",
    "model": {
        "id": 20,
        "name": "Ayush Raj",
        "email": "ayursaj@gmail.com"
    }
}

5. 



