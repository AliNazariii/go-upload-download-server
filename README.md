# Concurrent Upload Download HTTP Server

# Introduction

You are to use net/http package and concurrency in golang to implement a simple http server.
First, You are to write a program that reading from the specific file and writing into
another file concurrently.
Then, You should implement an api to get file or send file to the user.
REST APIs that you should implement are listed below:

## `localhost:8080/uploadFile`
    
***input formats*** :
    
1. **json** :
      ```json
       {
           "file" : "string"
       }
      ```

      In this format, `file` is an url that you should get file from.


2. **form** :

         file : []byte

      In this format, `file` is a byte array of the actual file.


<br />

***output format*** :

1. **json** :
   
      *successful upload* : 
      ```json
       {
           "file_id" : "string"
       }
      ```

      `file_id` is a unique identifier of uploaded file.
      It has 2 parts  : <br/>
      ***access_hash : file_name*** <br/>
      access_hash is an 64 bit hashed number, and the file_name is an 
      encrypted value of the actual file.\
      \
      \
      *failure upload* :
      ```json
       {
           "error" : "string"
       }
      ```
      
      `error` is a description of the occurred error.



<br />

## `localhost:8080/downloadFile` 

1. **json** :
      ```json
       {
           "file_id" : "string"
       }
      ```

   In this format, `file_id` is an id that we got in successful upload request.


2. **form** :

         file_id : string

   In this format, `file_id` is an id that we got in successful upload request.


<br />

***output format*** :

1. **json** :

   *successful download* :
      
         http response with actual file
   \
   *failure download* :
      ```json
       {
           "error" : "string"
       }
      ```

   `error` is a description of the occurred error.
   
      




# Roadmap

- [x] Implement reading concurrently from a file.
- [x] Implement writing concurrently to a file.
- [x] Examine the writing performance with different number of goroutines and number of bytes each goroutine writes in one access.
- [x] Add REST APIs.
- [x] Complete your upload/download http server.
