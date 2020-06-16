API Server written by Golang for Streetlity Project, aims to high performance in massive processing.

# API
`API Server` providing required function for working with services.

See the [API documentation](https://documenter.getpostman.com/view/4817676/Szezaqyj?fbclid=IwAR05uUEJ7p2dNONCc4kHf-LrD7wpBkLHE6RNp-A_HwdlAGT0e0EK28M3ZYw&version=latest) for more information about the published APIs and some useful Example

## Documentation

Before starting, we need to make some conventions for understanding exactly what we need to do.

- The requests that needs permission need to pass the `jwt-token` via key `Auth` in `Header`
- Parameters for an request will be passed in `query`, `x-www-form-urlencoded` (could add header `"Content-Type": "application/x-www-form-urlencoded"` if you don't know how to do)
### Service
Format: `/service/$serviceName/$method`.

You can take a look for service names that are listed below.

>**POST** /service/$`serviceName`/add  

<pre>
    <b>Header</b>  
    {  
    "timestamp": 12345679121  
    "authen-hash": "3vudiH0Kyo8c7Qa4ihIIvL/yO8fN+ondP6aEhFJlZTA="  
    }  

    <b>Params: JSON</b>  
    {  
        "name": "bank"  
        "id": 564521456  
    }  
</pre>  

# Config
Config are defined in `src/config/config.json`, include:

`database`:
- `server`: address of the host where database is located
- `dbname`: name of database which   
- `username`: username to connect to database server
- `password`: password to connect to database server

