// server.js
const http = require('http');
const url = require('url');
const crypto = require('crypto');

const server = http.createServer((req, res) => {
    const method = req.method;
    const parsedUrl = url.parse(req.url, true);
    const uri = parsedUrl.pathname;

    if (method === 'POST') {
        let body = '';

        req.on('data', chunk => {
            body += chunk;
        });

        req.on('end', () => {
            try {
                // Assume content-type is application/x-www-form-urlencoded
                const input = new URLSearchParams(body).get('input');

                if (!input) {
                    res.statusCode = 400;
                    res.end('Input is required');
                    return;
                }

                res.statusCode = 200;

                if (uri === '/sha256') {
                    const hash = crypto.createHash('sha256').update(input).digest('hex');
                    res.end(hash);
                } else if (uri === '/base64') {
                    const base64 = Buffer.from(input).toString('base64');
                    res.end(base64);
                } else if (uri === '/urlencode') {
                    const urlencode = encodeURIComponent(input);
                    res.end(urlencode);
                } else {
                    res.statusCode = 404;
                    res.end('Not Found');
                }
            } catch (err) {
                res.statusCode = 500;
                res.end('Internal Server Error');
                console.error(err);
            }
        });
    } else if (method === 'GET') {
        res.statusCode = 201;
        res.end('node');
    } else {
        res.statusCode = 405;
        res.end('Method Not Allowed');
    }
});

server.listen(8008, () => {
    console.log('Server is listening on port 8008');
});
