function HttpDns (gateway) {
    this.gateway = gateway;
}

HttpDns.prototype.now = function() {
    return parseInt(Date.parse(new Date()) / 1000);
};

HttpDns.prototype.updateDnsCache = function(region, serviceName) {
    let httpDns = this;
    fetch(httpDns.gateway + '/ips?region=' + region + '&service-name=' + serviceName, {"METHOD": "GET"}).then(function (response) {
        response.json().then(function (data) {
            if (data.code === 0) {
                let ips = data.data.ips;
                const nowSecond = httpDns.now();

                ips.forEach(function (ip, i) {
                    ips[i].ttl = nowSecond + parseInt(ip.ttl)
                });

                localStorage.setItem('dns:' + region + ':' + serviceName, JSON.stringify(ips));
            }
        });
    });
};

HttpDns.prototype.resolveName = function(region, serviceName, once = false) {
    const ipsStr = localStorage.getItem('dns:' + region + ':' + serviceName);
    if (ipsStr) {
        let ips = JSON.parse(ipsStr);
        const nowSecond = this.now();

        for (let i in ips) {
            if (ips[i].ttl >= nowSecond) {
                return ips[i].ip;
            }
        }
    }

    if (!once) {
        self.updateDnsCache(region, serviceName);
        return this.resolveName(region, serviceName, true)
    }

    return null;
};

// Demo
// let httpDnsIns = new HttpDns('http://localhost:53103');
