function HttpDns (gateway) {
    this.gateway = gateway;
}

HttpDns.prototype.now = function() {
    return parseInt(Date.parse(new Date()) / 1000);
};

HttpDns.prototype.storageKey = function(region, serviceName) {
    return 'http-dns:' + region + ':' + serviceName;
};

HttpDns.prototype.updateDnsCache = function(region, serviceName, callback) {
    let httpDns = this;
    fetch(httpDns.gateway + '/ips?region=' + region + '&service-name=' + serviceName, {"METHOD": "GET"}).then(function (response) {
        response.json().then(function (data) {
            if (data.code === 0) {
                let ips = data.data.ips;
                const nowSecond = httpDns.now();

                ips.forEach(function (ip, i) {
                    ips[i].ttl = nowSecond + parseInt(ip.ttl)
                });

                localStorage.setItem(httpDns.storageKey(region, serviceName), JSON.stringify(ips));

                callback();
            }
        });
    });
};

HttpDns.prototype.resolveName = function(region, serviceName, callback, once = false) {
    const ipsStr = localStorage.getItem(this.storageKey(region, serviceName));
    if (ipsStr) {
        let ips = JSON.parse(ipsStr);
        const nowSecond = this.now();

        for (let i in ips) {
            if (ips[i].ttl >= nowSecond) {
                callback(ips[i].ip);
                return;
            }
        }
    }

    if (!once) {
        let httpDns = this;
        this.updateDnsCache(region, serviceName, function () {
            httpDns.resolveName(region, serviceName, callback, true)
        });
    }
};

HttpDns.prototype.fetch = function (url, options, callback) {
    let urlObj = new URL(url);
    this.resolveName("sh", urlObj.hostname, function (ip) {
        if (typeof options.Headers == 'undefined') {
            options.Headers = {"Host": urlObj.hostname};
        } else {
            options.Headers.Host = urlObj.hostname;
        }
        const newUrl = urlObj.protocol + "//" + ip + (urlObj.port.length > 0 ? ':' + urlObj.port : "") + urlObj.pathname + urlObj.search + urlObj.hash;
        callback(fetch(newUrl, options));
    });
};

// Demo
// let httpDnsIns = new HttpDns('http://localhost:61596');
