function updateDnsCache(region, serviceName) {
    fetch('http://localhost:57267/ips?region=' + region + '&service-name=' + serviceName, {"METHOD": "GET"}).then(function (response) {
        response.json().then(function (data) {
            if (data.code === 0) {
                let ips = data.data.ips;
                ips.forEach(function (ip, i) {
                    ips[i].ttl = parseInt(Date.parse(new Date()) / 1000) + parseInt(ip.ttl)
                });

                localStorage.setItem('dns:' + region + ':' + serviceName, JSON.stringify(ips));
            }
        });
    });
}

function resolveName(region, serviceName, once = false) {
    const ipsStr = localStorage.getItem('dns:' + region + ':' + serviceName);
    if (ipsStr) {
        let ips = JSON.parse(ipsStr);
        const nowSecond = parseInt(Date.parse(new Date()) / 1000);

        for (let i in ips) {
            if (ips[i].ttl >= nowSecond) {
                return ips[i].ip;
            }
        }
    }

    if (!once) {
        updateDnsCache(region, serviceName);
        return resolveName(region, serviceName, true)
    }

    return null;
}
