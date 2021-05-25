/* quenubesHTTP simplifies usage of the XMLHttpRequest object. */
class quenubesHTTP {
    constructor(url,method, headers,data, callbckFunc) {
      this.url = url;
      this.method = method.toUpperCase();
      this.headers = headers;
      this.data = data;
      this.responseCallback = callbckFunc;
  
      /* Initialize the Response */
      this.res = {
        statusCode : 0,
        responseText : "",
        errMessage : "",
        timedOut :false  
      }        
    }
    call() {
        let resx = this.res;
        let h = this.headers;
        let rx = this.responseCallback;
        let xhttp = new XMLHttpRequest();
        xhttp.open(this.method, this.url, true);
        /* Add the headers */
        if(h != null){
          for(let i=0; i < h.length; i++){
            let v = h[i].split(":");
            xhttp.setRequestHeader(v[0].trim(),v[1].trim());
          }
        }
      if(this.method == 'POST' && this.data != null){
        let j = JSON.stringify(this.data);
        xhttp.send(j); /* Post data */
      }else{
        xhttp.send();  /* Get data */
      }      
      /* All of these use the same callback so that the caller does not have to 
       setup a separate function for each http result. */
      xhttp.ontimeout = () => {
        resx.timedOut = true;
        rx(resx);
      };
      xhttp.onerror = function () {
        /* The error is not always in the expected format. */
        try {
          let err = JSON.parse(xhttp.responseText);
          resx.errMessage = err.message;
        } catch {
          resx.errMessage = xhttp.responseText;
        }
        rx(resx);
      };
        xhttp.onreadystatechange = function () {
          if (this.readyState == 4){             
            resx.responseText = xhttp.responseText;  
            resx.statusCode = xhttp.status;  
            rx(resx);        
          }   
        }    
    } /* End of call function */
  }