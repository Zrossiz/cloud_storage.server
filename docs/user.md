# User
## Routes
* Create user  
POST /api/user/registration  
Accept json object:  
{  
    &nbsp;&nbsp;"name": "name_user",  
    &nbsp;&nbsp;"email": "email_user",  
    &nbsp;&nbsp;"password": "password_user"  
}
----
* Login user  
POST /api/user/login  
Accept json object:  
{  
    &nbsp;&nbsp;"name": "name_user",  
    &nbsp;&nbsp;"password": "password_user"  
}  
## API response
Example success registration/login 
 
Cookies:  
{  
&nbsp;&nbsp;userId: number HttpOnly  
&nbsp;&nbsp;session: string HttpOnly  
}  

JSON response:  
{  
    &nbsp;&nbsp;"success": boolean,  
    &nbsp;&nbsp;"data": {  
        &nbsp;&nbsp;&nbsp;&nbsp;"id": number,  
        &nbsp;&nbsp;&nbsp;&nbsp;"name": string,  
        &nbsp;&nbsp;&nbsp;&nbsp;"email": string,  
        &nbsp;&nbsp;&nbsp;&nbsp;"password": string  
        &nbsp;&nbsp;&nbsp;&nbsp;"created_at": Date,  
        &nbsp;&nbsp;&nbsp;&nbsp;"updated_at": Date,  
        &nbsp;&nbsp;&nbsp;&nbsp;"deleted_at": Date,  
    &nbsp;&nbsp;}  
}