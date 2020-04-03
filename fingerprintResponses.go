package lzr

import (
    "strings"
)


func fingerprintResponse( data string ) string {

    if strings.Contains( data, "HTTP" ){
         return "HTTP"
    } else if strings.Contains( data, "maxBsonObjectSize" ) {
        return "Mongo"
    } else {
        return ""
    }

}
