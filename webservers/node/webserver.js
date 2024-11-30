const http = require('http');
const url = require('url');
const querystring = require('querystring');
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
                const parsedBody = querystring.parse(body);
                const input = parsedBody['input'];

                if (!input) {
                    throw new Error('Input is required');
                }

                switch (uri) {
                    case '/sha256':
                        const hash = crypto.createHash('sha256').update(input).digest('hex');
                        res.statusCode = 200;
                        res.end(hash);
                        break;
                    case '/base64':
                        const base64 = Buffer.from(input).toString('base64');
                        res.statusCode = 200;
                        res.end(base64);
                        break;
                    case '/urlencode':
                        const urlencode = encodeURIComponent(input);
                        res.statusCode = 200;
                        res.end(urlencode);
                        break;
                    default:
                        res.statusCode = 404;
                        res.end('Not Found');
                }
            } catch (err) {
                res.statusCode = 500;
                res.end(err.message || 'Internal Server Error');
            }
        });
    } else if (method === 'GET') {
        try {
            res.statusCode = 201;
            res.end('php');
        } catch (err) {
            res.statusCode = 500;
            res.end(err.message || 'Internal Server Error');
        }
    } else {
        res.statusCode = 405;
        res.end('Method Not Allowed');
    }
});

server.listen(8008, () => {
    console.log('Server is listening on port 8008');
});
