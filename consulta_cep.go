package main 


import (
    
    "fmt"
    "net/http"
    "io/ioutil"
    "strings"
    "time"
)

func check(err error) {
    if err != nil {
        fmt.Println(err)
    }
}

func quant_goroutines(n int, t <-chan string, r chan<- []byte) {
   
    for c := 1; c <= n; c++ {
        
        go consulta(n, t, r)
    }
    
}

func coloca_url(ceps []string, canal chan<- string) {
    for _, cep := range ceps {
        
        url := fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cep)
        
        canal <- url
    }
    
    close(canal)
}

func consulta(n int, trabalho <-chan string, resultado chan<- []byte) {
    
    for url := range trabalho {
        
        req, err:= http.Get(url)
        check(err)
    
        body, err := ioutil.ReadAll(req.Body)
        check(err)
        
    resultado <- body
    }  
}

func tira_trabalho(ceps []string, canal <-chan []byte) {
    
    for _, cep := range ceps {
        
        err := ioutil.WriteFile("./" + cep + ".json", <-canal, 0644)
        check(err)
    }
}

func main() {
    
    now := time.Now()
    
    trabalho := make(chan string, 200)
    resultado := make(chan []byte, 200)
    
    leitor, err := ioutil.ReadFile("cep.txt")
    check(err)
    
    ceps := strings.Split(string(leitor), "\n")
    
    quant_goroutines(200, trabalho, resultado)
    
    coloca_url(ceps, trabalho)
    
    tira_trabalho(ceps, resultado)
    
    fmt.Println("tempo de execução:", time.Since(now).Seconds(),"s")
}