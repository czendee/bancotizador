package  util
import (
	"log"
)


//     month,year, errorValid := convertMMYYintoMonthYear(requestData.Valid) 
     func ConvertMMYYintoMonthYear (validall string) (string,string, string){
	      var month string
	      var errorGeneral string
	      var year string
	      	      
	      month = ""
	      errorGeneral=""
	      
	     if len(validall) > 2 {
	          month= validall[:2]
	     }else{
	     	errorGeneral ="Valid Trhu length is wrong, not including the month "
	     	log.Print("Valid Trhu length is wrong, not including the month ")
	     }
	      
	     if len(validall) > 3 {
	          year=  "20"+validall[len(validall)-2:]
	     }else{
	     	errorGeneral ="Valid Trhu length is wrong, not including the year "
	     }
         return month, year, errorGeneral     	
     }
      
//     dateDDMMYY, errorValid := convertMMYYintoDDMMYY(requestData.Valid)       
     func ConvertMMYYintoDDMMYY (validall string) (string, string){
	      var month string
	      var errorGeneral string
	      var year string
	      var dateValid string
	      	      
	      month = ""
	      errorGeneral=""
	      dateValid =""
	      year = ""
	      
	     if len(validall) > 2 {
	          month= validall[:2]
	     }else{
	     	errorGeneral ="Valid Trhu length is wrong, not including the month "
	     }
	      
	     if len(validall) > 3 {
	          year=  "20"+validall[len(validall)-2:]  //2000
	          
	          dateValid ="01/"+month+"/"+year     // DD/MM/YYYY
	          
	     }else{
	     	errorGeneral ="Valid Trhu length is wrong, not including the year "
	     }
	     
	     
         return dateValid, errorGeneral  
	}
	
     func ObtainBINfromCard (parmcard string) (string,string){
	      var bin string
	      var errorGeneral string
	      	      
	      bin = ""
	      errorGeneral=""
	      
	     if len(parmcard) > 6 {
	          bin= parmcard[:6]
	     }else{
	     	errorGeneral ="Card length is invalid, less than 6 "
	     }
	     
	     
         return bin, errorGeneral  
	}
	
     func ObtainLast4fromCard (parmcard string) (string,string){
	      var last string
	      var errorGeneral string
	      	      
	      last = ""
	      errorGeneral=""
	      log.Print(" ObtainLast4fromCard:"+parmcard)
	     if len(parmcard) > 6 {
	          last= parmcard[len(parmcard)-4:]
	          log.Print(" ObtainLast4fromCard 01:"+last)
	     }else{
	     	errorGeneral ="Card length is invalid, less than 6 "
	          log.Print("Card length is invalid, less than 6 ")
	     }
	     
	     
         return last, errorGeneral  
	}
	

func GetCardType(number string) string {
	return "VISA"
/*
 * 

// visa
var re = new RegExp("^4");
 if (number.match(re) != null)
     return "Visa"; 
     
// Mastercard
 // Updated for Mastercard 2017 BINs expansion
 if (/^(5[1-5][0-9]{14}|2(22[1-9][0-9]{12}|2[3-9][0-9]{13}|[3-6][0-9]{14}|7[0-1][0-9]{13}|720[0-9]{12}))$/.test(number)) 
  return "Mastercard"; 
  
  // AMEX 
  re = new RegExp("^3[47]");
  if (number.match(re) != null) return "AMEX";
  // Discover
  re = new RegExp("^(6011|622(12[6-9]|1[3-9][0-9]|[2-8][0-9]{2}|9[0-1][0-9]|92[0-5]|64[4-9])|65)");
  if (number.match(re) != null)
  return "Discover";
  // Diners
  re = new RegExp("^36"); 
  if (number.match(re) != null)
    return "Diners"; 
    // Diners - Carte Blanche 
    
    re = new RegExp("^30[0-5]");
    if (number.match(re) != null)
    return "Diners - Carte Blanche";
    
    // JCB
    re = new RegExp("^35(2[89]|[3-8][0-9])");
    if (number.match(re) != null)
    return "JCB";
    // Visa Electron
    re = new RegExp("^(4026|417500|4508|4844|491(3|7))"); 
    if (number.match(re) != null) 
    return "Visa Electron";
    return ""; 
 */
}
