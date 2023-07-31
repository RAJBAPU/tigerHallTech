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

# Tiger Sighting APIS

1. /tigers: Here the user can make a POST request to store the tiger details.
Request:

{

    "name": "baghh",
    "dob": "2021-07-29T21:25:56+05:30",
    "lastSteenTimeStamp": "2023-07-29T21:25:56+05:30",
    "lastSteenCoordinates": {
        "latitude": 40.7128,
        "longitude": -74.006
    }
}

SuccessResponse:

{

    "code": 200,
    "msg": "success",
    "model": null
}

2. /tigers: Its a GET request to get details of the all tigers with support for pagination.
{
  "code": 200,
  "msg": "success",
  "model": {
    "tigerDetails": [
      {
        "name": "tiger",
        "dob": "1998-07-29T00:00:00Z",
        "lastSteenTimeStamp": "2023-07-29T15:55:56+05:30",
        "lastSteenCoordinates": {
          "latitude": 13.7128,
          "longitude": -23.006
        }
      }
    ],
    "totalTigers": 6
  }
}

3. /tigers/sightings: POST to store tiger sightings. Images are stored like the asked size ie 250*200
Request:
{
    "tigerId": 1,
    "lastSteenTimeStamp": "2023-07-29T21:25:56+05:30",
    "lastSteenCoordinates": {
        "latitude": 45.7128,
        "longitude": -23.006
    },
    "image": ""
}

// to get image, Please upload a jpg file to https://codebeautify.org/jpg-to-base64-converter and get the base64 encoded string, BE is expecting fe to send such a string in my implementaion

Succes Response:
{

    "code": 200,
    "msg": "success",
    "model": null
}

Response of being sighted less than 5Km radius:
{

    "code": 452,
    "msg": "Tiger was already spotted in range",
    "model": null
}

4./tigers/sightings?tigerId= : GET request to get details of a particular tiger. In image fe will get the firebase url of the image but without the base url


Response:
{

    "code": 200,
    "msg": "success",
    "model": {
        "tigerId": 2,
        "sightingDetails": [
            {
                "lastSteenTimeStamp": "2023-07-29T15:55:56+05:30",
                "lastSteenCoordinates": {
                    "latitude": 0,
                    "longitude": 0
                },
                "image": "2%2FtigerImage%2Ff1fcd963-c0e5-049b-10cc-df099a850c7d"
            }
            
        ],
        "total": 6
    }
}

{

1.Kindly update ur firebase credentials in firebase/firebase_json to store and get image
2. Migration script is provided. Please run that to migrate the data base
Pre-Reqs: Beego and Golang has to be installed

}


