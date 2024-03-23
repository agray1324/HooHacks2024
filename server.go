package main



func Server(){
    r := router.Router()    
    r.Run(":8080")  
    
}