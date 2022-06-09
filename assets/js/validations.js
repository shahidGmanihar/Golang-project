$(document).ready(function() { 
 
    $('#btn-submitt').click(function() {  
 
        $(".error").hide();
        var hasError = false;
        //var emailReg = /^([\w-\.]+@([\w-]+\.)+[\w-]{2,4})?$/;
       
        var name = $("#name").val();
        var birthdate = $("#birthdate").val();
        var address = $("#address").val();        
        var pancard = $("#pancard").val();        
        var familydetails = $("#familydetails").val();        
        //var city = $("#city_dropdown").val();        


        if(name == '') {
            $("#Userfname").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error">*Enter your First name.</span>');
            hasError = true;
          } 

          if(birthdate == '') {
            $("#birthdate").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error">*Enter your Last name.</span>');
            hasError = true;
          } 


        if(address == '') {
            $("#address").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error">*Enter your email address.</span>');
            hasError = true;
        }
 
        // else if(!emailReg.test(emailaddressVal)) {
        //     $("#UserEmail").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error">*Enter a valid email address.</span>');
        //     hasError = true;
        // }
 
       
        if(pancard == '') {
            $("#pancard").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error">* Pancard details is required.</span>');
            hasError = true;
          }else if(String(pancard).length>10){
            $("#pancard").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error">* Pancard details should be max 10 character</span>');
            hasError = true;

          }

          if(familydetails == '') {
            $("#familydetails").after('<span style="color: rgb(228, 60, 18); font-size: 10pt" class="error">*Select state.</span>');
            hasError = true;
          } 

         
         
          if(hasError == true) { return false; }



      


    });
});