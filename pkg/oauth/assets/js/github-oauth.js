function githubOauth() {
    // 1: oath login
    const params = {
        client_id: 'Iv1.27e4e24d5cf774cc',
        redirect_uri: 'http://localhost:8080/v1/oauth/callback/github',
    }
    const options = { height: 800, width: 1200 }
    const url = `https://github.com/login/oauth/authorize`
    const id = "github-oauth"

    var code, user, err

    login(params, options, url, id).then(
        res => {
            code = res.code;
            user = res.user;
            document.getElementById("user").innerHTML = res.user
            document.getElementById("code").innerHTML = res.code
        })


    // 2: list available orgs and packages
}

const toParams = query => {
    const q = query.replace(/^\??\//, '');

    return q.split('&').reduce((values, param) => {
        const [key, value] = param.split('=');

        values[key] = value;

        return values;
    }, {});
};

const toQuery = (params, delimiter = '&') => {
    const keys = Object.keys(params);

    return keys.reduce((str, key, index) => {
        let query = `${str}${key}=${params[key]}`;

        if (index < keys.length - 1) {
            query += delimiter;
        }

        return query;
    }, '');
};

class PopupWindow {
    constructor(params, options, url, id) {
        this.id = id;
        this.url = url + '?' + toQuery(params);
        this.options = options;
    }

    open() {
        const { url, id, options } = this;
        this.window = window.open(url, id, toQuery(options, ','));
    }

    close() {
        this.cancel();
        this.window.close();
    }

    poll() {
        this.promise = new Promise((resolve, reject) => {
            this._iid = window.setInterval(() => {
                try {
                    const popup = this.window;
                    // pop up closed, finish promise
                    if (!popup || popup.closed !== false) {
                        this.close();
                        reject(new Error('The popup was closed'));
                        return;
                    }

                    // empty path name, finish promise
                    if (popup.location.href === this.url || popup.location.pathname === 'blank') {
                        return;
                    }

                    // returned from callback, code and user will be returned from server.
                    const params = toParams(popup.location.search.replace(/^\?/, ''));

                    // server already list and force user select installations, we are done here.
                    // resolve params
                    resolve(params);

                    // sleep for 2 seconds in order to load success image.
                    sleep(5000)

                    // close window
                    this.close();
                } catch (error) {}
            }, 2000);
        });
    }

    cancel() {
        if (this._iid) {
            window.clearInterval(this._iid);
            this._iid = null;
        }
    }

    then(...args) {
        return this.promise.then(...args);
    }

    catch(...args) {
        return this.promise.then(...args);
    }

    static open(...args) {
        const popup = new this(...args);

        popup.open();
        popup.poll();

        return popup;
    }
}

const login = (params, options, url, id) => {
    return new Promise(  async (resolve, reject) => {
        const popup = PopupWindow.open(params, options, url, id);
        popup.then(resolve, reject);
    });
};

function sleep(milliseconds) {
    const date = Date.now();
    let currentDate = null;
    do {
        currentDate = Date.now();
    } while (currentDate - date < milliseconds);
}
