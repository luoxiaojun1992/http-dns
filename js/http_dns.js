fetch('http://localhost:52273/ips?region=sh&service-name=user-service', {"METHOD":"GET"}).
then(function(response){
    response.json().
    then(function(data) {
        localStorage.setItem('dns:sh:user-service', data)
    });
});
