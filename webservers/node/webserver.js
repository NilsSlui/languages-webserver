const http = require('http');
const crypto = require('crypto');

const server = http.createServer((req, res) => {
    if (req.method === 'POST') {
        let body = '';
        req.on('data', chunk => {
            body += chunk.toString();
        });

        req.on('end', () => {
            const input = new URLSearchParams(body).get('input');
            if (!input) {
                res.writeHead(500, { 'Content-Type': 'text/plain' });
                res.end('Input is required');
                return;
            }

            const urlParts = req.url.split('/');
            switch (urlParts[1]) {
                case 'sha256':
                    res.writeHead(200, { 'Content-Type': 'text/plain' });
                    res.end(crypto.createHash('sha256').update(input).digest('hex'));
                    break;
                case 'base64':
                    res.writeHead(200, { 'Content-Type': 'text/plain' });
                    res.end(Buffer.from(input).toString('base64'));
                    break;
                case 'urlencode':
                    res.writeHead(200, { 'Content-Type': 'text/plain' });
                    res.end(encodeURIComponent(input));
                    break;
                default:
                    res.writeHead(404, { 'Content-Type': 'text/plain' });
                    res.end('Not Found');
            }
        });
    } else {
        res.writeHead(201, { 'Content-Type': 'text/plain' });
        res.end('node');
    }
});

server.listen(8008, () => {
    console.log('Server is running on port 8008');
});
