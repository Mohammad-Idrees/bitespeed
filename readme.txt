
email/phone

Contact
{
    id                   Int                   
    phoneNumber          String?
    email                String?
    linkedId             Int? // the ID of another Contact linked to this one
    linkPrecedence       "secondary"|"primary" // "primary" if it's the first Contact in the link
    createdAt            DateTime              
    updatedAt            DateTime              
    deletedAt            DateTime?
}
Note: Question mark represents optional fields


eg1

If a customer placed an order with 
email=lorraine@hillvalley.edu & phoneNumber=123456 
and later came back to place another order with 
email=mcfly@hillvalley.edu & phoneNumber=123456 ,
database will have the following rows:

Note: Both contacts are linked together because they share the same phone number.

{
    id                   1                   
    phoneNumber          "123456"
    email                "lorraine@hillvalley.edu"
    linkedId             null
    linkPrecedence       "primary"
    createdAt            2023-04-01 00:00:00.374+00              
    updatedAt            2023-04-01 00:00:00.374+00              
    deletedAt            null
},
{
    id                   23                   
    phoneNumber          "123456"
    email                "mcfly@hillvalley.edu"
    linkedId             1
    linkPrecedence       "secondary"
    createdAt            2023-04-20 05:30:00.11+00              
    updatedAt            2023-04-20 05:30:00.11+00              
    deletedAt            null
}



API
/identify
{
	"email"?: string,
	"phoneNumber"?: number
}


{
    "contact":{
        "primaryContatctId": number,
        "emails": string[], // first element being email of primary contact 
        "phoneNumbers": string[], // first element being phoneNumber of primary contact
        "secondaryContactIds": number[] // Array of all Contact IDs that are "secondary" to the primary contact
    }
}



1. check if email/phone exists
does not exists
    1. insert
exists
    1. only phone exists
        1. get linkedId
        2. insert secondary and set linkedId

        
    
