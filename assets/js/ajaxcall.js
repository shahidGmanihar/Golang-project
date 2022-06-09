$(document).ready(function () {
    $('#country_dropdown').on('change', function () {
         $('#state_dropdown').empty();
         $('#city_dropdown').empty();
         $('#state_dropdown').html('<option value="">Select State </option>');
         var country_id = this.value;
         $.ajax({
              url: "/states",
              type: "POST",
              data: {
                   country_id: country_id
              },
              dataType: 'json',
              cache: false,
              success: function (result) {
                   console.log(result);
                   var selOpts = "";
                   for (i = 0; i < result.length; i++) {
                        var id = result[i]['Stateid'];
                        var val = result[i]['Statename'];
                        selOpts += "<option  value='" + id + "' >" + val + "</option>";
                   }
                   $('#state_dropdown').append(selOpts);
                   $('#city_dropdown').html('<option value="">Select State First</option>');
              }
         });
    });
    $('#state_dropdown').on('change', function () {
         $('#city_dropdown').empty();
         $('#city_dropdown').html('<option value="">Select City </option>');
         var state_id = this.value;
         $.ajax({
              url: "/cities",
              type: "POST",
              data: {
                   state_id: state_id
              },
              dataType: 'json',
              cache: false,
              success: function (result) {
                   console.log(result);
                   var selOpts = "";
                   for (i = 0; i < result.length; i++) {
                        var id = result[i]['Cityid'];
                        var val = result[i]['Cityname'];
                        selOpts += "<option   value='" + id + "' >" + val + "</option>";
                   }
                   $('#city_dropdown').append(selOpts);


              }
         });
    });
  

   
});

