<?php
$uri = $_SERVER['REQUEST_URI'];
$method = $_SERVER['REQUEST_METHOD'];

try {
    if ($method === 'POST') {
        $input = $_POST['input'] ?? throw new Exception('Input is required');
        echo match ($uri) {
            '/sha256' => hash('sha256', $input),
            '/base64' => base64_encode($input),
            '/urlencode' => urlencode($input),
            default => http_response_code(404) . 'Not Found',
        };
    } elseif ($method === 'GET') {
        http_response_code(201);
        echo 'php';
    } else {
        http_response_code(405);
        echo 'Method Not Allowed';
    }
} catch (Exception $e) {
    http_response_code(500);
    echo $e->getMessage() ?? 'Internal Server Error';
}
