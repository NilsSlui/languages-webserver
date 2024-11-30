const http = require('http');
const url = require('url');
const crypto = require('crypto');

const server = http.createServer((req, res) => {
    const parsedUrl = url.parse(req.url, true);
    const method = req.method;
    const uri = parsedUrl.pathname;

    try {
        if (method === 'POST') {
            let body = '';
            req.on('data', chunk => {
                body += chunk;
            });

            req.on('end', () => {
                const input = new URLSearchParams(body).get('input');
                if (!input) {
                    throw new Error('Input is required');
                }

                let result;
                switch (uri) {
                    case '/sha256':
                        result = crypto.createHash('sha256').update(input).digest('hex');
                        break;
                    case '/base64':
                        result = Buffer.from(input).toString('base64');
                        break;
                    case '/urlencode':
                        result = encodeURIComponent(input);
                        break;
                    default:
                        res.writeHead(404, { 'Content-Type': 'text/plain' });
                        return res.end('Not Found');
                }

                res.writeHead(200, { 'Content-Type': 'text/plain' });
                res.end(result);
            });
        } else if (method === 'GET') {
            res.writeHead(201, { 'Content-Type': 'text/plain' });
            res.end('node');
        } else {
            res.writeHead(405, { 'Content-Type': 'text/plain' });
            res.end('Method Not Allowed');
        }
    } catch (error) {
        res.writeHead(500, { 'Content-Type': 'text/plain' });
        res.end(error.message || 'Internal Server Error');
    }
});

server.listen(8008, () => {
    console.log('Server running');
});
